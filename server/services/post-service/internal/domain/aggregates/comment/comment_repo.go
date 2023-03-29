package commentaggre

import (
	"context"

	"github.com/google/uuid"
)

type CommentRepository interface {
	GetByParentID(ctx context.Context, parentID uuid.UUID, pageSize int, afterCommentID uuid.UUID) ([]*Comment, error)
	GetByUUID(ctx context.Context, id uuid.UUID) (*Comment, error)
	GetCommentAncestors(ctx context.Context, id uuid.UUID) ([]*Comment, error)
	CountCommentByParentID(ctx context.Context, parentID uuid.UUID) (int, error)
	SaveComment(ctx context.Context, comment Comment) (*Comment, error)
	UpdateComment(ctx context.Context, comment Comment) (*Comment, error)
	DeleteByUUID(ctx context.Context, id uuid.UUID) (*Comment, error)
}
