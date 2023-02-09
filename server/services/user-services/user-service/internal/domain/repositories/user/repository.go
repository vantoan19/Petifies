package user

import (
	"github.com/google/uuid"

	userAggre "github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/domain/valueobjects"
)

type UserRepository interface {
	GetByUUID(uuid.UUID) (userAggre.User, error)
	GetByEmail(string) (userAggre.User, error)
	Add(userAggre.User) (userAggre.User, error)
	UpdateName(uuid.UUID, valueobjects.Name) (userAggre.User, error)
	DeleteByUUID(uuid.UUID) error
	DeleteByEmail(string) error
}
