package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*UserAggre, error)
	Save(ctx context.Context, user UserAggre) (*UserAggre, error)
	Update(ctx context.Context, user UserAggre) (*UserAggre, error)
	DeleteByID(ctx context.Context, id uuid.UUID) (*UserAggre, error)
}
