package reviewaggre

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

type ReviewAggre struct {
	review *entities.Review
}

func NewReviewAggregate(
	id,
	petifiesID,
	authorID uuid.UUID,
	image valueobjects.Image,
	review string,
	createdAt,
	updatedAt time.Time,
) (*ReviewAggre, error) {
	reviewEntity := &entities.Review{
		ID:         id,
		PetifiesID: petifiesID,
		AuthorID:   authorID,
		Image:      image,
		Review:     review,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	if errs := reviewEntity.Validate(); errs.Exist() {
		return nil, errs[0]
	}

	return &ReviewAggre{
		review: reviewEntity,
	}, nil
}

// ======== Aggregate Root Getters ========

func (r *ReviewAggre) GetID() uuid.UUID {
	return r.review.ID
}

func (r *ReviewAggre) GetPetifiesID() uuid.UUID {
	return r.review.PetifiesID
}

func (r *ReviewAggre) GetAuthorID() uuid.UUID {
	return r.review.AuthorID
}

func (r *ReviewAggre) GetImage() valueobjects.Image {
	return r.review.Image
}

func (r *ReviewAggre) SetImage(image valueobjects.Image) error {
	if errs := image.Validate(); errs != nil {
		return errs[0]
	}

	r.review.Image = image
	return nil
}

func (r *ReviewAggre) GetReview() string {
	return r.review.Review
}

func (r *ReviewAggre) SetReview(review string) error {
	if len(strings.TrimSpace(review)) == 0 {
		return entities.ErrReviewRequired
	}
	if len(review) > 5000 {
		return entities.ErrReviewExceedsLimit
	}

	r.review.Review = review
	return nil
}

func (r *ReviewAggre) GetCreatedAt() time.Time {
	return r.review.CreatedAt
}

func (r *ReviewAggre) GetUpdatedAt() time.Time {
	return r.review.UpdatedAt
}
