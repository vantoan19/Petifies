package listener

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	useraggre "github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/entities"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/infra/repositories/user/neo4j"
)

var logger = logging.New("RelationshipService.UserEventListener")

type UserEventListener struct {
	userRepo *neo4j.UserRepository
}

func NewUserEventListener(userRepo *neo4j.UserRepository) *UserEventListener {
	return &UserEventListener{
		userRepo: userRepo,
	}
}

func (ul *UserEventListener) UserCreated(ctx context.Context, event models.UserEvent) (*useraggre.UserAggregate, error) {
	logger.Info("Start UserCreated")

	user, err := useraggre.NewUserAggregate(entities.User{
		ID:    event.ID,
		Email: event.Email,
	})
	if err != nil {
		logger.ErrorData("Finished UserCreated: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = ul.userRepo.SaveUser(ctx, user)
	if err != nil {
		logger.ErrorData("Finished UserCreated: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	createdUser, err := ul.userRepo.GetByUUID(ctx, user.GetID())
	if err != nil {
		logger.ErrorData("Finished UserCreated: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish UserCreated: SUCCESS")
	return createdUser, nil
}
