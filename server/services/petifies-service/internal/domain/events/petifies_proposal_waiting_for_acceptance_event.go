package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

type PetifiesProposalWaitingForAcceptanceEvent models.PetifiesProposalEvent

func NewPetifiesProposalWaitingForAcceptanceEvent(petifiesProposal *petifiesproposalaggre.PetifiesProposalAggre) PetifiesProposalWaitingForAcceptanceEvent {
	return PetifiesProposalWaitingForAcceptanceEvent{
		ID:                petifiesProposal.GetID(),
		UserID:            petifiesProposal.GetUserID(),
		PetifiesSessionID: petifiesProposal.GetPetifiesSessionID(),
		Status:            string(valueobjects.PetifiesProposalStatusWaitingForAcceptance),
		EventType:         models.PETIFIES_PROPOSAL_STATUS_CHANGED,
		CreatedAt:         time.Now(),
	}
}
