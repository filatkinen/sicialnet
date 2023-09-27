package internalhttp

import (
	"context"
	"github.com/filatkinen/socialnet/internal/grpc/dialog"
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
