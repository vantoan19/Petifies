package publishers

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type PetifiesEventMessagePublisher interface {
	Publish(ctx context.Context, event models.PetifiesEvent) error
}
