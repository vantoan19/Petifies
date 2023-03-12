package publisher

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type PostEventMessagePublisher interface {
	Publish(ctx context.Context, event models.PostEvent) error
}
