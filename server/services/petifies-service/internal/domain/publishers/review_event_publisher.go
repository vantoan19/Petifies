package publishers

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type ReviewEventMessagePublisher interface {
	Publish(ctx context.Context, event models.ReviewEvent) error
}
