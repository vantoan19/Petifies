package relationshipclient

import (
	"context"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"

	relationshipProtoV1 "github.com/vantoan19/Petifies/proto/relationship-service/v1"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/models"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/translator"
)

var logger = logging.New("RelationshipClient")

type relationshipClient struct {
	listFollowers endpoint.Endpoint
}

type RelationshipClient interface {
	ListFollowers(ctx context.Context, userID uuid.UUID) (*models.ListFollowersResp, error)
}

func New(conn *grpc.ClientConn) RelationshipClient {
	return &relationshipClient{
		listFollowers: grpctransport.NewClient(
			conn,
			"RelationshipService",
			"ListFollowers",
			translator.EncodeListFollowersRequest,
			translator.DecodeListFollowersResponse,
			relationshipProtoV1.ListFollowersResponse{},
		).Endpoint(),
	}
}

func (c *relationshipClient) ListFollowers(ctx context.Context, userID uuid.UUID) (*models.ListFollowersResp, error) {
	logger.Info("Start ListFollowers")

	req := &models.ListFollowersReq{
		UserID: userID,
	}
	resp, err := c.listFollowers(ctx, req)
	if err != nil {
		logger.ErrorData("Finished ListFollowers: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished ListFollowers: SUCCESSFUL")
	return resp.(*models.ListFollowersResp), nil
}
