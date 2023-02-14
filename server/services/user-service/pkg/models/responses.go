package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserResp struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LoginResp struct {
	AccessToken string
}
