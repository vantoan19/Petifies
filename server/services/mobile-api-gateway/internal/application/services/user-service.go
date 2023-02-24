package services

import (
	"context"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	userclient "github.com/vantoan19/Petifies/server/services/grpc-clients/user-client"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
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
	RefreshToken(ctx context.Context, req *commonProto.RefreshTokenRequest) (*commonProto.RefreshTokenResponse, error)
	GetMyInfo(ctx context.Context) (*models.User, error)
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

func (s *userService) RefreshToken(ctx context.Context, req *commonProto.RefreshTokenRequest) (*commonProto.RefreshTokenResponse, error) {
	logger.Info("Start UserService.RefreshToken")

	logger.Info("Executing UserService.RefreshToken: forwarding the request to UserService")
	resp, err := s.userClient.RefreshTokenForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserService.RefreshToken: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserService.RefreshToken: SUCCESSFUL")
	return resp, nil
}

func (s *userService) GetMyInfo(ctx context.Context) (*models.User, error) {
	logger.Info("Start UserService.GetMyInfo")

	logger.Info("Executing UserService.GetMyInfo: forwarding the request to UserService")
	userID, err := commonutils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.userClient.GetUser(ctx, userID)
	if err != nil {
		logger.ErrorData("Finished UserService.GetMyInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserService.GetMyInfo: SUCCESSFUL")
	return resp, nil
}
