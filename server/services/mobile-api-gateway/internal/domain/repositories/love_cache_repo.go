package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type LoveCacheRepository interface {
	GetLove(ctx context.Context, authorID, targetID uuid.UUID) (*models.Love, error)
	SetLove(ctx context.Context, authorID, targetID uuid.UUID, love models.Love) error
	RemoveLove(ctx context.Context, authorID, targetID uuid.UUID) error
	ExistsLove(ctx context.Context, authorID, targetID uuid.UUID) (bool, error)
}
