package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

type UserCacheRepository interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error)
	SetUser(ctx context.Context, userID uuid.UUID, user models.User) error
	ExistsUser(ctx context.Context, userID uuid.UUID) (bool, error)
}
