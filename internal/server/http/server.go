package internalhttp

import (
	"context"
	"errors"
	"github.com/filatkinen/socialnet/internal/config/server"
	socialapp "github.com/filatkinen/socialnet/internal/server/app"
	"github.com/filatkinen/socialnet/internal/storage/caching"
	"log"
	"net"
	"net/http"
	"time"
)

const timeDelayRedisCache = time.Second * 2

type Server struct {
	srv        *http.Server
	log        *log.Logger
	config     server.Config
	httplog    *httplog
	app        *socialapp.App
	reqCounter *RID
	promData   *promData
	cache      *caching.RedisCache
}

func NewServer(config server.Config, log *log.Logger) (*Server, error) {
	httpsrv := &http.Server{
		Addr: net.JoinHostPort(config.ServerAddress, config.ServerPort),
	}
	hlog := newHTTPLogger(config.ServerHTTPLogfile, log)

	app, err := socialapp.New(log, config)
	if err != nil {
		return nil, err
	}

	s := &Server{
		srv:        httpsrv,
		log:        log,
		config:     config,
		httplog:    hlog,
		app:        app,
		reqCounter: NewRID(),
	}
	s.promData = NewPromData()
	s.srv.Handler = s.NewRouter()
	s.cache, err = caching.NewCache(config, log)
	if err != nil {
		log.Printf("Error connecting to redis cache. App will not use caching %s\n", err)
	} else {
		log.Print("Using  redis cache for post(with additional postgres db connection)\n")
	}
	return s, nil
}

func (s *Server) Start() error {
	//go func() {
	//	if s.cache != nil {
	//		time.Sleep(timeDelayRedisCache)
	//		s.cache.UpdatePostAll()
	//	}
	//}()
	s.log.Printf("Starting HTTP server at:%s", net.JoinHostPort(s.config.ServerAddress, s.config.ServerPort))
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Printf("Error starting HTTP server at:%s with error:%s\n",
			net.JoinHostPort(s.config.ServerAddress, s.config.ServerPort), err)
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Printf("HTTP shutdown error: %s\n", err)
		return err
	}
	s.log.Println("HTTP graceful shutdown complete.")
	return nil
}

func (s *Server) Close() error {
	err := s.app.Close(context.Background())
	if s.cache != nil {
		s.log.Print("Closing redis. (with additional postgres db connection)\n")
		s.cache.Close()
	}
	if e := s.httplog.close(); e != nil {
		s.log.Printf("error closing httplog %s\n", e)
		err = errors.Join(err, e)
	}
	return err
}
