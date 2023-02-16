package userclient

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/translator"
)

const (
	serviceName = "UserService"
)

var logger = logging.New("Clients.UserClient")

type userClient struct {
	createUser        endpoint.Endpoint
	createUserForward endpoint.Endpoint
	login             endpoint.Endpoint
	loginForward      endpoint.Endpoint
}

type UserClient interface {
	CreateUser(ctx context.Context, email, password, firstName, lastName string) (*models.CreateUserResp, error)
	CreateUserForward(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	LoginForward(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error)
}

func New(conn *grpc.ClientConn) UserClient {
	return &userClient{
		createUser: grpctransport.NewClient(
			conn,
			serviceName,
			"CreateUser",
			translator.EncodeCreateUserRequest,
			translator.DecodeCreateUserResponse,
			commonProto.User{},
		).Endpoint(),
		createUserForward: grpctransport.NewClient(
			conn,
			serviceName,
			"CreateUser",
			common.CreateClientForwardEncodeRequestFunc[*commonProto.CreateUserRequest](),
			common.CreateClientForwardDecodeResponseFunc[*commonProto.User](),
			commonProto.User{},
		).Endpoint(),
		login: grpctransport.NewClient(
			conn,
			serviceName,
			"Login",
			translator.EncodeLoginRequest,
			translator.DecodeLoginResponse,
			commonProto.LoginResponse{},
		).Endpoint(),
		loginForward: grpctransport.NewClient(
			conn,
			serviceName,
			"Login",
			common.CreateClientForwardEncodeRequestFunc[*commonProto.LoginRequest](),
			common.CreateClientForwardDecodeResponseFunc[*commonProto.LoginResponse](),
			commonProto.LoginResponse{},
		).Endpoint(),
	}
}

func (c *userClient) CreateUser(ctx context.Context, email, password, firstName, lastName string) (*models.CreateUserResp, error) {
	logger.Info("Start UserClient.CreateUser")

	req := &models.CreateUserReq{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}
	resp, err := c.createUser(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.CreateUser: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserClient.CreateUser: SUCCESSFUL")
	return resp.(*models.CreateUserResp), nil
}

func (c *userClient) CreateUserForward(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error) {
	logger.Info("Start UserClient.CreateUserForward")

	resp, err := c.createUserForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.CreateUserForward: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserClient.CreateUserForward: SUCCESSFUL")
	return resp.(*commonProto.User), nil
}

func (c *userClient) Login(ctx context.Context, email, password string) (string, error) {
	logger.Info("Start UserClient.Login")

	req := &models.LoginReq{
		Email:    email,
		Password: password,
	}
	resp, err := c.login(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.Login: FAILED", logging.Data{"error": err.Error()})
		return "", err
	}

	logger.Info("Finished UserClient.Login: SUCCESSFUL")
	loginResp := resp.(*models.LoginResp)
	return loginResp.AccessToken, nil
}

func (c *userClient) LoginForward(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error) {
	logger.Info("Start UserClient.LoginForward")

	resp, err := c.loginForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.LoginForward: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserClient.LoginForward: SUCCESSFUL")
	return resp.(*commonProto.LoginResponse), nil
}
