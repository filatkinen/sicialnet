package internalhttp

import (
	"context"
	pb "github.com/filatkinen/socialnet/internal/grpc/dialog"
)

func (s *DialogService) UserLogin(ctx context.Context, userID string, userPassword string) (string, error) {
	token, err := s.pbclient.GetToken(ctx, &pb.Cred{
		UserID:   userID,
		Password: userPassword,
		Logging:  "client ID=",
	})
	if err != nil {
		return "", err
	}
	return token.Token, err
}

func (s *DialogService) CheckToken(ctx context.Context, token string) (string, error) {
	user, err := s.pbclient.CheckToken(ctx, &pb.Token{Token: token})
	if err != nil {
		return "", err
	}
	return user.UserID, err
}
