package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
)

type PetifiesSessionUpdatedEvent models.PetifiesSessionEvent

func NewPetifiesSessionUpdatedEvent(petifiesSession *petifiessessionaggre.PetifiesSessionAggre) PetifiesSessionUpdatedEvent {
	return PetifiesSessionUpdatedEvent{
		ID:         petifiesSession.GetID(),
		PetifiesID: petifiesSession.GetPetifiesID(),
		Status:     string(petifiesSession.GetStatus()),
		EventType:  models.PETIFIES_SESSION_UPDATED,
		CreatedAt:  time.Now(),
	}
}
