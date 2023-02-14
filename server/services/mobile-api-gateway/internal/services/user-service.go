package userservice

import (
	"context"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	userclient "github.com/vantoan19/Petifies/server/services/grpc-clients/user-client"
	"google.golang.org/grpc"
)

var logger = logging.New("MobileGateway.Service")

type UserConfiguration func(us *userService) error

type userService struct {
	userClient userclient.UserClient
}

type UserService interface {
	CreateUser(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error)
	Login(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error)
}

func NewUserService(conn *grpc.ClientConn, cfgs ...UserConfiguration) (UserService, error) {
	us := &userService{
		userClient: userclient.New(conn),
	}
	for _, cfg := range cfgs {
		err := cfg(us)
		if err != nil {
			return nil, err
		}
	}
	return us, nil
}

func (s *userService) CreateUser(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error) {
	logger.Info("Forwarding CreateUser request to User Service")
	resp, err := s.userClient.CreateUserForward(ctx, req)
	if err != nil {
		return nil, err
	}
	logger.Info("Returning CreateUser response from User Service")
	return resp, nil
}

func (s *userService) Login(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error) {
	logger.Info("Forwarding Login request to User Service")
	resp, err := s.userClient.LoginForward(ctx, req)
	if err != nil {
		return nil, err
	}
	logger.Info("Returning Login response from User Service")
	return resp, nil
}
