package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	publicProtoV1 "github.com/vantoan19/Petifies/proto/public-gateway/v1"
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	publicEndpointV1 "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/endpoints/grpc/v1"
)

type gRPCPublicServer struct {
	createUser grpctransport.Handler
	login      grpctransport.Handler
}

func NewPublicServer(endpoints publicEndpointV1.UserEndpoints) publicProtoV1.PublicGatewayServer {
	return &gRPCPublicServer{
		createUser: grpctransport.NewServer(
			endpoints.CreateUser,
			common.CreateClientForwardDecodeRequestFunc[*commonProto.CreateUserRequest](),
			common.CreateClientForwardEncodeResponseFunc[*commonProto.User](),
		),
		login: grpctransport.NewServer(
			endpoints.Login,
			common.CreateClientForwardDecodeRequestFunc[*commonProto.LoginRequest](),
			common.CreateClientForwardEncodeResponseFunc[*commonProto.LoginResponse](),
		),
	}
}

// ============================================================
// Definition of the endpoints
// ============================================================

// Create User endpoint

func (s *gRPCPublicServer) CreateUser(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.User), nil
}

// Login endpoint

func (s *gRPCPublicServer) Login(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error) {
	_, resp, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.LoginResponse), nil
}
