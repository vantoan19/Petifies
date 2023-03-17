package translators

import (
	"context"

	"github.com/google/uuid"
	newfeedProtoV1 "github.com/vantoan19/Petifies/proto/newfeed-service/v1"
	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/pkg/models"
)

func EncodeListPostFeedsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListPostFeedsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &newfeedProtoV1.ListPostFeedsRequest{
		UserId: req.UserID.String(),
	}, nil
}

func EncodeListPostFeedsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListPostFeedsResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &newfeedProtoV1.ListPostFeedsResponse{
		PostIds: utils.Map2(resp.PostIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}

func EncodeListStoryFeedsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListStoryFeedsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &newfeedProtoV1.ListStoryFeedsRequest{
		UserId: req.UserID.String(),
	}, nil
}

func EncodeListStoryFeedsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListStoryFeedsResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &newfeedProtoV1.ListStoryFeedsResponse{
		StoryIds: utils.Map2(resp.StoryIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}
