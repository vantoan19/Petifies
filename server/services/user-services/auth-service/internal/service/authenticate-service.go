package authservice

import (
	"context"
)

// var logger = logging.New("AuthService.Authenticate")

type authenticateService struct{}

type AuthenticateService interface {
	SayHello(ctx context.Context) (string, error)
}

func NewAuthenticateService() AuthenticateService {
	return &authenticateService{}
}

func (s *authenticateService) SayHello(ctx context.Context) (string, error) {
	return "Hello World", nil
}
