package mappers

import (
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/models"
)

func DbModelToPetifiesAggregate(p *models.Petifies) (*petifiesaggre.PetifiesAggre, error) {
	petifies, err := petifiesaggre.NewPetifiesAggregate(
		p.ID,
		p.OwnerID,
		valueobjects.PetifiesType(p.Type),
		p.Title,
		p.Description,
		p.PetName,
		commonutils.Map2(p.Images, func(i models.Image) valueobjects.Image { return valueobjects.NewImage(i.URI, i.Description) }),
		valueobjects.PetifiesStatus(p.Status),
		p.CreatedAt,
		p.UpdatedAt,
		entities.Address{
			AddressLineOne: p.Address.AddressLineOne,
			AddressLineTwo: p.Address.AddressLineTwo,
			Street:         p.Address.Street,
			District:       p.Address.District,
			City:           p.Address.City,
			Region:         p.Address.Region,
			PostalCode:     p.Address.PostalCode,
			Country:        p.Address.Country,
			Coordinates:    valueobjects.NewCoordinates(p.Address.Coordinates.Longitude, p.Address.Coordinates.Latitude),
		},
	)
	if err != nil {
		return nil, err
	}

	return petifies, nil
}

func AggregatePetifiesToDbPetifies(p *petifiesaggre.PetifiesAggre) *models.Petifies {
	address := p.GetAddress()

	return &models.Petifies{
		ID:          p.GetID(),
		OwnerID:     p.GetOwnerID(),
		Type:        string(p.GetType()),
		Title:       p.GetTitle(),
		Description: p.GetDescription(),
		Address: models.Address{
			AddressLineOne: address.AddressLineOne,
			AddressLineTwo: address.AddressLineTwo,
			Street:         address.Street,
			District:       address.District,
			City:           address.City,
			Region:         address.Region,
			PostalCode:     address.PostalCode,
			Country:        address.Country,
			Coordinates: models.Coordinates{
				Longitude: address.Coordinates.GetLongitude(),
				Latitude:  address.Coordinates.GetLatitude(),
			},
		},
		PetName: p.GetPetName(),
		Images: commonutils.Map2(p.GetImages(), func(i valueobjects.Image) models.Image {
			return models.Image{URI: i.GetURI(), Description: i.GetDescription()}
		}),
		Status:    string(p.GetStatus()),
		CreatedAt: p.GetCreatedAt(),
		UpdatedAt: p.GetUpdatedAt(),
	}
}
