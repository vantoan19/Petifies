package translators

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	postProtoV1 "github.com/vantoan19/Petifies/proto/post-service/v1"
	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

func EncodeCreatePostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.CreatePostReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.CreatePostRequest{
		AuthorId: req.AuthorID.String(),
		Content:  req.TextContent,
		Images: utils.Map2(req.Images, func(i models.Image) *commonProto.Image {
			return &commonProto.Image{Uri: i.URL, Description: i.Description}
		}),
		Videos: utils.Map2(req.Videos, func(v models.Video) *commonProto.Video {
			return &commonProto.Video{Uri: v.URL, Description: v.Description}
		}),
	}, nil
}

func EncodeCreatePostResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.Post)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &commonProto.Post{
		Id:       resp.ID.String(),
		AuthorId: resp.AuthorID.String(),
		Content:  resp.Content,
		Images: utils.Map2(resp.Images, func(i models.Image) *commonProto.Image {
			return &commonProto.Image{
				Uri:         i.URL,
				Description: i.Description,
			}
		}),
		Videos: utils.Map2(resp.Videos, func(v models.Video) *commonProto.Video {
			return &commonProto.Video{
				Uri:         v.URL,
				Description: v.Description,
			}
		}),
		Loves: utils.Map2(resp.Loves, func(l models.Love) *commonProto.Love {
			return &commonProto.Love{
				Id:        l.ID.String(),
				PostId:    l.PostID.String(),
				CommentId: l.CommentID.String(),
				AuthorId:  l.AuthorID.String(),
				CreatedAt: timestamppb.New(l.CreatedAt),
			}
		}),
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}

func EncodeCreateCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.CreateCommentReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.CreateCommentRequest{
		PostId:       req.PostID.String(),
		AuthorId:     req.AuthorID.String(),
		ParentId:     req.ParentID.String(),
		IsPostParent: req.IsPostParent,
		Content:      req.Content,
		Image: &commonProto.Image{
			Uri:         req.ImageContent.URL,
			Description: req.ImageContent.Description,
		},
		Video: &commonProto.Video{
			Uri:         req.VideoContent.URL,
			Description: req.VideoContent.Description,
		},
	}, nil
}

func EncodeCreateCommentResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.Comment)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &commonProto.Comment{
		Id:           resp.ID.String(),
		PostId:       resp.PostID.String(),
		AuthorId:     resp.AuthorID.String(),
		ParentId:     resp.ParentID.String(),
		IsPostParent: resp.IsPostParent,
		Content:      resp.Content,
		Image: &commonProto.Image{
			Uri:         resp.Image.URL,
			Description: resp.Image.Description,
		},
		Video: &commonProto.Video{
			Uri:         resp.Video.URL,
			Description: resp.Video.Description,
		},
		Loves: utils.Map2(resp.Loves, func(l models.Love) *commonProto.Love {
			return &commonProto.Love{
				Id:        l.ID.String(),
				PostId:    l.PostID.String(),
				CommentId: l.CommentID.String(),
				AuthorId:  l.AuthorID.String(),
				CreatedAt: timestamppb.New(l.CreatedAt),
			}
		}),
		SubcommentCount: int32(resp.SubcommentCount),
		CreatedAt:       timestamppb.New(resp.CreatedAt),
		UpdatedAt:       timestamppb.New(resp.UpdatedAt),
	}, nil
}
