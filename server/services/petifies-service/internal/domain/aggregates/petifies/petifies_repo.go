package petifiesaggre

import (
	"context"

	"github.com/google/uuid"
)

type PetifiesRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*PetifiesAggre, error)
	GetByID(ctx context.Context, id uuid.UUID) (*PetifiesAggre, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*PetifiesAggre, error)
	Save(ctx context.Context, petifies PetifiesAggre) (*PetifiesAggre, error)
	Update(ctx context.Context, petifies PetifiesAggre) (*PetifiesAggre, error)

	GetByUserIDWithSession(ctx context.Context, userID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*PetifiesAggre, error)
	GetByIDWithSession(ctx context.Context, id uuid.UUID) (*PetifiesAggre, error)
	ListByIDsWithSession(ctx context.Context, ids []uuid.UUID) ([]*PetifiesAggre, error)
	SaveWithSession(ctx context.Context, petifies PetifiesAggre) (*PetifiesAggre, error)
	UpdateWithSession(ctx context.Context, petifies PetifiesAggre) (*PetifiesAggre, error)
}
