package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       string
	FirstName   string
	LastName    string
	IsActivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
