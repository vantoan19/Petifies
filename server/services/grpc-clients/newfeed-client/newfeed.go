package newfeedclient

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	newfeedProtoV1 "github.com/vantoan19/Petifies/proto/newfeed-service/v1"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/pkg/models"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/pkg/translators"
)

var logger = logging.New("NewfeedClient")

const newfeedService = "newfeed_service.v1.NewfeedService"

type newfeedClient struct {
	listPostFeeds endpoint.Endpoint
}

type NewfeedClient interface {
	ListPostFeeds(ctx context.Context, req *models.ListPostFeedsReq) (*models.ListPostFeedsResp, error)
}

func New(conn *grpc.ClientConn) NewfeedClient {
	return &newfeedClient{
		listPostFeeds: grpctransport.NewClient(
			conn,
			newfeedService,
			"ListPostFeeds",
			translators.EncodeListPostFeedsRequest,
			translators.DecodeListPostFeedsResponse,
			newfeedProtoV1.ListPostFeedsResponse{},
		).Endpoint(),
	}
}

func (c *newfeedClient) ListPostFeeds(ctx context.Context, req *models.ListPostFeedsReq) (*models.ListPostFeedsResp, error) {
	logger.Info("Start ListPostFeeds")

	resp, err := c.listPostFeeds(ctx, req)
	if err != nil {
		logger.ErrorData("Finished ListPostFeeds: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished ListPostFeeds: SUCCESSFUL")
	return resp.(*models.ListPostFeedsResp), nil
}
