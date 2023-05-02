package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	locationservice "github.com/vantoan19/Petifies/server/services/location-service/internal/application/services/location"
	locationaggre "github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location"
	"github.com/vantoan19/Petifies/server/services/location-service/pkg/models"
)

type LocationEndpoints struct {
	ListNearByLocationsByType endpoint.Endpoint
}

func NewLocationEndpoints(ls locationservice.LocationService) LocationEndpoints {
	return LocationEndpoints{
		ListNearByLocationsByType: makeListNearByLocationsByTypeEndpoint(ls),
	}
}

// ================ Endpoint Makers ==================

func makeListNearByLocationsByTypeEndpoint(ls locationservice.LocationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListNearByLocationsByTypeReq)
		results, err := ls.ListNearByLocationsByType(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.ListNearByLocationsByTypeResp{
			Locations: commonutils.Map2(results, func(l *locationaggre.LocationAggre) *models.Location { return mapLocationAggregateToLocationModel(l) }),
		}, nil
	}
}

// ================ Mappers ======================

func mapLocationAggregateToLocationModel(location *locationaggre.LocationAggre) *models.Location {
	return &models.Location{
		ID:           location.ID(),
		EntityID:     location.EntityID(),
		LocationType: string(location.EntityType()),
	}
}
