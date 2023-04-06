package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ReviewEventType string

const (
	REVIEW_CREATED ReviewEventType = "REVIEW_CREATED"
	REVIEW_UPDATED ReviewEventType = "REVIEW_UPDATED"
	REVIEW_DELETED ReviewEventType = "REVIEW_DELETED"
)

type ReviewEvent struct {
	ID         uuid.UUID       `json:"id"`
	PetifiesID uuid.UUID       `json:"petifies_id"`
	AuthorID   uuid.UUID       `json:"author_id"`
	EventType  ReviewEventType `json:"event_type"`
	CreatedAt  time.Time       `json:"created_at"`
}

func (u *ReviewEvent) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *ReviewEvent) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}
