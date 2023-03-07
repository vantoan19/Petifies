package publisher

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type UserRequestMessagePublisher interface {
	Publish(ctx context.Context, event models.UserEvent) error
}
