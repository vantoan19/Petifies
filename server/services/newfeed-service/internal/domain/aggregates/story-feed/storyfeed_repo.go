package storyfeedaggre

import (
	"context"

	"github.com/google/uuid"
	paginateutils "github.com/vantoan19/Petifies/server/libs/paginate-utils"
)

type StoryFeedRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID, pageToken paginateutils.PageToken) ([]*StoryFeedAggre, error)
	GetByID(ctx context.Context, id uuid.UUID) (*StoryFeedAggre, error)
	ExistsStoryFeed(ctx context.Context, userID uuid.UUID, storyID uuid.UUID) (bool, error)
	Save(ctx context.Context, story StoryFeedAggre) (*StoryFeedAggre, error)
	Update(ctx context.Context, story StoryFeedAggre) (*StoryFeedAggre, error)
	DeleteByID(ctx context.Context, id uuid.UUID) (*StoryFeedAggre, error)
}
