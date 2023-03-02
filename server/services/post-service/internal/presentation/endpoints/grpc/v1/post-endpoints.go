package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/application/services"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type PostEndpoints struct {
	CreatePost    endpoint.Endpoint
	CreateComment endpoint.Endpoint
}

func NewPostEndpoints(ps services.PostService, cs services.CommentService) PostEndpoints {
	return PostEndpoints{
		CreatePost:    makeCreatePostEndpoint(ps),
		CreateComment: makeCreateCommentEndpoint(cs),
	}
}

func makeCreatePostEndpoint(ps services.PostService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CreatePostReq)
		result, err := ps.CreatePost(ctx, req)
		if err != nil {
			return nil, err
		}

		images := utils.Map2(result.GetImages(), func(i valueobjects.ImageContent) models.Image {
			return models.Image{URL: i.URL(), Description: i.Description()}
		})
		videos := utils.Map2(result.GetVideos(), func(v valueobjects.VideoContent) models.Video {
			return models.Video{URL: v.URL(), Description: v.Description()}
		})
		loves := utils.Map2(result.GetLoves(), func(l entities.Love) models.Love {
			return models.Love{
				ID:        l.ID,
				PostID:    l.PostID,
				CommentID: l.CommentID,
				AuthorID:  l.AuthorID,
				CreatedAt: l.CreatedAt,
			}
		})

		return &models.Post{
			ID:        result.GetPostID(),
			AuthorID:  result.GetAuthorID(),
			Content:   result.GetPostTextContent().Content(),
			Images:    images,
			Videos:    videos,
			Loves:     loves,
			CreatedAt: result.GetCreatedAt(),
			UpdatedAt: result.GetUpdatedAt(),
		}, nil
	}
}

func makeCreateCommentEndpoint(cs services.CommentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CreateCommentReq)
		result, err := cs.CreateComment(ctx, req)
		if err != nil {
			return nil, err
		}

		loves := utils.Map2(result.GetLoves(), func(l entities.Love) models.Love {
			return models.Love{
				ID:        l.ID,
				PostID:    l.PostID,
				CommentID: l.CommentID,
				AuthorID:  l.AuthorID,
				CreatedAt: l.CreatedAt,
			}
		})

		return &models.Comment{
			ID:           result.GetID(),
			PostID:       result.GetPostID(),
			AuthorID:     result.GetAuthorID(),
			ParentID:     result.GetParentID(),
			IsPostParent: result.GetIsPostParent(),
			Image: models.Image{
				URL:         result.GetImageContent().URL(),
				Description: result.GetImageContent().Description(),
			},
			Video: models.Video{
				URL:         result.GetVideoContent().URL(),
				Description: result.GetVideoContent().Description(),
			},
			Loves:           loves,
			SubcommentCount: result.CountSubcomments(),
			CreatedAt:       result.GetCreatedAt(),
			UpdatedAt:       result.GetUpdatedAt(),
		}, nil
	}
}
