package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	reviewservice "github.com/vantoan19/Petifies/server/services/petifies-service/internal/application/services/petifies-review-service"
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

type ReviewEndpoints struct {
	CreateReview            endpoint.Endpoint
	EditReview              endpoint.Endpoint
	GetReviewById           endpoint.Endpoint
	ListReviewsByIds        endpoint.Endpoint
	ListReviewsByPetifiesId endpoint.Endpoint
}

func NewReviewEndpoints(rs reviewservice.ReviewService) ReviewEndpoints {
	return ReviewEndpoints{
		CreateReview:            makeCreateReviewEndpoint(rs),
		EditReview:              makeEditReviewEndpoint(rs),
		GetReviewById:           makeGetReviewByIdEndpoint(rs),
		ListReviewsByIds:        makeListReviewsByIdsEndpoint(rs),
		ListReviewsByPetifiesId: makeListReviewsByPetifiesIdEndpoint(rs),
	}
}

// ================ Endpoint Makers ==================

func makeCreateReviewEndpoint(rs reviewservice.ReviewService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CreateReviewReq)
		result, err := rs.CreateReview(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapReviewAggregateToReviewModel(result), nil
	}
}

func makeEditReviewEndpoint(rs reviewservice.ReviewService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.EditReviewReq)
		result, err := rs.EditReview(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapReviewAggregateToReviewModel(result), nil
	}
}

func makeGetReviewByIdEndpoint(rs reviewservice.ReviewService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.GetReviewByIdReq)
		result, err := rs.GetById(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return mapReviewAggregateToReviewModel(result), nil
	}
}

func makeListReviewsByIdsEndpoint(rs reviewservice.ReviewService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListReviewsByIdsReq)
		results, err := rs.ListByIds(ctx, req.ReviewIDs)
		if err != nil {
			return nil, err
		}

		return &models.ManyReviews{
			Reviews: commonutils.Map2(results, func(r *reviewaggre.ReviewAggre) *models.Review {
				return mapReviewAggregateToReviewModel(r)
			}),
		}, nil
	}
}

func makeListReviewsByPetifiesIdEndpoint(rs reviewservice.ReviewService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListReviewsByPetifiesIdReq)
		results, err := rs.ListByPetifiesId(ctx, req.PetifiesID, req.PageSize, req.AfterID)
		if err != nil {
			return nil, err
		}

		return &models.ManyReviews{
			Reviews: commonutils.Map2(results, func(r *reviewaggre.ReviewAggre) *models.Review {
				return mapReviewAggregateToReviewModel(r)
			}),
		}, nil
	}
}

// ================ Mappers ======================

func mapReviewAggregateToReviewModel(review *reviewaggre.ReviewAggre) *models.Review {
	return &models.Review{
		ID:         review.GetID(),
		PetifiesID: review.GetPetifiesID(),
		AuthorID:   review.GetAuthorID(),
		Review:     review.GetReview(),
		Image:      models.Image{URI: review.GetImage().GetURI(), Description: review.GetImage().GetDescription()},
		CreatedAt:  review.GetCreatedAt(),
		UpdatedAt:  review.GetUpdatedAt(),
	}
}
