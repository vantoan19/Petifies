package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	commentservice "github.com/vantoan19/Petifies/server/services/post-service/internal/application/services/comment"
	postservice "github.com/vantoan19/Petifies/server/services/post-service/internal/application/services/post"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	loveaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/love"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type PostEndpoints struct {
	CreatePost      endpoint.Endpoint
	CreateComment   endpoint.Endpoint
	LoveReact       endpoint.Endpoint
	EditPost        endpoint.Endpoint
	EditComment     endpoint.Endpoint
	ListComments    endpoint.Endpoint
	ListPosts       endpoint.Endpoint
	GetLoveCount    endpoint.Endpoint
	GetCommentCount endpoint.Endpoint
	GetPost         endpoint.Endpoint
	GetComment      endpoint.Endpoint
}

func NewPostEndpoints(ps postservice.PostService, cs commentservice.CommentService) PostEndpoints {
	return PostEndpoints{
		CreatePost:      makeCreatePostEndpoint(ps),
		CreateComment:   makeCreateCommentEndpoint(cs),
		LoveReact:       makeLoveReactEndpoint(ps, cs),
		EditPost:        makeEditPostEndpoint(ps),
		EditComment:     makeEditCommentEndpoint(cs),
		ListComments:    makeListCommentsEndpoint(cs),
		ListPosts:       makeListPostsEndpoint(ps),
		GetLoveCount:    makeGetLoveCountEndpoint(ps, cs),
		GetCommentCount: makeGetCommentCountEndpoint(ps, cs),
		GetComment:      makeGetCommentEndpoint(cs),
		GetPost:         makeGetPostEndpoint(ps),
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

		var result *loveaggre.Love
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
		var postModels []*models.Post
		for _, p := range results {
			postModels = append(postModels, mapPostAggregateToPostModel(p))
		}

		return &models.ListPostsResp{
			Posts: postModels,
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
		var commentModels []*models.Comment
		for _, p := range results {
			commentModels = append(commentModels, mapCommentAggregateToCommentModel(p))
		}

		return &models.ListCommentsResp{
			Comments: commentModels,
		}, nil
	}
}

func makeGetLoveCountEndpoint(ps postservice.PostService, cs commentservice.CommentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.GetLoveCountReq)
		var result int

		if req.IsPostParent {
			result, err = ps.GetLoveCount(ctx, req.TargetID)
		} else {
			result, err = cs.GetLoveCount(ctx, req.TargetID)
		}

		return &models.GetLoveCountResp{
			Count: result,
		}, err
	}
}

func makeGetCommentCountEndpoint(ps postservice.PostService, cs commentservice.CommentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.GetCommentCountReq)
		var result int

		if req.IsPostParent {
			result, err = ps.GetCommentCount(ctx, req.ParentID)
		} else {
			result, err = cs.GetSubcommentCount(ctx, req.ParentID)
		}

		return &models.GetCommentCountResp{
			Count: result,
		}, err
	}
}

func makeGetPostEndpoint(ps postservice.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.GetPostReq)
		result, err := ps.GetPost(ctx, req.PostID)

		return mapPostAggregateToPostModel(result), err
	}
}

func makeGetCommentEndpoint(cs commentservice.CommentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.GetCommentReq)
		result, err := cs.GetComment(ctx, req.CommentID)

		return mapCommentAggregateToCommentModel(result), err
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

	return &models.Post{
		ID:         post.GetPostID(),
		AuthorID:   post.GetAuthorID(),
		Content:    post.GetPostTextContent().Content(),
		Images:     images,
		Videos:     videos,
		Visibility: string(post.GetVisibility()),
		Activity:   post.GetActivity(),
		CreatedAt:  post.GetCreatedAt(),
		UpdatedAt:  post.GetUpdatedAt(),
	}
}

func mapCommentAggregateToCommentModel(comment *commentaggre.Comment) *models.Comment {
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
		CreatedAt: comment.GetCreatedAt(),
		UpdatedAt: comment.GetUpdatedAt(),
	}
}

func mapLoveEntityToLoveModel(love *loveaggre.Love) *models.Love {
	return &models.Love{
		ID:           love.GetID(),
		TargetID:     love.GetTargetID(),
		IsPostTarget: love.GetIsPostTarget(),
		AuthorID:     love.GetAuthorID(),
		CreatedAt:    love.GetCreatedAt(),
	}
}
