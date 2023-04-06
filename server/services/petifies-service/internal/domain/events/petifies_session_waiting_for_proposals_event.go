package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

type PetifiesSessionWaitingForProposalsEvent models.PetifiesSessionEvent

func NewPetifiesSessionWaitingForProposalsEvent(petifiesSession *petifiessessionaggre.PetifiesSessionAggre) PetifiesSessionWaitingForProposalsEvent {
	return PetifiesSessionWaitingForProposalsEvent{
		ID:         petifiesSession.GetID(),
		PetifiesID: petifiesSession.GetPetifiesID(),
		Status:     string(valueobjects.PetifiesSessionStatusWaitingForProposal),
		EventType:  models.PETIFIES_SESSION_STATUS_CHANGED,
		CreatedAt:  time.Now(),
	}
}
