package translators

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
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

func EncodePostResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.Post)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return encodePostModel(resp), nil
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

func EncodeCommentResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.Comment)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return encodeCommentModel(resp), nil
}

func EncodeLoveReactRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.LoveReactReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &commonProto.LoveReactRequest{
		TargetId:     req.TargetID.String(),
		AuthorId:     req.AuthorID.String(),
		IsTargetPost: req.IsTargetPost,
	}, nil
}

func EncodeLoveResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.Love)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return encodeLoveModel(resp), nil
}

func EncodeEditPostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.EditPostReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.EditPostRequest{
		Id:      req.ID.String(),
		Content: req.Content,
		Images: utils.Map2(req.Images, func(i models.Image) *commonProto.Image {
			return &commonProto.Image{
				Uri:         i.URL,
				Description: i.Description,
			}
		}),
		Videos: utils.Map2(req.Videos, func(v models.Video) *commonProto.Video {
			return &commonProto.Video{
				Uri:         v.URL,
				Description: v.Description,
			}
		}),
	}, nil
}

func EncodeEditCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.EditCommentReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.EditCommentRequest{
		Id:      req.ID.String(),
		Content: req.Content,
		Image: &commonProto.Image{
			Uri:         req.Image.URL,
			Description: req.Image.Description,
		},
		Video: &commonProto.Video{
			Uri:         req.Video.URL,
			Description: req.Video.Description,
		},
	}, nil
}

func EncodeListCommentsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListCommentsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.ListCommentsRequest{
		CommentIds: utils.Map2(req.CommentIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}

func EncodeListCommentsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListCommentsResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &postProtoV1.ListCommentsResponse{
		Comments: utils.Map2(resp.Comments, func(c *models.Comment) *commonProto.Comment { return encodeCommentModel(c) }),
	}, nil
}

func EncodeListPostsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListPostsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.ListPostsRequest{
		PostIds: utils.Map2(req.PostIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}

func EncodeListPostsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListPostsResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &postProtoV1.ListPostsResponse{
		Posts: utils.Map2(resp.Posts, func(p *models.Post) *commonProto.Post { return encodePostModel(p) }),
	}, nil
}

func encodePostModel(post *models.Post) *commonProto.Post {
	return &commonProto.Post{
		Id:       post.ID.String(),
		AuthorId: post.AuthorID.String(),
		Content:  post.Content,
		Images: utils.Map2(post.Images, func(i models.Image) *commonProto.Image {
			return &commonProto.Image{
				Uri:         i.URL,
				Description: i.Description,
			}
		}),
		Videos: utils.Map2(post.Videos, func(v models.Video) *commonProto.Video {
			return &commonProto.Video{
				Uri:         v.URL,
				Description: v.Description,
			}
		}),
		Loves: utils.Map2(post.Loves, func(l models.Love) *commonProto.Love {
			return &commonProto.Love{
				Id:        l.ID.String(),
				PostId:    l.PostID.String(),
				CommentId: l.CommentID.String(),
				AuthorId:  l.AuthorID.String(),
				CreatedAt: timestamppb.New(l.CreatedAt),
			}
		}),
		CreatedAt: timestamppb.New(post.CreatedAt),
		UpdatedAt: timestamppb.New(post.UpdatedAt),
	}
}

func encodeCommentModel(comment *models.Comment) *commonProto.Comment {
	return &commonProto.Comment{
		Id:           comment.ID.String(),
		PostId:       comment.PostID.String(),
		AuthorId:     comment.AuthorID.String(),
		ParentId:     comment.ParentID.String(),
		IsPostParent: comment.IsPostParent,
		Content:      comment.Content,
		Image: &commonProto.Image{
			Uri:         comment.Image.URL,
			Description: comment.Image.Description,
		},
		Video: &commonProto.Video{
			Uri:         comment.Video.URL,
			Description: comment.Video.Description,
		},
		Loves: utils.Map2(comment.Loves, func(l models.Love) *commonProto.Love {
			return &commonProto.Love{
				Id:        l.ID.String(),
				PostId:    l.PostID.String(),
				CommentId: l.CommentID.String(),
				AuthorId:  l.AuthorID.String(),
				CreatedAt: timestamppb.New(l.CreatedAt),
			}
		}),
		SubcommentCount: int32(comment.SubcommentCount),
		CreatedAt:       timestamppb.New(comment.CreatedAt),
		UpdatedAt:       timestamppb.New(comment.UpdatedAt),
	}
}

func encodeLoveModel(love *models.Love) *commonProto.Love {
	return &commonProto.Love{
		Id:        love.ID.String(),
		PostId:    love.PostID.String(),
		CommentId: love.CommentID.String(),
		AuthorId:  love.AuthorID.String(),
		CreatedAt: timestamppb.New(love.CreatedAt),
	}
}
