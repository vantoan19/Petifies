package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

type PetifiesProposalAcceptedEvent models.PetifiesProposalEvent

func NewPetifiesProposalAcceptedEvent(petifiesProposal *petifiesproposalaggre.PetifiesProposalAggre) PetifiesProposalAcceptedEvent {
	return PetifiesProposalAcceptedEvent{
		ID:                petifiesProposal.GetID(),
		UserID:            petifiesProposal.GetUserID(),
		PetifiesSessionID: petifiesProposal.GetPetifiesSessionID(),
		Status:            string(valueobjects.PetifiesProposalStatusAccepted),
		EventType:         models.PETIFIES_PROPOSAL_STATUS_CHANGED,
		CreatedAt:         time.Now(),
	}
}
