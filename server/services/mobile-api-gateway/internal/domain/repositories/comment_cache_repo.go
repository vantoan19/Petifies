package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type CommentCacheRepository interface {
	GetCommentContent(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)
	SetCommentContent(ctx context.Context, commentID uuid.UUID, comment models.Comment) error
	ExistsCommentContent(ctx context.Context, commentID uuid.UUID) (bool, error)

	GetCommentLoveCount(ctx context.Context, commentID uuid.UUID) (int, error)
	SetCommentLoveCount(ctx context.Context, commentID uuid.UUID, loveCount int) error
	ExistsCommentLoveCount(ctx context.Context, commentID uuid.UUID) (bool, error)

	GetCommentSubCommentCount(ctx context.Context, commentID uuid.UUID) (int, error)
	SetCommentSubCommentCount(ctx context.Context, commentID uuid.UUID, commentCount int) error
	ExistsCommentSubCommentCount(ctx context.Context, commentID uuid.UUID) (bool, error)
}
