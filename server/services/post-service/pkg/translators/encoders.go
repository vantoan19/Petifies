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
		AuthorId:   req.AuthorID.String(),
		Visibility: req.Visibility,
		Actitivty:  req.Activity,
		Content:    req.TextContent,
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

	return &postProtoV1.LoveReactRequest{
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
		Id:         req.ID.String(),
		Content:    req.Content,
		Visibility: req.Visibility,
		Activity:   req.Activity,
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

func EncodeGetLoveCountRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetLoveCountReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.GetLoveCountRequest{
		TargetId:     req.TargetID.String(),
		IsPostTarget: req.IsPostParent,
	}, nil
}

func EncodeGetLoveCountResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.GetLoveCountResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &postProtoV1.GetLoveCountReponse{
		Count: int32(resp.Count),
	}, nil
}

func EncodeGetCommentCountRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetCommentCountReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.GetCommentCountRequest{
		ParentId:     req.ParentID.String(),
		IsPostParent: req.IsPostParent,
	}, nil
}

func EncodeGetCommentCountResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.GetCommentCountResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &postProtoV1.GetCommentCountReponse{
		Count: int32(resp.Count),
	}, nil
}

func EncodeGetPostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetPostReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.GetPostRequest{
		PostId: req.PostID.String(),
	}, nil
}

func EncodeGetCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetCommentReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.GetCommentRequest{
		CommentId: req.CommentID.String(),
	}, nil
}

func EncodeRemoveLoveReactRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.RemoveLoveReactReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.RemoveLoveReactRequest{
		TargetId:     req.TargetID.String(),
		AuthorId:     req.AuthorID.String(),
		IsTargetPost: req.IsTargetPost,
	}, nil
}

func EncodeRemoveLoveReactResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(*models.RemoveLoveReactResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &postProtoV1.RemoveLoveReactResponse{}, nil
}

func EncodeGetLoveRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetLoveReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.GetLoveRequest{
		AuthorId: req.AuthorID.String(),
		TargetId: req.TargetID.String(),
	}, nil
}

func EncodeListCommentIDsByParentIDRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListCommentIDsByParentIDReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.ListCommentIDsByParentIDRequest{
		ParentId:       req.ParentID.String(),
		PageSize:       int32(req.PageSize),
		AfterCommentId: req.AfterCommentID.String(),
	}, nil
}

func EncodeListCommentIDsByParentIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListCommentIDsByParentIDResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &postProtoV1.ListCommentIDsByParentIDResponse{
		CommentIds: utils.Map2(resp.CommentIDs, func(c uuid.UUID) string { return c.String() }),
	}, nil
}

func EncodeListCommentAncestorsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListCommentAncestorsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &postProtoV1.ListCommentAncestorsRequest{
		CommentId: req.CommentID.String(),
	}, nil
}

func EncodeListCommentAncestorsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListCommentAncestorsResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &postProtoV1.ListCommentAncestorsResponse{
		AncestorComments: utils.Map2(resp.AncestorComments, func(c *models.Comment) *commonProto.Comment { return encodeCommentModel(c) }),
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
		LoveCount:    int32(post.LoveCount),
		CommentCount: int32(post.CommentCount),
		Visibility:   post.Visibility,
		Activity:     post.Activity,
		CreatedAt:    timestamppb.New(post.CreatedAt),
		UpdatedAt:    timestamppb.New(post.UpdatedAt),
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
		LoveCount:       int32(comment.LoveCount),
		SubcommentCount: int32(comment.SubcommentCount),
		CreatedAt:       timestamppb.New(comment.CreatedAt),
		UpdatedAt:       timestamppb.New(comment.UpdatedAt),
	}
}

func encodeLoveModel(love *models.Love) *commonProto.Love {
	return &commonProto.Love{
		Id:           love.ID.String(),
		TargetId:     love.TargetID.String(),
		IsPostTarget: love.IsPostTarget,
		AuthorId:     love.AuthorID.String(),
		CreatedAt:    timestamppb.New(love.CreatedAt),
	}
}
