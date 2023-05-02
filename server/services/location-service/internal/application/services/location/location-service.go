package locationservice

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	locationaggre "github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location/valueobjects"
	locationmongo "github.com/vantoan19/Petifies/server/services/location-service/internal/infra/repositories/location/mongo"
	"github.com/vantoan19/Petifies/server/services/location-service/pkg/models"
)

var logger = logging.New("LocationService.LocationSvc")

type locationService struct {
	mongoClient  *mongo.Client
	locationRepo locationaggre.LocationRepository
}

type LocationConfiguration func(ps *locationService) error

type LocationService interface {
	ListNearByLocationsByType(ctx context.Context, req *models.ListNearByLocationsByTypeReq) ([]*locationaggre.LocationAggre, error)
}

func NewLocationService(cfgs ...LocationConfiguration) (LocationService, error) {
	ls := &locationService{}
	for _, cfg := range cfgs {
		if err := cfg(ls); err != nil {
			return nil, err
		}
	}
	return ls, nil
}

func WithMongoLocationRepository(client *mongo.Client) LocationConfiguration {
	return func(ls *locationService) error {
		repo := locationmongo.New(client)
		ls.locationRepo = repo
		return nil
	}
}

func (ls *locationService) ListNearByLocationsByType(ctx context.Context, req *models.ListNearByLocationsByTypeReq) ([]*locationaggre.LocationAggre, error) {
	logger.Info("Start ListNearByLocationsByType")

	locations, err := ls.locationRepo.FindNearbyLocationsByEntityType(
		ctx,
		req.Longitude,
		req.Latitude,
		req.Radius,
		valueobjects.EntityType(req.LocationType),
		req.PageSize,
		req.Offset,
	)
	if err != nil {
		return nil, err
	}

	logger.Info("Finish ListNearByLocationsByType: Successful")
	return locations, nil
}
