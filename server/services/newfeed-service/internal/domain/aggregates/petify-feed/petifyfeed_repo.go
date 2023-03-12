package petifyfeedaggre

import (
	"context"

	"github.com/google/uuid"
	paginateutils "github.com/vantoan19/Petifies/server/libs/paginate-utils"
)

type PetifyFeedRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID, pageToken paginateutils.PageToken) ([]*PetifyFeedAggre, error)
	GetByID(ctx context.Context, id uuid.UUID) (*PetifyFeedAggre, error)
	ExistsPetifyFeed(ctx context.Context, userID uuid.UUID, petifyID uuid.UUID) (bool, error)
	Save(ctx context.Context, petify PetifyFeedAggre) (*PetifyFeedAggre, error)
	Update(ctx context.Context, petify PetifyFeedAggre) (*PetifyFeedAggre, error)
	DeleteByID(ctx context.Context, id uuid.UUID) (*PetifyFeedAggre, error)
}
