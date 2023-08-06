package internalhttp

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
)

type promResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

type promData struct {
	totalRequests  *prometheus.CounterVec
	responseStatus *prometheus.CounterVec
	httpDuration   *prometheus.HistogramVec
}

func NewPromData() *promData {

	pr := promData{
		totalRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goapp_http_requests_total",
				Help: "Number of get requests.",
			},
			[]string{"path"},
		),
		responseStatus: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "goapp_response_status",
				Help: "Status of HTTP response",
			},
			[]string{"status"},
		),
		httpDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "goapp_http_response_time_seconds",
			Help: "Duration of HTTP requests.",
		}, []string{"path"}),
	}
	prometheus.Register(pr.totalRequests)
	prometheus.Register(pr.httpDuration)
	prometheus.Register(pr.totalRequests)
	return &pr
}

func NewPromResponseWriter(w http.ResponseWriter) *promResponseWriter {
	return &promResponseWriter{w, http.StatusOK}
}

func (rw *promResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *Server) prometheusMiddleware(next http.Handler, routePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		timer := prometheus.NewTimer(s.promData.httpDuration.WithLabelValues(routePath))
		rw := NewPromResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		s.promData.responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		s.promData.totalRequests.WithLabelValues(routePath).Inc()

		timer.ObserveDuration()
	})
}
