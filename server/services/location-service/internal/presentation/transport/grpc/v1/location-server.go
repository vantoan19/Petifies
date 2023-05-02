package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	locationProtoV1 "github.com/vantoan19/Petifies/proto/location-service/v1"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/location-service/internal/presentation/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/location-service/pkg/translators"
)

type gRPCLocationServer struct {
	listNearByLocationsByType grpctransport.Handler
}

func NewLocationServer(
	locationEndpoints endpointsV1.LocationEndpoints,
) locationProtoV1.LocationServiceServer {
	return &gRPCLocationServer{
		listNearByLocationsByType: grpctransport.NewServer(
			locationEndpoints.ListNearByLocationsByType,
			translators.DecodeListNearByLocationByTypeRequest,
			translators.EncodeListNearByLocationsByTypeResponse,
		),
	}
}

func (s *gRPCLocationServer) ListNearByLocationsByType(
	ctx context.Context,
	req *locationProtoV1.ListNearByLocationsByTypeRequest,
) (*locationProtoV1.ListNearByLocationsByTypeResponse, error) {
	_, resp, err := s.listNearByLocationsByType.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*locationProtoV1.ListNearByLocationsByTypeResponse), nil
}
