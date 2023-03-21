package translators

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	newfeedProtoV1 "github.com/vantoan19/Petifies/proto/newfeed-service/v1"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/pkg/models"
)

var (
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
	MustBeProtoReqErr     = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr    = status.Error(codes.InvalidArgument, "must be proto response")
)

func DecodeListPostFeedsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*newfeedProtoV1.ListPostFeedsRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &models.ListPostFeedsReq{
		UserID:     id,
		PageSize:   int(req.PageSize),
		BeforeTime: req.BeforeTime.AsTime(),
	}, nil
}

func DecodeListPostFeedsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*newfeedProtoV1.ListPostFeedsResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	postIDs := make([]uuid.UUID, 0)
	for _, id := range resp.PostIds {
		id_, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		postIDs = append(postIDs, id_)
	}

	return &models.ListPostFeedsResp{
		PostIDs: postIDs,
	}, nil
}

func DecodeListStoryFeedsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*newfeedProtoV1.ListStoryFeedsRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &models.ListStoryFeedsReq{
		UserID:     id,
		PageSize:   int(req.PageSize),
		BeforeTime: req.BeforeTime.AsTime(),
	}, nil
}

func DecodeListStoryFeedsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*newfeedProtoV1.ListStoryFeedsResponse)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	storyIDs := make([]uuid.UUID, 0)
	for _, id := range resp.StoryIds {
		id_, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		storyIDs = append(storyIDs, id_)
	}

	return &models.ListStoryFeedsResp{
		StoryIDs: storyIDs,
	}, nil
}
