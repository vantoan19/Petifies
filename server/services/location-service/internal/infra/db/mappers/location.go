package mappers

import (
	locationaggre "github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location/valueobjects"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/infra/db/models"
)

func DbModelToLocationAggregate(l *models.Location) (*locationaggre.LocationAggre, error) {
	location, err := locationaggre.NewLocation(
		l.ID,
		l.EntityID,
		l.Location.Coordinates[0],
		l.Location.Coordinates[1],
		valueobjects.EntityType(l.EntityType),
		valueobjects.LocationStatus(l.Status),
		l.CreatedAt,
		l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func AggregateLocationToDbLocation(l *locationaggre.LocationAggre) *models.Location {
	return &models.Location{
		ID: l.ID(),
		Location: models.GeoJSON{
			Type:        "Point",
			Coordinates: []float64{l.Longitude(), l.Latitude()},
		},
		Status:     string(l.Status()),
		EntityType: string(l.EntityType()),
		EntityID:   l.EntityID(),
		CreatedAt:  l.CreatedAt(),
		UpdatedAt:  l.UpdatedAt(),
	}
}
