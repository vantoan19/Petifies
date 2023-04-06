package listeners

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type PetifiesEventListener interface {
	Receive(ctx context.Context, event models.PetifiesEvent) error
}
