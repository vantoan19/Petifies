package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	petifiesservice "github.com/vantoan19/Petifies/server/services/petifies-service/internal/application/services/petifies-service"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

type PetifiesEndpoints struct {
	CreatePetify          endpoint.Endpoint
	EditPetify            endpoint.Endpoint
	GetPetifyById         endpoint.Endpoint
	ListPetifiesByIds     endpoint.Endpoint
	ListPetifiesByOwnerId endpoint.Endpoint
}

func NewPetifiesEndpoints(ps petifiesservice.PetifesService) PetifiesEndpoints {
	return PetifiesEndpoints{
		CreatePetify:          makeCreatePetifyEndpoint(ps),
		EditPetify:            makeEditPetifyEndpoint(ps),
		GetPetifyById:         makeGetByIdEndpoint(ps),
		ListPetifiesByIds:     makeListByIdsEndpoint(ps),
		ListPetifiesByOwnerId: makeListByOwnerIdEndpoint(ps),
	}
}

// ================ Endpoint Makers ==================

func makeCreatePetifyEndpoint(ps petifiesservice.PetifesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CreatePetifiesReq)
		result, err := ps.CreatePetify(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapPetifiesAggregateToPetifiesModel(result), nil
	}
}

func makeEditPetifyEndpoint(ps petifiesservice.PetifesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.EditPetifiesReq)
		result, err := ps.EditPetify(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapPetifiesAggregateToPetifiesModel(result), nil
	}
}

func makeGetByIdEndpoint(ps petifiesservice.PetifesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.GetPetifiesByIdReq)
		result, err := ps.GetById(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return mapPetifiesAggregateToPetifiesModel(result), nil
	}
}

func makeListByIdsEndpoint(ps petifiesservice.PetifesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListPetifiesByIdsReq)
		results, err := ps.ListByIds(ctx, req.PetifiesIDs)
		if err != nil {
			return nil, err
		}

		return &models.ManyPetifies{
			Petifies: commonutils.Map2(results, func(p *petifiesaggre.PetifiesAggre) *models.Petifies { return mapPetifiesAggregateToPetifiesModel(p) }),
		}, nil
	}
}

func makeListByOwnerIdEndpoint(ps petifiesservice.PetifesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListPetifiesByOwnerIdReq)
		results, err := ps.ListByOwnerId(ctx, req.OwnerID, req.PageSize, req.AfterID)
		if err != nil {
			return nil, err
		}

		return &models.ManyPetifies{
			Petifies: commonutils.Map2(results, func(p *petifiesaggre.PetifiesAggre) *models.Petifies { return mapPetifiesAggregateToPetifiesModel(p) }),
		}, nil
	}
}

// ================ Mappers ======================

func mapPetifiesAggregateToPetifiesModel(petify *petifiesaggre.PetifiesAggre) *models.Petifies {
	address := petify.GetAddress()

	return &models.Petifies{
		ID:          petify.GetID(),
		OwnerID:     petify.GetOwnerID(),
		Title:       string(petify.GetType()),
		Description: petify.GetDescription(),
		Address: models.Address{
			AddressLineOne: address.AddressLineOne,
			AddressLineTwo: address.AddressLineTwo,
			Street:         address.Street,
			District:       address.District,
			City:           address.City,
			Region:         address.Region,
			PostalCode:     address.PostalCode,
			Country:        address.Country,
			Longitude:      address.Coordinates.GetLongitude(),
			Latitude:       address.Coordinates.GetLatitude(),
		},
		PetName: petify.GetPetName(),
		Images: commonutils.Map2(petify.GetImages(), func(i valueobjects.Image) models.Image {
			return models.Image{URI: i.GetURI(), Description: i.GetDescription()}
		}),
		Status:    string(petify.GetStatus()),
		CreatedAt: petify.GetCreatedAt(),
		UpdatedAt: petify.GetUpdatedAt(),
	}
}
