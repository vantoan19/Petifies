package mappers

import (
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/models"
)

func DbModelToPetifiesSessionAggregate(p *models.PetifiesSession) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	petifies, err := petifiessessionaggre.NewPetitifesSessionAggre(
		p.ID,
		p.PetifiesID,
		p.Time.FromTime,
		p.Time.ToTime,
		valueobjects.PetifiesSessionStatus(p.Status),
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return petifies, nil
}

func AggregatePetifiesSessionToDbPetifiesSession(p *petifiessessionaggre.PetifiesSessionAggre) *models.PetifiesSession {
	return &models.PetifiesSession{
		ID:         p.GetID(),
		PetifiesID: p.GetPetifiesID(),
		Time: models.TimeRange{
			FromTime: p.GetTime().GetFromTime(),
			ToTime:   p.GetTime().GetToTime(),
		},
		Status:    string(p.GetStatus()),
		CreatedAt: p.GetCreatedAt(),
		UpdatedAt: p.GetUpdatedAt(),
	}
}
