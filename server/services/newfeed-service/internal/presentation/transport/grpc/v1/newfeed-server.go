package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	newfeedProtoV1 "github.com/vantoan19/Petifies/proto/newfeed-service/v1"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/presentation/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/pkg/translators"
)

type gRPCNewfeedServer struct {
	listPostFeeds  grpctransport.Handler
	listStoryFeeds grpctransport.Handler
}

func NewNewfeedServer(endpoints endpointsV1.NewfeedEndpoints) newfeedProtoV1.NewfeedServiceServer {
	return &gRPCNewfeedServer{
		listPostFeeds: grpctransport.NewServer(
			endpoints.ListPostFeeds,
			translators.DecodeListPostFeedsRequest,
			translators.EncodeListPostFeedsResponse,
		),
		listStoryFeeds: grpctransport.NewServer(
			endpoints.ListStoryFeeds,
			translators.DecodeListStoryFeedsRequest,
			translators.EncodeListStoryFeedsResponse,
		),
	}
}

func (s *gRPCNewfeedServer) ListPostFeeds(ctx context.Context, req *newfeedProtoV1.ListPostFeedsRequest) (*newfeedProtoV1.ListPostFeedsResponse, error) {
	_, resp, err := s.listPostFeeds.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*newfeedProtoV1.ListPostFeedsResponse), nil
}

func (s *gRPCNewfeedServer) ListStoryFeeds(ctx context.Context, req *newfeedProtoV1.ListStoryFeedsRequest) (*newfeedProtoV1.ListStoryFeedsResponse, error) {
	_, resp, err := s.listStoryFeeds.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*newfeedProtoV1.ListStoryFeedsResponse), nil
}
