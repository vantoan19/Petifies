package models

import (
	"time"

	"github.com/google/uuid"
)

type PetifiesProposal struct {
	ID                uuid.UUID `bson:"id"`
	UserID            uuid.UUID `bson:"user_id"`
	PetifiesSessionID uuid.UUID `bson:"petifies_session_id"`
	Proposal          string    `bson:"proposal"`
	Status            string    `bson:"status"`
	CreatedAt         time.Time `bson:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at"`
}
