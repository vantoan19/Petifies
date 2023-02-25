package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-gateway/v1"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	endpoints "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/presentation/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/translator"
)

type gRPCAuthServer struct {
	getMyInfo grpctransport.Handler
}

func NewAuthServer(endpoints endpoints.UserEndpoints) authProtoV1.AuthGatewayServer {
	return &gRPCAuthServer{
		getMyInfo: grpctransport.NewServer(
			endpoints.GetMyInfo,
			decodeGetMyInfoRequest,
			translator.EncodeGetUserResponse,
		),
	}
}

// =============================

func decodeGetMyInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	return nil, nil
}

func (s *gRPCAuthServer) GetMyInfo(ctx context.Context, req *authProtoV1.GetMyInfoRequest) (*commonProto.User, error) {
	_, resp, err := s.getMyInfo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.User), nil
}
