package postfeedaggre

import (
	"context"

	"github.com/google/uuid"
	paginateutils "github.com/vantoan19/Petifies/server/libs/paginate-utils"
)

type PostFeedRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID, pageToken paginateutils.PageToken) ([]*PostFeedAggre, error)
	ExistsPostFeed(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (bool, error)
	Save(ctx context.Context, post PostFeedAggre) (*PostFeedAggre, error)
	Update(ctx context.Context, post PostFeedAggre) (*PostFeedAggre, error)
	DeleteByUserAndPostID(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (*PostFeedAggre, error)
}
