package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	userEndpointsV1 "github.com/vantoan19/Petifies/server/services/user-service/internal/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/translator"
)

type gRPCUserServer struct {
	createUser grpctransport.Handler
	login      grpctransport.Handler
}

func NewUserServer(endpoints userEndpointsV1.UserEndpoints) userProtoV1.UserServiceServer {
	return &gRPCUserServer{
		createUser: grpctransport.NewServer(
			endpoints.CreateUser,
			translator.DecodeCreateUserRequest,
			translator.EncodeCreateUserResponse,
		),
		login: grpctransport.NewServer(
			endpoints.Login,
			translator.DecodeLoginRequest,
			translator.EncodeLoginResponse,
		),
	}
}

// ============================================================
// Definition of the endpoints
// ============================================================

// Create User endpoint

func (s *gRPCUserServer) CreateUser(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.User), nil
}

// Login endpoint

func (s *gRPCUserServer) Login(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error) {
	_, resp, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.LoginResponse), nil
}
