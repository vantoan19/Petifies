package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PetifiesEventType string

const (
	PETIFIES_CREATED        PetifiesEventType = "PETIFIES_CREATED"
	PETIFIES_UPDATED        PetifiesEventType = "PETIFIES_UPDATED"
	PETIFIES_STATUS_CHANGED PetifiesEventType = "PETIFIES_STATUS_CHANGED"
	PETIFIES_DELETED        PetifiesEventType = "PETIFIES_DELETED"
)

type PetifiesEvent struct {
	ID          uuid.UUID         `json:"id"`
	OwnerID     uuid.UUID         `json:"owner_id"`
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Longitude   float64           `json:"longitude"`
	Latitude    float64           `json:"latitude"`
	EventType   PetifiesEventType `json:"event_type"`
	CreatedAt   time.Time         `json:"created_at"`
}

func (u *PetifiesEvent) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *PetifiesEvent) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}
