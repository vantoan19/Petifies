package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type PostCacheRepository interface {
	GetPostContent(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	ListPostContents(ctx context.Context, postIds []uuid.UUID) ([]*models.Post, error)
	SetPostContent(ctx context.Context, postID uuid.UUID, post models.Post) error
	ExistsPostContent(ctx context.Context, postID uuid.UUID) (bool, error)
}
