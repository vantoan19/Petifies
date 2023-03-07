package translator

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	relationshipProtoV1 "github.com/vantoan19/Petifies/proto/relationship-service/v1"
	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/models"
)

var (
	MustBeProtoReqErr     = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr    = status.Error(codes.InvalidArgument, "must be proto response")
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
)

func EncodeAddRelationshipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.AddRelationshipReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &relationshipProtoV1.AddRelationshipRequest{
		FromUserId:       req.FromUserID.String(),
		ToUserId:         req.ToUserID.String(),
		RelationshipType: req.RelationshipType,
	}, nil
}

func EncodeAddRelationshipResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.AddRelationshipResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &relationshipProtoV1.AddRelationshipResponse{
		Message: resp.Message,
	}, nil
}

func EncodeRemoveRelationshipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.RemoveRelationshipReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &relationshipProtoV1.RemoveRelationshipRequest{
		FromUserId:       req.FromUserID.String(),
		ToUserId:         req.ToUserID.String(),
		RelationshipType: req.RelationshipType,
	}, nil
}

func EncodeRemoveRelationshipResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.RemoveRelationshipResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &relationshipProtoV1.RemoveRelationshipResponse{
		Message: resp.Message,
	}, nil
}

func EncodeListFollowersRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListFollowersReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &relationshipProtoV1.ListFollowersRequest{
		UserId: req.UserID.String(),
	}, nil
}

func EncodeListFollowersResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListFollowersResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &relationshipProtoV1.ListFollowersResponse{
		FollowerIds: utils.Map2(resp.FollowerIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}

func EncodeListFollowingsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListFollowingsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &relationshipProtoV1.ListFollowingsRequest{
		UserId: req.UserID.String(),
	}, nil
}

func EncodeListFollowingsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListFollowingsResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &relationshipProtoV1.ListFollowingsResponse{
		FollowingIds: utils.Map2(resp.FollowingIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}
