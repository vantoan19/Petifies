package loveaggre

import (
	"context"

	"github.com/google/uuid"
	paginateutils "github.com/vantoan19/Petifies/server/libs/paginate-utils"
)

type LoveRepository interface {
	GetByTargetID(ctx context.Context, targetID uuid.UUID, pageToken paginateutils.PageToken) ([]*Love, error)
	GetByTargetIDAndAuthorID(ctx context.Context, authorID uuid.UUID, targetID uuid.UUID) (*Love, error)
	CountLoveByTargetID(ctx context.Context, targetID uuid.UUID) (int, error)
	SaveLove(ctx context.Context, love Love) (*Love, error)
	DeleteByUUID(ctx context.Context, id uuid.UUID) error
	DeleteByTargetID(ctx context.Context, targetID uuid.UUID) error
}
