package userclient

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/translator"
)

const (
	userService = "UserService"
	authService = "AuthService"
)

var logger = logging.New("Clients.UserClient")

type userClient struct {
	createUser         endpoint.Endpoint
	createUserForward  endpoint.Endpoint
	login              endpoint.Endpoint
	loginForward       endpoint.Endpoint
	verifyToken        endpoint.Endpoint
	verifyTokenForward endpoint.Endpoint
}

type UserClient interface {
	CreateUser(ctx context.Context, email, password, firstName, lastName string) (*models.CreateUserResp, error)
	CreateUserForward(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	LoginForward(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error)
	VerifyToken(ctx context.Context, token string) (string, error)
	VerifyTokenForward(ctx context.Context, req *userProtoV1.VerifyTokenRequest) (*userProtoV1.VerifyTokenResponse, error)
}

func New(conn *grpc.ClientConn) UserClient {
	return &userClient{
		createUser: grpctransport.NewClient(
			conn,
			userService,
			"CreateUser",
			translator.EncodeCreateUserRequest,
			translator.DecodeCreateUserResponse,
			commonProto.User{},
		).Endpoint(),
		createUserForward: grpctransport.NewClient(
			conn,
			userService,
			"CreateUser",
			common.CreateClientForwardEncodeRequestFunc[*commonProto.CreateUserRequest](),
			common.CreateClientForwardDecodeResponseFunc[*commonProto.User](),
			commonProto.User{},
		).Endpoint(),
		login: grpctransport.NewClient(
			conn,
			userService,
			"Login",
			translator.EncodeLoginRequest,
			translator.DecodeLoginResponse,
			commonProto.LoginResponse{},
		).Endpoint(),
		loginForward: grpctransport.NewClient(
			conn,
			userService,
			"Login",
			common.CreateClientForwardEncodeRequestFunc[*commonProto.LoginRequest](),
			common.CreateClientForwardDecodeResponseFunc[*commonProto.LoginResponse](),
			commonProto.LoginResponse{},
		).Endpoint(),
		verifyToken: grpctransport.NewClient(
			conn,
			authService,
			"VerifyToken",
			translator.EncodeVerifyTokenRequest,
			translator.DecodeVerifyTokenResponse,
			userProtoV1.VerifyTokenResponse{},
		).Endpoint(),
		verifyTokenForward: grpctransport.NewClient(
			conn,
			authService,
			"VerifyToken",
			common.CreateClientForwardEncodeRequestFunc[*userProtoV1.VerifyTokenRequest](),
			common.CreateClientForwardDecodeResponseFunc[*userProtoV1.VerifyTokenResponse](),
			userProtoV1.VerifyTokenResponse{},
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

func (c *userClient) VerifyToken(ctx context.Context, token string) (string, error) {
	logger.Info("Start UserClient.VerifyToken")

	req := &models.VerifyTokenReq{
		Token: token,
	}
	resp, err := c.verifyToken(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.VerifyToken: FAILED", logging.Data{"error": err.Error()})
		return "", err
	}

	logger.Info("Finished UserClient.VerifyToken: SUCCESSFUL")
	loginResp := resp.(*models.LoginResp)
	return loginResp.AccessToken, nil
}

func (c *userClient) VerifyTokenForward(ctx context.Context, req *userProtoV1.VerifyTokenRequest) (*userProtoV1.VerifyTokenResponse, error) {
	logger.Info("Start UserClient.VerifyTokenForward")

	resp, err := c.verifyTokenForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.VerifyTokenForward: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserClient.VerifyTokenForward: SUCCESSFUL")
	return resp.(*userProtoV1.VerifyTokenResponse), nil
}
