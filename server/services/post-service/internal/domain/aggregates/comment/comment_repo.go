package commentaggre

import (
	"context"

	"github.com/google/uuid"
)

type CommentRepository interface {
	GetByUUID(ctx context.Context, id uuid.UUID) (*Comment, error)
	SaveComment(ctx context.Context, comment Comment) (*Comment, error)
	UpdateComment(ctx context.Context, comment Comment) (*Comment, error)
	DeleteByUUID(ctx context.Context, id uuid.UUID) (*Comment, error)
}
