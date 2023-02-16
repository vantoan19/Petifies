package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/user-service/internal/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/translator"
)

type gRPCAuthServer struct {
	verifyToken grpctransport.Handler
}

func NewAuthServer(endpoints endpointsV1.AuthEndpoints) userProtoV1.AuthServiceServer {
	return &gRPCAuthServer{
		verifyToken: grpctransport.NewServer(
			endpoints.VerifyToken,
			translator.DecodeVerifyTokenRequest,
			translator.EncodeVerifyTokenResponse,
		),
	}
}

// ============================================================
// Definition of the endpoints
// ============================================================

// Create User endpoint

func (s *gRPCAuthServer) VerifyToken(ctx context.Context, req *userProtoV1.VerifyTokenRequest) (*userProtoV1.VerifyTokenResponse, error) {
	_, resp, err := s.verifyToken.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*userProtoV1.VerifyTokenResponse), nil
}
