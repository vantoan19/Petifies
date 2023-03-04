package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type UserStatus string

const (
	USER_CREATED UserStatus = "USER_CREATED"
	USER_DELETED UserStatus = "USER_DELETED"
)

type UserRequest struct {
	ID     uuid.UUID  `json:"id"`
	Email  string     `json:"email"`
	Status UserStatus `json:"status"`
}

func (u *UserRequest) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserRequest) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}
