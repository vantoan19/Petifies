package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
)

type PetifiesProposalUpdatedEvent models.PetifiesProposalEvent

func NewPetifiesProposalUpdatedEvent(petifiesProposal *petifiesproposalaggre.PetifiesProposalAggre) PetifiesProposalUpdatedEvent {
	return PetifiesProposalUpdatedEvent{
		ID:                petifiesProposal.GetID(),
		UserID:            petifiesProposal.GetUserID(),
		PetifiesSessionID: petifiesProposal.GetPetifiesSessionID(),
		Status:            string(petifiesProposal.GetStatus()),
		EventType:         models.PETIFIES_PROPOSAL_UPDATED,
		CreatedAt:         time.Now(),
	}
}
