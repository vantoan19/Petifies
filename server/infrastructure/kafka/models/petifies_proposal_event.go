package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PetifiesProposalEventType string

const (
	PETIFIES_PROPOSAL_CREATED        PetifiesProposalEventType = "PETIFIES_PROPOSAL_CREATED"
	PETIFIES_PROPOSAL_STATUS_CHANGED PetifiesProposalEventType = "PETIFIES_PROPOSAL_STATUS_CHANGED"
	PETIFIES_PROPOSAL_UPDATED        PetifiesProposalEventType = "PETIFIES_PROPOSAL_UPDATED"
	PETIFIES_PROPOSAL_DELETED        PetifiesProposalEventType = "PETIFIES_PROPOSAL_DELETED"
)

type PetifiesProposalEvent struct {
	ID                uuid.UUID                 `json:"id"`
	UserID            uuid.UUID                 `json:"user_id"`
	PetifiesSessionID uuid.UUID                 `json:"petifies_session_id"`
	Status            string                    `json:"status"`
	EventType         PetifiesProposalEventType `json:"event_type"`
	CreatedAt         time.Time                 `json:"created_at"`
}

func (u *PetifiesProposalEvent) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *PetifiesProposalEvent) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}
