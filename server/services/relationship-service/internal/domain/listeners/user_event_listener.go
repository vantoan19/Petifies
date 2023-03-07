package listeners

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	useraggre "github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user"
)

type UserEventListener interface {
	UserCreated(ctx context.Context, event models.UserEvent) (*useraggre.UserAggregate, error)
}
