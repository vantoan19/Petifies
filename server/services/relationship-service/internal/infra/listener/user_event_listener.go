package listener

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	useraggre "github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/entities"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/infra/repositories/user/neo4j"
)

type UserEventListener struct {
	userRepo *neo4j.UserRepository
}

func NewUserEventListener(userRepo *neo4j.UserRepository) *UserEventListener {
	return &UserEventListener{
		userRepo: userRepo,
	}
}

func (ul *UserEventListener) UserCreated(ctx context.Context, event models.UserEvent) (*useraggre.UserAggregate, error) {
	user, err := useraggre.NewUserAggregate(entities.User{
		ID:    event.ID,
		Email: event.Email,
	})
	if err != nil {
		return nil, err
	}

	err = ul.userRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser, err := ul.userRepo.GetByUUID(ctx, user.GetID())
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
