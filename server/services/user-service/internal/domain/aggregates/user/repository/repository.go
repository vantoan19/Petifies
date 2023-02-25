package user

import (
	"context"

	"github.com/google/uuid"

	userAggre "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user"
)

type UserRepository interface {
	GetByUUID(ctx context.Context, id uuid.UUID) (*userAggre.User, error)
	GetByEmail(ctx context.Context, email string) (*userAggre.User, error)
	SaveUser(ctx context.Context, user userAggre.User) (*userAggre.User, error)
	UpdateUser(ctx context.Context, user *userAggre.User) (*userAggre.User, error)
	DeleteByUUID(ctx context.Context, id uuid.UUID) (*userAggre.User, error)
	DeleteByEmail(ctx context.Context, email string) (*userAggre.User, error)
}
