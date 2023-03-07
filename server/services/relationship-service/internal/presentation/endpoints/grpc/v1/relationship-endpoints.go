package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	relationshipservice "github.com/vantoan19/Petifies/server/services/relationship-service/internal/application/services/relationship"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/models"
)

type RelationshipEndpoints struct {
	AddRelationship    endpoint.Endpoint
	RemoveRelationship endpoint.Endpoint
	ListFollowers      endpoint.Endpoint
	ListFollowings     endpoint.Endpoint
}

func NewRelationshipEndpoints(rs relationshipservice.RelationshipService) RelationshipEndpoints {
	return RelationshipEndpoints{
		AddRelationship:    makeAddRelationshipEndpoint(rs),
		RemoveRelationship: makeRemoveRelationshipEndpoint(rs),
		ListFollowers:      makeListFollowersEndpoint(rs),
		ListFollowings:     makeListFollowingsEndpoint(rs),
	}
}

// ============ Endpoint Makers ============

func makeAddRelationshipEndpoint(rs relationshipservice.RelationshipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.AddRelationshipReq)
		_, err = rs.AddRelationship(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.AddRelationshipResp{
			Message: "Add relationship successfully",
		}, nil
	}
}

func makeRemoveRelationshipEndpoint(rs relationshipservice.RelationshipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.RemoveRelationshipReq)
		_, err = rs.RemoveRelationship(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.RemoveRelationshipResp{
			Message: "Remove relation successfully",
		}, nil
	}
}

func makeListFollowersEndpoint(rs relationshipservice.RelationshipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListFollowersReq)
		result, err := rs.ListFollowers(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.ListFollowersResp{
			FollowerIDs: result,
		}, nil
	}
}

func makeListFollowingsEndpoint(rs relationshipservice.RelationshipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListFollowingsReq)
		result, err := rs.ListFollowings(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.ListFollowingsResp{
			FollowingIDs: result,
		}, nil
	}
}
