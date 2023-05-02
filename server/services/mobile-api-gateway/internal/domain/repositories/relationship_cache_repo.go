package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/models"
)

type RelationshipCacheRepository interface {
	GetFollowingsInfo(ctx context.Context, userID uuid.UUID) (*models.ListFollowingsResp, error)
	SetFollowingsInfo(ctx context.Context, userID uuid.UUID, followings *models.ListFollowingsResp) error
	ExistsFollowingsInfo(ctx context.Context, userID uuid.UUID) (bool, error)

	GetFollowersInfo(ctx context.Context, userID uuid.UUID) (*models.ListFollowersResp, error)
	SetFollowersInfo(ctx context.Context, userID uuid.UUID, followers *models.ListFollowersResp) error
	ExistsFollowersInfo(ctx context.Context, userID uuid.UUID) (bool, error)
}
