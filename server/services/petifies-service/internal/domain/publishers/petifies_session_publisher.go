package publishers

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type PetifiesSessionEventMessagePublisher interface {
	Publish(ctx context.Context, event models.PetifiesSessionEvent) error
}
