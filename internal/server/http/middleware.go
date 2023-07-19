package internalhttp

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

type ContextKey string

const (
	ContextUserKey ContextKey = "user_id"
	RequestID      ContextKey = "request_id"
)

type RID struct {
	reqCounter int
	lock       sync.Mutex
}

func NewRID() *RID {
	return &RID{reqCounter: 0}
}

func (r *RID) nextRequestID() int {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.reqCounter++
	return r.reqCounter
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

func (s *Server) Logger(next http.Handler, routeName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)

		tstart := time.Now()
		next.ServeHTTP(lrw, r)
		code := lrw.statusCode
		size := lrw.size
		timeToTakeServ := time.Since(tstart)
		raddrlog := r.RemoteAddr[0:strings.Index(r.RemoteAddr, ":")]
		timelog := tstart.UTC().Format("02/01/2006 15:04:05 UTC")

		rID := r.Context().Value(RequestID).(string)
		s.httplog.Httplogger.Printf("%s [%s] reqID=%s (%s) %s %s %s %d %d %s %s\n",
			raddrlog, timelog, rID, routeName, r.Method, r.URL.Path, r.Proto, code, size, timeToTakeServ, r.UserAgent())
	})
}

func (s *Server) CheckSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := "Bearer "
		authHeader := r.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, prefix)
		if authHeader == "" || reqToken == authHeader || len(authHeader) != 26 {
			s.ClientError(w, http.StatusBadRequest, "bad auth header")
			return
		}
		userID, err := s.app.CheckToken(r.Context(), reqToken)
		if err != nil {
			s.ClientError(w, http.StatusBadRequest, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) addRequestCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RequestID, strconv.Itoa(s.reqCounter.nextRequestID()))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
