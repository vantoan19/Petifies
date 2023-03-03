package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	commentservice "github.com/vantoan19/Petifies/server/services/post-service/internal/application/services/comment"
	postservice "github.com/vantoan19/Petifies/server/services/post-service/internal/application/services/post"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type PostEndpoints struct {
	CreatePost    endpoint.Endpoint
	CreateComment endpoint.Endpoint
	LoveReact     endpoint.Endpoint
	EditPost      endpoint.Endpoint
	EditComment   endpoint.Endpoint
	ListComments  endpoint.Endpoint
	ListPosts     endpoint.Endpoint
}

func NewPostEndpoints(ps postservice.PostService, cs commentservice.CommentService) PostEndpoints {
	return PostEndpoints{
		CreatePost:    makeCreatePostEndpoint(ps),
		CreateComment: makeCreateCommentEndpoint(cs),
		LoveReact:     makeLoveReactEndpoint(ps, cs),
		EditPost:      makeEditPostEndpoint(ps),
		EditComment:   makeEditCommentEndpoint(cs),
		ListComments:  makeListCommentsEndpoint(cs),
		ListPosts:     makeListPostsEndpoint(ps),
	}
}

// ================ Endpoint Makers ==================

func makeCreatePostEndpoint(ps postservice.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CreatePostReq)
		result, err := ps.CreatePost(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapPostAggregateToPostModel(result), nil
	}
}

func makeCreateCommentEndpoint(cs commentservice.CommentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CreateCommentReq)
		result, err := cs.CreateComment(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapCommentAggregateToCommentModel(result), nil
	}
}

func makeLoveReactEndpoint(ps postservice.PostService, cs commentservice.CommentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.LoveReactReq)

		var result *entities.Love
		if req.IsTargetPost {
			result, err = ps.LoveReactPost(ctx, req)
			if err != nil {
				return nil, err
			}
		} else {
			result, err = cs.LoveReactComment(ctx, req)
			if err != nil {
				return nil, err
			}
		}

		return mapLoveEntityToLoveModel(result), nil
	}
}

func makeEditPostEndpoint(ps postservice.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.EditPostReq)
		result, err := ps.EditPost(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapPostAggregateToPostModel(result), nil
	}
}

func makeEditCommentEndpoint(cs commentservice.CommentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.EditCommentReq)
		result, err := cs.EditComment(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapCommentAggregateToCommentModel(result), nil
	}
}

func makeListPostsEndpoint(ps postservice.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListPostsReq)
		results, err := ps.ListPosts(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.ListPostsResp{
			Posts: utils.Map2(results, func(p *postaggre.Post) *models.Post { return mapPostAggregateToPostModel(p) }),
		}, nil
	}
}

func makeListCommentsEndpoint(cs commentservice.CommentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListCommentsReq)
		results, err := cs.ListComments(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.ListCommentsResp{
			Comments: utils.Map2(results, func(c *commentaggre.Comment) *models.Comment { return mapCommentAggregateToCommentModel(c) }),
		}, nil
	}
}

// ================ Mappers ======================

func mapPostAggregateToPostModel(post *postaggre.Post) *models.Post {
	images := utils.Map2(post.GetImages(), func(i valueobjects.ImageContent) models.Image {
		return models.Image{URL: i.URL(), Description: i.Description()}
	})
	videos := utils.Map2(post.GetVideos(), func(v valueobjects.VideoContent) models.Video {
		return models.Video{URL: v.URL(), Description: v.Description()}
	})
	loves := utils.Map2(post.GetLoves(), func(l entities.Love) models.Love {
		return *mapLoveEntityToLoveModel(&l)
	})

	return &models.Post{
		ID:        post.GetPostID(),
		AuthorID:  post.GetAuthorID(),
		Content:   post.GetPostTextContent().Content(),
		Images:    images,
		Videos:    videos,
		Loves:     loves,
		CreatedAt: post.GetCreatedAt(),
		UpdatedAt: post.GetUpdatedAt(),
	}
}

func mapCommentAggregateToCommentModel(comment *commentaggre.Comment) *models.Comment {
	loves := utils.Map2(comment.GetLoves(), func(l entities.Love) models.Love {
		return *mapLoveEntityToLoveModel(&l)
	})

	return &models.Comment{
		ID:           comment.GetID(),
		PostID:       comment.GetPostID(),
		AuthorID:     comment.GetAuthorID(),
		ParentID:     comment.GetParentID(),
		IsPostParent: comment.GetIsPostParent(),
		Content:      comment.GetContent().Content(),
		Image: models.Image{
			URL:         comment.GetImageContent().URL(),
			Description: comment.GetImageContent().Description(),
		},
		Video: models.Video{
			URL:         comment.GetVideoContent().URL(),
			Description: comment.GetVideoContent().Description(),
		},
		Loves:           loves,
		SubcommentCount: comment.CountSubcomments(),
		CreatedAt:       comment.GetCreatedAt(),
		UpdatedAt:       comment.GetUpdatedAt(),
	}
}

func mapLoveEntityToLoveModel(love *entities.Love) *models.Love {
	return &models.Love{
		ID:        love.ID,
		PostID:    love.PostID,
		CommentID: love.CommentID,
		AuthorID:  love.AuthorID,
		CreatedAt: love.CreatedAt,
	}
}
