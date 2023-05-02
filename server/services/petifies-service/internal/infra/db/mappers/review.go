package mappers

import (
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/models"
)

func DbModelToReviewAggregate(r *models.Review) (*reviewaggre.ReviewAggre, error) {
	review, err := reviewaggre.NewReviewAggregate(
		r.ID,
		r.PetifiesID,
		r.AuthorID,
		valueobjects.NewImage(r.Image.URI, r.Image.Description),
		r.Review,
		r.CreatedAt,
		r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func AggregateReviewToDbReview(r *reviewaggre.ReviewAggre) *models.Review {
	return &models.Review{
		ID:         r.GetID(),
		PetifiesID: r.GetPetifiesID(),
		AuthorID:   r.GetAuthorID(),
		Review:     r.GetReview(),
		Image:      models.Image{URI: r.GetImage().GetURI(), Description: r.GetImage().GetDescription()},
		CreatedAt:  r.GetCreatedAt(),
		UpdatedAt:  r.GetUpdatedAt(),
	}
}
