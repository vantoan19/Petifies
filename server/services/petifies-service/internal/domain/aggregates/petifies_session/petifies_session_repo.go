package petifiessessionaggre

import (
	"context"

	"github.com/google/uuid"
)

type PetifiesSessionRepository interface {
	GetByPetifiesID(ctx context.Context, petifiesID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*PetifiesSessionAggre, error)
	GetByID(ctx context.Context, id uuid.UUID) (*PetifiesSessionAggre, error)
	ListByIds(ctx context.Context, ids []uuid.UUID) ([]*PetifiesSessionAggre, error)
	Save(ctx context.Context, session PetifiesSessionAggre) (*PetifiesSessionAggre, error)
	Update(ctx context.Context, session PetifiesSessionAggre) (*PetifiesSessionAggre, error)

	GetByPetifiesIDWithSession(ctx context.Context, petifiesID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*PetifiesSessionAggre, error)
	GetByIDWithSession(ctx context.Context, id uuid.UUID) (*PetifiesSessionAggre, error)
	ListByIdsWithSession(ctx context.Context, ids []uuid.UUID) ([]*PetifiesSessionAggre, error)
	SaveWithSession(ctx context.Context, session PetifiesSessionAggre) (*PetifiesSessionAggre, error)
	UpdateWithSession(ctx context.Context, session PetifiesSessionAggre) (*PetifiesSessionAggre, error)
}
