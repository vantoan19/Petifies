package useraggre

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByUUID(ctx context.Context, id uuid.UUID) (*UserAggregate, error)
	SaveUser(ctx context.Context, user UserAggregate) error
	UpdateUser(ctx context.Context, user UserAggregate) error
	DeleteByUUID(ctx context.Context, id uuid.UUID) error
}
