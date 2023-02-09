package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/domain/valueobjects"
)

type User struct {
	// Identifier of the user
	ID uuid.UUID `validate:"required"`
	// Email of the user
	Email string `validate:"required,email,max=300"`
	// Username
	Password string `validate:"required,max=50"`
	// Name
	Name valueobjects.Name
	// Timestamp when the user is created
	CreatedAt time.Time `validate:"required"`
	// Timestamp when the user is updated
	UpdatedAt time.Time `validate:"required"`
}
