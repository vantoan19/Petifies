package postaggre

import (
	"context"

	"github.com/google/uuid"
)

type PostRepository interface {
	GetByUUID(ctx context.Context, id uuid.UUID) (Post, error)
	SavePost(ctx context.Context, user Post) (Post, error)
	UpdatePost(ctx context.Context, user Post) (Post, error)
	DeleteByUUID(ctx context.Context, id uuid.UUID) (Post, error)
}
