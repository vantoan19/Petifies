package repository

import (
	"context"

	aggregates "github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media"
)

type MediaRepository interface {
	Save(ctx context.Context, media *aggregates.Media) (string, error)
	RemoveByUri(ctx context.Context, uri string) error
}
