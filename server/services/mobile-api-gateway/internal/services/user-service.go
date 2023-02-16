package services

import (
	"context"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	userclient "github.com/vantoan19/Petifies/server/services/grpc-clients/user-client"
	"google.golang.org/grpc"
)

var logger = logging.New("MobileGateway.UserService")

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
	logger.Info("Start UserService.CreateUser")

	logger.Info("Executing UserService.CreateUser: forwarding the request to UserService")
	resp, err := s.userClient.CreateUserForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserService.CreateUser: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserService.CreateUser: SUCCESSFUL")
	return resp, nil
}

func (s *userService) Login(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error) {
	logger.Info("Start UserService.Login")

	logger.Info("Executing UserService.Login: forwarding the request to UserService")
	resp, err := s.userClient.LoginForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserService.Login: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserService.Login: SUCCESSFUL")
	return resp, nil
}
