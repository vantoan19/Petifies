package listeners

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	postfeedaggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/post-feed"
)

type PostEventListener interface {
	PostCreated(ctx context.Context, event models.PostEvent) (*postfeedaggre.PostFeedAggre, error)
}
