package internalhttp

import (
	"context"
	"errors"
	"github.com/filatkinen/socialnet/internal/config/server"
	"github.com/filatkinen/socialnet/internal/grpc/dialog"
	socialapp "github.com/filatkinen/socialnet/internal/server/app"
	"github.com/filatkinen/socialnet/internal/storage/caching"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	srv        *http.Server
	log        *log.Logger
	config     server.Config
	httplog    *httplog
	app        *socialapp.App
	reqCounter *RID
	promData   *promData
	cache      *caching.RedisCache

	conn     *grpc.Server
	connLock sync.Mutex
}

func NewServer(config server.Config, log *log.Logger) (*Server, error) {
	httpsrv := &http.Server{
		Addr:              net.JoinHostPort(config.ServerAddress, config.ServerPort),
		ReadHeaderTimeout: time.Second * 10,
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

func (s *Server) startGRPC() error {
	s.connLock.Lock()
	s.conn = grpc.NewServer()
	s.connLock.Unlock()
	lis, err := net.Listen("tcp", net.JoinHostPort(s.config.ServerAddress, s.config.ServerGRPCPort))
	if err != nil {
		return err
	}
	dialog.RegisterDialogServer(s.conn, s)
	if err := s.conn.Serve(lis); err != nil {
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}
	}
	return nil
}

func (s *Server) Start() error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var errGRPC, errHTTP error
	go func() {
		defer wg.Done()
		log.Printf("Starting  GRPC subsystem...\n")
		errGRPC = s.startGRPC()
		if errGRPC != nil {
			log.Printf("Failed to start GRPC server: %s ", errGRPC.Error())
		}
	}()
	go func() {
		defer wg.Done()
		s.log.Printf("Starting HTTP server at:%s", net.JoinHostPort(s.config.ServerAddress, s.config.ServerPort))
		errHTTP = s.srv.ListenAndServe()
		if errHTTP != nil && !errors.Is(errHTTP, http.ErrServerClosed) {
			s.log.Printf("Error starting HTTP server at:%s with error:%s\n",
				net.JoinHostPort(s.config.ServerAddress, s.config.ServerPort), errHTTP)
		}
	}()
	wg.Wait()
	return errors.Join(errHTTP, errGRPC)
}

func (s *Server) Stop(ctx context.Context) error {
	log.Printf("Stopping service...\n")
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Printf("HTTP shutdown error: %s\n", err)
		return err
	}

	log.Printf("Stopping sysmon service. GRPC subsystem...\n")
	s.connLock.Lock()
	s.conn.Stop()
	s.connLock.Unlock()

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
