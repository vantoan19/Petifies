package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	USER_CREATED UserStatus = "USER_CREATED"
	USER_UPDATED UserStatus = "USER_UPDATED"
	USER_DELETED UserStatus = "USER_DELETED"
)

type UserEvent struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	Status    UserStatus `json:"status"`
}

func (u *UserEvent) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserEvent) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}
