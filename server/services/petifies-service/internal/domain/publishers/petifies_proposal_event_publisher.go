package publishers

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type PetifiesProposalEventMessagePublisher interface {
	Publish(ctx context.Context, event models.PetifiesProposalEvent) error
}
