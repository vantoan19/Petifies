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
	addRelationship    endpoint.Endpoint
	removeRelationship endpoint.Endpoint
	listFollowers      endpoint.Endpoint
	listFollowings     endpoint.Endpoint
}

type RelationshipClient interface {
	AddRelationship(ctx context.Context, req *models.AddRelationshipReq) (*models.AddRelationshipResp, error)
	RemoveRelationship(ctx context.Context, req *models.RemoveRelationshipReq) (*models.RemoveRelationshipResp, error)
	ListFollowers(ctx context.Context, userID uuid.UUID) (*models.ListFollowersResp, error)
	ListFollowings(ctx context.Context, userID uuid.UUID) (*models.ListFollowingsResp, error)
}

func New(conn *grpc.ClientConn) RelationshipClient {
	return &relationshipClient{
		addRelationship: grpctransport.NewClient(
			conn,
			"RelationshipService",
			"AddRelationship",
			translator.EncodeAddRelationshipRequest,
			translator.DecodeAddRelationshipResponse,
			relationshipProtoV1.AddRelationshipResponse{},
		).Endpoint(),
		removeRelationship: grpctransport.NewClient(
			conn,
			"RelationshipService",
			"RemoveRelationship",
			translator.EncodeRemoveRelationshipRequest,
			translator.DecodeRemoveRelationshipResponse,
			relationshipProtoV1.RemoveRelationshipResponse{},
		).Endpoint(),
		listFollowers: grpctransport.NewClient(
			conn,
			"RelationshipService",
			"ListFollowers",
			translator.EncodeListFollowersRequest,
			translator.DecodeListFollowersResponse,
			relationshipProtoV1.ListFollowersResponse{},
		).Endpoint(),
		listFollowings: grpctransport.NewClient(
			conn,
			"RelationshipService",
			"ListFollowings",
			translator.EncodeListFollowingsRequest,
			translator.DecodeListFollowingsResponse,
			relationshipProtoV1.ListFollowingsResponse{},
		).Endpoint(),
	}
}

func (c *relationshipClient) AddRelationship(ctx context.Context, req *models.AddRelationshipReq) (*models.AddRelationshipResp, error) {
	logger.Info("Start AddRelationship")

	resp, err := c.addRelationship(ctx, req)
	if err != nil {
		logger.ErrorData("Finished AddRelationship: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished AddRelationship: SUCCESSFUL")
	return resp.(*models.AddRelationshipResp), nil
}

func (c *relationshipClient) RemoveRelationship(ctx context.Context, req *models.RemoveRelationshipReq) (*models.RemoveRelationshipResp, error) {
	logger.Info("Start RemoveRelationship")

	resp, err := c.removeRelationship(ctx, req)
	if err != nil {
		logger.ErrorData("Finished RemoveRelationship: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished RemoveRelationship: SUCCESSFUL")
	return resp.(*models.RemoveRelationshipResp), nil
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

func (c *relationshipClient) ListFollowings(ctx context.Context, userID uuid.UUID) (*models.ListFollowingsResp, error) {
	logger.Info("Start ListFollowings")

	req := &models.ListFollowingsReq{
		UserID: userID,
	}
	resp, err := c.listFollowings(ctx, req)
	if err != nil {
		logger.ErrorData("Finished ListFollowings: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished ListFollowings: SUCCESSFUL")
	return resp.(*models.ListFollowingsResp), nil
}
