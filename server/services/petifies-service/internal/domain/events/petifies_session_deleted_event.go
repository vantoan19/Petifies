package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
)

type PetifiesSessionDeletedEvent models.PetifiesSessionEvent

func NewPetifiesSessionDeletedEvent(petifiesSession *petifiessessionaggre.PetifiesSessionAggre) PetifiesSessionDeletedEvent {
	return PetifiesSessionDeletedEvent{
		ID:         petifiesSession.GetID(),
		PetifiesID: petifiesSession.GetPetifiesID(),
		Status:     string(petifiesSession.GetStatus()),
		EventType:  models.PETIFIES_SESSION_DELETED,
		CreatedAt:  time.Now(),
	}
}
