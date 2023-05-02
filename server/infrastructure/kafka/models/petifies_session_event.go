package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PetifiesSessionEventType string

const (
	PETIFIES_SESSION_CREATED        PetifiesSessionEventType = "PETIFIES_SESSION_CREATED"
	PETIFIES_SESSION_STATUS_CHANGED PetifiesSessionEventType = "PETIFIES_SESSION_STATUS_CHANGED"
	PETIFIES_SESSION_UPDATED        PetifiesSessionEventType = "PETIFIES_SESSION_UPDATED"
	PETIFIES_SESSION_DELETED        PetifiesSessionEventType = "PETIFIES_SESSION_DELETED"
)

type PetifiesSessionEvent struct {
	ID         uuid.UUID                `json:"id"`
	PetifiesID uuid.UUID                `json:"petifies_id"`
	Status     string                   `json:"status"`
	EventType  PetifiesSessionEventType `json:"event_type"`
	CreatedAt  time.Time                `json:"created_at"`
}

func (u *PetifiesSessionEvent) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *PetifiesSessionEvent) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}
