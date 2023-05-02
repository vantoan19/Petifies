package reviewaggre

import (
	"context"

	"github.com/google/uuid"
)

type ReviewRepository interface {
	GetByPetifiesID(ctx context.Context, petifiesID uuid.UUID, pageSize int, afterId uuid.UUID) ([]*ReviewAggre, error)
	GetByUserId(ctx context.Context, userId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*ReviewAggre, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ReviewAggre, error)
	ListByIds(ctx context.Context, ids []uuid.UUID) ([]*ReviewAggre, error)
	Save(ctx context.Context, review ReviewAggre) (*ReviewAggre, error)
	Update(ctx context.Context, review ReviewAggre) (*ReviewAggre, error)

	GetByPetifiesIDWithSession(ctx context.Context, petifiesID uuid.UUID, pageSize int, afterId uuid.UUID) ([]*ReviewAggre, error)
	GetByUserIdWithSession(ctx context.Context, userId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*ReviewAggre, error)
	GetByIDWithSession(ctx context.Context, id uuid.UUID) (*ReviewAggre, error)
	ListByIdsWithSession(ctx context.Context, ids []uuid.UUID) ([]*ReviewAggre, error)
	SaveWithSession(ctx context.Context, review ReviewAggre) (*ReviewAggre, error)
	UpdateWithSession(ctx context.Context, review ReviewAggre) (*ReviewAggre, error)
}
