package outbox_repo

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
)

type State string

const (
	StartedState   = "STARTED"
	CompletedState = "COMPLETED"
)

type Event struct {
	ID          uuid.UUID
	Payload     models.KafkaMessage
	OutboxState State
	LockedBy    *uuid.UUID
	LockedAt    *time.Time
	Error       *string
	CompletedAt *time.Time
	CreatedAt   time.Time
}

type EventRepository interface {
	AddEvent(event Event) (*Event, error)
	GetEventsByLockerID(lockerID uuid.UUID) ([]*Event, error)
	LockStartedEvents(lockerID uuid.UUID) error
	UpdateEvent(event Event) error
	UnlockEventsByLockerID(lockerID uuid.UUID) error
	UnlockEventsBeforeDatetime(t time.Time) error
	DeleteEventsBeforeDatetime(t time.Time) error
}
