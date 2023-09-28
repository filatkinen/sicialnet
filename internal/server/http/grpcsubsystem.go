package internalhttp

import (
	"context"
	"github.com/filatkinen/socialnet/internal/grpc/dialog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"strings"
	"time"
)

func (s *Server) GetToken(ctx context.Context, cred *dialog.Cred) (*dialog.Token, error) {
	token, err := s.app.UserLogin(ctx, cred.UserID, cred.Password)
	if err != nil {
		return nil, err
	}
	return &dialog.Token{Token: token}, nil
}

func (s *Server) CheckToken(ctx context.Context, token *dialog.Token) (*dialog.UserID, error) {
	userID, err := s.app.CheckToken(ctx, token.Token)
	if err != nil {
		return nil, err
	}
	return &dialog.UserID{UserID: userID}, nil

}

func (s *Server) eventUnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	defer s.loggingGRPC(ctx, time.Now(), info.FullMethod)
	m, err := handler(ctx, req)
	return m, err
}

func (s *Server) loggingGRPC(ctx context.Context, timeStart time.Time, method string) {
	timeToTakeServ := time.Since(timeStart)
	timelog := timeStart.UTC().Format("02/01/2006 15:04:05 UTC")
	var remoteIP string
	var ua string
	r, ok := peer.FromContext(ctx)
	if ok {
		remoteIP = r.Addr.String()
		remoteIP = remoteIP[0:strings.Index(remoteIP, ":")]
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		_ua := md["user-agent"]
		if len(_ua) > 0 {
			ua = _ua[0]
		}
	}
	s.httplog.Httplogger.Printf("%s [%s] %s %s %s %s %s %s %s\n",
		remoteIP, timelog, "RPC", method, "-", "-", "-", timeToTakeServ, ua)
}
