package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

type PetifiesProposalCacheRepository interface {
	GetPetifiesProposal(ctx context.Context, proposalId uuid.UUID) (*models.PetifiesProposal, error)
	SetPetifiesProposal(ctx context.Context, proposalId uuid.UUID, proposal *models.PetifiesProposal) error
}
