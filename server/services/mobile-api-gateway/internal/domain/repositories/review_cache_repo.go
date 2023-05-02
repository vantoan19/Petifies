package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

type ReviewCacheRepository interface {
	GetReview(ctx context.Context, reviewId uuid.UUID) (*models.Review, error)
	SetReview(ctx context.Context, reviewId uuid.UUID, review *models.Review) error
}
