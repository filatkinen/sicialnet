package internalhttp

import (
	"context"
	"errors"
	"github.com/filatkinen/socialnet/internal/config/dialog"
	"github.com/filatkinen/socialnet/internal/dialog/app"
	pb "github.com/filatkinen/socialnet/internal/grpc/dialog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"time"
)

type DialogService struct {
	srv        *http.Server
	log        *log.Logger
	config     dialog.Config
	httplog    *httplog
	app        *dialogapp.DialogApp
	reqCounter *RID
	promData   *promData

	pbclient pb.DialogClient
	pbconn   *grpc.ClientConn
}

func NewServer(config dialog.Config, log *log.Logger) (*DialogService, error) {
	httpsrv := &http.Server{
		Addr:              net.JoinHostPort(config.ServerAddress, config.ServerPort),
		ReadHeaderTimeout: time.Second * 10,
	}
	hlog := newHTTPLogger(config.ServerHTTPLogfile, log)

	app, err := dialogapp.New(log, config)
	if err != nil {
		return nil, err
	}

	s := &DialogService{
		srv:        httpsrv,
		log:        log,
		config:     config,
		httplog:    hlog,
		app:        app,
		reqCounter: NewRID(),
	}
	s.promData = NewPromData()
	s.srv.Handler = s.NewRouter()

	conn, err := grpc.Dial(net.JoinHostPort(s.config.ServiceGRPCAddress, s.config.ServiceGRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	s.pbconn = conn
	s.pbclient = pb.NewDialogClient(conn)
	return s, nil
}

func (s *DialogService) Start() error {
	s.log.Printf("Starting HTTP server at:%s", net.JoinHostPort(s.config.ServerAddress, s.config.ServerPort))
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Printf("Error starting HTTP server at:%s with error:%s\n",
			net.JoinHostPort(s.config.ServerAddress, s.config.ServerPort), err)
		return err
	}
	return nil
}

func (s *DialogService) Stop(ctx context.Context) error {
	var errHTTP error

	s.log.Println("Stopping HTTP...")
	if errHTTP = s.srv.Shutdown(ctx); errHTTP != nil {
		s.log.Printf("HTTP shutdown error: %s\n", errHTTP)
	} else {
		s.log.Println("HTTP graceful shutdown complete.")
	}

	return errHTTP
}

func (s *DialogService) Close() error {
	var errApp, errLog, errGRPC error

	errApp = s.app.Close(context.Background())

	if errLog = s.httplog.close(); errLog != nil {
		s.log.Printf("error closing httplog %s\n", errLog)
	}

	s.log.Println("Closing GRPC...")
	if errGRPC = s.pbconn.Close(); errGRPC != nil {
		s.log.Printf("GRPC closing error: %s\n", errGRPC)
	} else {
		s.log.Println("GRPC closing  complete.")
	}

	return errors.Join(errApp, errLog, errGRPC)
}
