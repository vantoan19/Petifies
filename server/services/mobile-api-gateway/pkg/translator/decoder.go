package translator

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-gateway/v1"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
	postModels "github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var (
	MustBeProtoReqErr     = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr    = status.Error(codes.InvalidArgument, "must be proto response")
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
)

func DecodeUserCreatePostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserCreatePostRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.UserCreatePostReq{
		TextContent: req.Content,
		Images: utils.Map2(req.Images, func(i *commonProto.Image) postModels.Image {
			return postModels.Image{URL: i.Uri, Description: i.Description}
		}),
		Videos: utils.Map2(req.Videos, func(v *commonProto.Video) postModels.Video {
			return postModels.Video{URL: v.Uri, Description: v.Description}
		}),
	}, nil
}

func DecodeUserCreateCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserCreateCommentRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}
	parentID, err := uuid.Parse(req.ParentId)
	if err != nil {
		return nil, err
	}

	return &models.UserCreateCommentReq{
		PostID:       postID,
		ParentID:     parentID,
		IsParentPost: req.IsPostParent,
		Content:      req.Content,
		Image:        postModels.Image{URL: req.Image.Uri, Description: req.Image.Description},
		Video:        postModels.Video{URL: req.Video.Uri, Description: req.Video.Description},
	}, nil
}

func DecodeUserEditPostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserEditPostRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	postID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	return &models.UserEditPostReq{
		PostID:  postID,
		Content: req.Content,
		Images: utils.Map2(req.Images, func(i *commonProto.Image) postModels.Image {
			return postModels.Image{URL: i.Uri, Description: i.Description}
		}),
		Videos: utils.Map2(req.Videos, func(v *commonProto.Video) postModels.Video {
			return postModels.Video{URL: v.Uri, Description: v.Description}
		}),
	}, nil
}

func DecodeUserEditCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserEditCommentRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	commentID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	return &models.UserEditCommentReq{
		CommentID: commentID,
		Content:   req.Content,
		Image:     postModels.Image{URL: req.Image.Uri, Description: req.Image.Description},
		Video:     postModels.Video{URL: req.Video.Uri, Description: req.Video.Description},
	}, nil
}
