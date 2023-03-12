package mapper

import (
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	useraggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/db/models"
)

func DbUserToUserAggregate(user *models.User) (*useraggre.UserAggre, error) {
	id := uuid.UUID(user.ID)

	return useraggre.NewUserAggregate(id, user.Email)
}

func UserAggregateToUserDb(user *useraggre.UserAggre) (*models.User, error) {
	return &models.User{
		ID:    gocql.UUID(user.GetID()),
		Email: user.GetEmail(),
	}, nil
}
