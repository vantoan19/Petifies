package postfeedaggre

import (
	"context"

	"github.com/google/uuid"
)

type PostFeedRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID, pageSize int, afterPostID uuid.UUID) ([]*PostFeedAggre, error)
	ExistsPostFeed(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (bool, error)
	Save(ctx context.Context, post PostFeedAggre) (*PostFeedAggre, error)
	Update(ctx context.Context, post PostFeedAggre) (*PostFeedAggre, error)
	DeleteByUserAndPostID(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (*PostFeedAggre, error)
}
