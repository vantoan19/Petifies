package repositories

import (
	"context"

	"github.com/google/uuid"
)

type FeedCacheRepository interface {
	GetPostFeedIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	SetPostFeedIDs(ctx context.Context, userID uuid.UUID, feedIDs []uuid.UUID) error
	ExistsPostFeedIDs(ctx context.Context, userID uuid.UUID) (bool, error)
}
