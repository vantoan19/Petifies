package translators

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	postProtoV1 "github.com/vantoan19/Petifies/proto/post-service/v1"
	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var (
	MustBeProtoReqErr     = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr    = status.Error(codes.InvalidArgument, "must be proto response")
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
)

func DecodeCreatePostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.CreatePostRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	authorId, err := uuid.Parse(req.GetAuthorId())
	if err != nil {
		return nil, err
	}

	return &models.CreatePostReq{
		AuthorID:    authorId,
		TextContent: req.GetContent(),
		Images:      utils.Map2(req.GetImages(), func(i *commonProto.Image) models.Image { return models.Image{URL: i.Uri, Description: i.Description} }),
		Videos:      utils.Map2(req.GetVideos(), func(v *commonProto.Video) models.Video { return models.Video{URL: v.Uri, Description: v.Description} }),
	}, nil
}

func DecodeCreatePostResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.Post)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	id, err := uuid.Parse(resp.GetId())
	if err != nil {
		return nil, err
	}
	authorID, err := uuid.Parse(resp.GetAuthorId())
	if err != nil {
		return nil, err
	}
	loves := make([]models.Love, 0)
	for _, l := range resp.GetLoves() {
		loveId, err := uuid.Parse(l.Id)
		if err != nil {
			return nil, err
		}
		postId := uuid.Nil
		if l.PostId != "" {
			postId, err = uuid.Parse(l.PostId)
			if err != nil {
				return nil, err
			}
		}
		commentId := uuid.Nil
		if l.CommentId != "" {
			commentId, err = uuid.Parse(l.CommentId)
			if err != nil {
				return nil, err
			}
		}
		authorId, err := uuid.Parse(l.AuthorId)
		if err != nil {
			return nil, err
		}
		loves = append(loves, models.Love{
			ID:        loveId,
			PostID:    postId,
			CommentID: commentId,
			AuthorID:  authorId,
			CreatedAt: l.CreatedAt.AsTime(),
		})
	}

	return &models.Post{
		ID:        id,
		AuthorID:  authorID,
		Content:   resp.GetContent(),
		Images:    utils.Map2(resp.GetImages(), func(i *commonProto.Image) models.Image { return models.Image{URL: i.Uri, Description: i.Description} }),
		Videos:    utils.Map2(resp.GetVideos(), func(v *commonProto.Video) models.Video { return models.Video{URL: v.Uri, Description: v.Description} }),
		Loves:     loves,
		CreatedAt: resp.GetCreatedAt().AsTime(),
		UpdatedAt: resp.GetUpdatedAt().AsTime(),
	}, nil
}

func DecodeCreateCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*postProtoV1.CreateCommentRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	postId, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}
	authorId, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return nil, err
	}
	parentId, err := uuid.Parse(req.ParentId)
	if err != nil {
		return nil, err
	}

	return &models.CreateCommentReq{
		PostID:       postId,
		AuthorID:     authorId,
		ParentID:     parentId,
		IsPostParent: req.IsPostParent,
		Content:      req.Content,
		ImageContent: models.Image{
			URL:         req.Image.Uri,
			Description: req.Image.Description,
		},
		VideoContent: models.Video{
			URL:         req.Video.Uri,
			Description: req.Video.Description,
		},
	}, nil
}

func DecodeCreateCommentResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*commonProto.Comment)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	id, err := uuid.Parse(resp.Id)
	if err != nil {
		return nil, err
	}
	postId, err := uuid.Parse(resp.PostId)
	if err != nil {
		return nil, err
	}
	authorId, err := uuid.Parse(resp.AuthorId)
	if err != nil {
		return nil, err
	}
	parentId, err := uuid.Parse(resp.ParentId)
	if err != nil {
		return nil, err
	}
	loves := make([]models.Love, 0)
	for _, l := range resp.GetLoves() {
		loveId, err := uuid.Parse(l.Id)
		if err != nil {
			return nil, err
		}
		postId := uuid.Nil
		if l.PostId != "" {
			postId, err = uuid.Parse(l.PostId)
			if err != nil {
				return nil, err
			}
		}
		commentId := uuid.Nil
		if l.CommentId != "" {
			commentId, err = uuid.Parse(l.CommentId)
			if err != nil {
				return nil, err
			}
		}
		authorId, err := uuid.Parse(l.AuthorId)
		if err != nil {
			return nil, err
		}
		loves = append(loves, models.Love{
			ID:        loveId,
			PostID:    postId,
			CommentID: commentId,
			AuthorID:  authorId,
			CreatedAt: l.CreatedAt.AsTime(),
		})
	}

	return &models.Comment{
		ID:           id,
		PostID:       postId,
		AuthorID:     authorId,
		ParentID:     parentId,
		IsPostParent: resp.GetIsPostParent(),
		Content:      resp.GetContent(),
		Image: models.Image{
			URL:         resp.GetImage().Uri,
			Description: resp.GetImage().Description,
		},
		Video: models.Video{
			URL:         resp.GetVideo().Uri,
			Description: resp.GetVideo().Description,
		},
		Loves:           loves,
		SubcommentCount: int(resp.GetSubcommentCount()),
		CreatedAt:       resp.CreatedAt.AsTime(),
		UpdatedAt:       resp.UpdatedAt.AsTime(),
	}, nil
}
