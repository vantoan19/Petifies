package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

type PetifiesSessionProposalAcceptedEvent models.PetifiesSessionEvent

func NewPetifiesSessionProposalAcceptedEvent(petifiesSession *petifiessessionaggre.PetifiesSessionAggre) PetifiesSessionProposalAcceptedEvent {
	return PetifiesSessionProposalAcceptedEvent{
		ID:         petifiesSession.GetID(),
		PetifiesID: petifiesSession.GetPetifiesID(),
		Status:     string(valueobjects.PetifiesSessionStatusProposalAccepted),
		EventType:  models.PETIFIES_SESSION_STATUS_CHANGED,
		CreatedAt:  time.Now(),
	}
}
