package userclient

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/translator"
)

const (
	userService = "user_service.v1.UserService"
	authService = "auth_service.v1.AuthService"
)

var logger = logging.New("Clients.UserClient")

type userClient struct {
	createUser          endpoint.Endpoint
	createUserForward   endpoint.Endpoint
	login               endpoint.Endpoint
	loginForward        endpoint.Endpoint
	verifyToken         endpoint.Endpoint
	verifyTokenForward  endpoint.Endpoint
	refreshToken        endpoint.Endpoint
	refreshTokenForward endpoint.Endpoint
	getUser             endpoint.Endpoint
	listUsersByIds      endpoint.Endpoint
}

type UserClient interface {
	CreateUser(ctx context.Context, email, password, firstName, lastName string) (*models.User, error)
	CreateUserForward(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error)
	Login(ctx context.Context, email, password string) (*models.LoginResp, error)
	LoginForward(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error)
	VerifyToken(ctx context.Context, token string) (string, error)
	VerifyTokenForward(ctx context.Context, req *userProtoV1.VerifyTokenRequest) (*userProtoV1.VerifyTokenResponse, error)
	RefreshToken(ctx context.Context, token string) (string, time.Time, error)
	RefreshTokenForward(ctx context.Context, req *commonProto.RefreshTokenRequest) (*commonProto.RefreshTokenResponse, error)
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUsersByIds(ctx context.Context, ids []uuid.UUID) ([]*models.User, error)
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
		refreshToken: grpctransport.NewClient(
			conn,
			authService,
			"RefreshToken",
			translator.EncodeRefreshTokenRequest,
			translator.DecodeRefreshTokenResponse,
			commonProto.RefreshTokenResponse{},
		).Endpoint(),
		refreshTokenForward: grpctransport.NewClient(
			conn,
			authService,
			"RefreshToken",
			common.CreateClientForwardEncodeRequestFunc[*commonProto.RefreshTokenRequest](),
			common.CreateClientForwardDecodeResponseFunc[*commonProto.RefreshTokenResponse](),
			commonProto.RefreshTokenResponse{},
		).Endpoint(),
		getUser: grpctransport.NewClient(
			conn,
			userService,
			"GetUser",
			translator.EncodeGetUserRequest,
			translator.DecodeGetUserResponse,
			commonProto.User{},
		).Endpoint(),
		listUsersByIds: grpctransport.NewClient(
			conn,
			userService,
			"ListUsersByIds",
			translator.EncodeListUsersByIdsRequest,
			translator.DecodeListUsersByIdsResponse,
			commonProto.User{},
		).Endpoint(),
	}
}

func (c *userClient) CreateUser(ctx context.Context, email, password, firstName, lastName string) (*models.User, error) {
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
	return resp.(*models.User), nil
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

func (c *userClient) Login(ctx context.Context, email, password string) (*models.LoginResp, error) {
	logger.Info("Start UserClient.Login")

	req := &models.LoginReq{
		Email:    email,
		Password: password,
	}
	resp, err := c.login(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.Login: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserClient.Login: SUCCESSFUL")
	loginResp := resp.(*models.LoginResp)
	return loginResp, nil
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
	verifyResp := resp.(*models.VerifyTokenResp)
	return verifyResp.UserID, nil
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

func (c *userClient) RefreshToken(ctx context.Context, token string) (string, time.Time, error) {
	logger.Info("Start UserClient.RefreshToken")

	req := &models.RefreshTokenReq{
		RefreshToken: token,
	}
	resp, err := c.refreshToken(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.RefreshToken: FAILED", logging.Data{"error": err.Error()})
		return "", time.Time{}, err
	}

	logger.Info("Finished UserClient.RefreshToken: SUCCESSFUL")
	refreshTokenResp := resp.(*models.RefreshTokenResp)
	return refreshTokenResp.AccessToken, refreshTokenResp.AccessTokenExpiresAt, nil
}

func (c *userClient) RefreshTokenForward(ctx context.Context, req *commonProto.RefreshTokenRequest) (*commonProto.RefreshTokenResponse, error) {
	logger.Info("Start UserClient.RefreshTokenForward")

	resp, err := c.refreshTokenForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.RefreshTokenForward: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserClient.RefreshTokenForward: SUCCESSFUL")
	return resp.(*commonProto.RefreshTokenResponse), nil
}

func (c *userClient) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	logger.Info("Start UserClient.GetMyInfo")

	req := &models.GetUserReq{
		ID: id,
	}
	resp, err := c.getUser(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.GetMyInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserClient.GetMyInfo: SUCCESSFUL")
	user := resp.(*models.User)
	return user, nil
}

func (c *userClient) ListUsersByIds(ctx context.Context, ids []uuid.UUID) ([]*models.User, error) {
	logger.Info("Start UserClient.ListUsersByIds")

	req := &models.ListUsersByIdsReq{
		Ids: ids,
	}
	resp, err := c.listUsersByIds(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserClient.ListUsersByIds: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserClient.ListUsersByIds: SUCCESSFUL")
	users := resp.(*models.ListUsersByIdsResp).Users
	return users, nil
}
