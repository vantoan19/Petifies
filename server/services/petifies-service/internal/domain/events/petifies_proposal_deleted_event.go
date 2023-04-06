package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
)

type PetifiesProposalDeletedEvent models.PetifiesProposalEvent

func NewPetifiesProposalDeletedEvent(petifiesProposal *petifiesproposalaggre.PetifiesProposalAggre) PetifiesProposalDeletedEvent {
	return PetifiesProposalDeletedEvent{
		ID:                petifiesProposal.GetID(),
		UserID:            petifiesProposal.GetUserID(),
		PetifiesSessionID: petifiesProposal.GetPetifiesSessionID(),
		Status:            string(petifiesProposal.GetStatus()),
		EventType:         models.PETIFIES_PROPOSAL_DELETED,
		CreatedAt:         time.Now(),
	}
}
