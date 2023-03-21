package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type PostCacheRepository interface {
	GetPostContent(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	SetPostContent(ctx context.Context, postID uuid.UUID, post models.Post) error
	ExistsPostContent(ctx context.Context, postID uuid.UUID) (bool, error)

	GetPostLoveCount(ctx context.Context, postID uuid.UUID) (int, error)
	SetPostLoveCount(ctx context.Context, postID uuid.UUID, loveCount int) error
	ExistsPostLoveCount(ctx context.Context, postID uuid.UUID) (bool, error)

	GetPostCommentCount(ctx context.Context, postID uuid.UUID) (int, error)
	SetPostCommentCount(ctx context.Context, postID uuid.UUID, commentCount int) error
	ExistsPostCommentCount(ctx context.Context, postID uuid.UUID) (bool, error)
}
