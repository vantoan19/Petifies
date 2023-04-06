package models

import (
	"time"

	"github.com/google/uuid"
)

type OutboxState string

const (
	OutboxStateSTARTED   OutboxState = "STARTED"
	OutboxStateCOMPLETED OutboxState = "COMPLETED"
)

type PetifiesEvent struct {
	ID          uuid.UUID   `bson:"id"`
	Payload     string      `bson:"payload"`
	OutboxState OutboxState `bson:"outbox_state"`
	LockedBy    *uuid.UUID  `bson:"locked_by"`
	LockedAt    *time.Time  `bson:"locked_at"`
	Error       *string     `bson:"error"`
	CompletedAt *time.Time  `bson:"completed_at"`
	CreatedAt   time.Time   `bson:"created_at"`
}
