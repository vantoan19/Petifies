package userservice

import (
	"context"
)

// var logger = logging.NewLogger("UserService.Service")

type userService struct{}

type UserService interface {
	SayHello(ctx context.Context) (string, error)
}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) SayHello(ctx context.Context) (string, error) {
	return "Hello World", nil
}
