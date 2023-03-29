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
}
