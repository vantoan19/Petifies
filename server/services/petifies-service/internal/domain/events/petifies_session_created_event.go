package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
)

type PetifiesSessionCreatedEvent models.PetifiesSessionEvent

func NewPetifiesSessionCreatedEvent(petifiesSession *petifiessessionaggre.PetifiesSessionAggre) PetifiesSessionCreatedEvent {
	return PetifiesSessionCreatedEvent{
		ID:         petifiesSession.GetID(),
		PetifiesID: petifiesSession.GetPetifiesID(),
		Status:     string(petifiesSession.GetStatus()),
		EventType:  models.PETIFIES_SESSION_CREATED,
		CreatedAt:  time.Now(),
	}
}
