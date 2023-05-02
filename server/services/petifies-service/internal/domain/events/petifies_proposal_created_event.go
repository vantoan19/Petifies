package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
)

type PetifiesProposalCreatedEvent models.PetifiesProposalEvent

func NewPetifiesProposalCreatedEvent(petifiesProposal *petifiesproposalaggre.PetifiesProposalAggre) PetifiesProposalCreatedEvent {
	return PetifiesProposalCreatedEvent{
		ID:                petifiesProposal.GetID(),
		UserID:            petifiesProposal.GetUserID(),
		PetifiesSessionID: petifiesProposal.GetPetifiesSessionID(),
		Status:            string(petifiesProposal.GetStatus()),
		EventType:         models.PETIFIES_PROPOSAL_CREATED,
		CreatedAt:         time.Now(),
	}
}
