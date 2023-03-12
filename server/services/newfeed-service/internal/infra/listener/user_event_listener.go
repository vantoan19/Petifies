package listener

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	useraggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/user/cassandra"
)

type UserEventListener struct {
	userRepo *cassandra.UserRepository
}

func NewUserEventListener(userRepo *cassandra.UserRepository) *UserEventListener {
	return &UserEventListener{
		userRepo: userRepo,
	}
}

func (ul *UserEventListener) UserCreated(ctx context.Context, event models.UserEvent) (*useraggre.UserAggre, error) {
	logger.Info("Start UserCreated")

	user, err := useraggre.NewUserAggregate(event.ID, event.Email)
	if err != nil {
		logger.ErrorData("Finish UserCreated: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish UserCreated: Success")
	return ul.userRepo.Save(ctx, *user)
}
