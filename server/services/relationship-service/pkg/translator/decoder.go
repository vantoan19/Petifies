package translator

import (
	"context"

	"github.com/google/uuid"

	relationshipProtoV1 "github.com/vantoan19/Petifies/proto/relationship-service/v1"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/models"
)

func DecodeAddRelationshipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*relationshipProtoV1.AddRelationshipRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	fromUserID, err := uuid.Parse(req.FromUserId)
	if err != nil {
		return nil, err
	}
	toUserID, err := uuid.Parse(req.ToUserId)
	if err != nil {
		return nil, err
	}

	return &models.AddRelationshipReq{
		FromUserID:       fromUserID,
		ToUserID:         toUserID,
		RelationshipType: req.RelationshipType,
	}, nil
}

func DecodeAddRelationshipResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*relationshipProtoV1.AddRelationshipResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &models.AddRelationshipResp{
		Message: resp.Message,
	}, nil
}

func DecodeRemoveRelationshipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*relationshipProtoV1.RemoveRelationshipRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	fromUserID, err := uuid.Parse(req.FromUserId)
	if err != nil {
		return nil, err
	}
	toUserID, err := uuid.Parse(req.ToUserId)
	if err != nil {
		return nil, err
	}

	return &models.RemoveRelationshipReq{
		FromUserID:       fromUserID,
		ToUserID:         toUserID,
		RelationshipType: req.RelationshipType,
	}, nil
}

func DecodeRemoveRelationshipResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*relationshipProtoV1.RemoveRelationshipResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &models.RemoveRelationshipResp{
		Message: resp.Message,
	}, nil
}

func DecodeListFollowersRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*relationshipProtoV1.ListFollowersRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	return &models.ListFollowersReq{
		UserID: userID,
	}, nil
}

func DecodeListFollowersResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*relationshipProtoV1.ListFollowersResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	followerIDs := make([]uuid.UUID, 0)
	for _, id := range resp.FollowerIds {
		id_, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		followerIDs = append(followerIDs, id_)
	}

	return &models.ListFollowersResp{
		FollowerIDs: followerIDs,
	}, nil
}

func DecodeListFollowingsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*relationshipProtoV1.ListFollowingsRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	return &models.ListFollowingsReq{
		UserID: userID,
	}, nil
}

func DecodeListFollowingsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*relationshipProtoV1.ListFollowingsResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	followingIDs := make([]uuid.UUID, 0)
	for _, id := range resp.FollowingIds {
		id_, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		followingIDs = append(followingIDs, id_)
	}

	return &models.ListFollowingsResp{
		FollowingIDs: followingIDs,
	}, nil
}
