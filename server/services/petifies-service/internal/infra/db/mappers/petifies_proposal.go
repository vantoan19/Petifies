package mappers

import (
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/models"
)

func DbModelToPetifiesProposalAggregate(p *models.PetifiesProposal) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	petifies, err := petifiesproposalaggre.NewPetifiesProposalAggregate(
		p.ID,
		p.UserID,
		p.PetifiesSessionID,
		p.Proposal,
		valueobjects.PetifiesProposalStatus(p.Status),
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return petifies, nil
}

func AggregatePetifiesProposalToDbPetifiesProposal(p *petifiesproposalaggre.PetifiesProposalAggre) *models.PetifiesProposal {
	return &models.PetifiesProposal{
		ID:                p.GetID(),
		UserID:            p.GetUserID(),
		PetifiesSessionID: p.GetPetifiesSessionID(),
		Proposal:          p.GetProposal(),
		Status:            string(p.GetStatus()),
		CreatedAt:         p.GetCreatedAt(),
		UpdatedAt:         p.GetUpdatedAt(),
	}
}
