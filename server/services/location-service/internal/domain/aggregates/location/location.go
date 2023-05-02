package locationaggre

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location/entities"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location/valueobjects"
)

type LocationAggre struct {
	location *entities.Location
}

func NewLocation(
	id,
	entityID uuid.UUID,
	longitude,
	latitude float64,
	entityType valueobjects.EntityType,
	status valueobjects.LocationStatus,
	createdAt,
	updatedAt time.Time,
) (*LocationAggre, error) {
	locationEntity := &entities.Location{
		ID:         id,
		EntityID:   entityID,
		Longitude:  longitude,
		Latitude:   latitude,
		EntityType: entityType,
		Status:     status,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	if errs := locationEntity.Validate(); errs.Exist() {
		return nil, errs[0]
	}

	return &LocationAggre{
		location: locationEntity,
	}, nil
}

// ======== Aggregate Root Getter ========

func (l *LocationAggre) ID() uuid.UUID {
	return l.location.ID
}

func (l *LocationAggre) EntityID() uuid.UUID {
	return l.location.EntityID
}

func (l *LocationAggre) Longitude() float64 {
	return l.location.Longitude
}

func (l *LocationAggre) SetLongitude(longitude float64) error {
	if longitude < -180 || longitude > 180 {
		return entities.ErrInvalidLongitude
	}
	l.location.Longitude = longitude
	return nil
}

func (l *LocationAggre) Latitude() float64 {
	return l.location.Latitude
}

func (l *LocationAggre) SetLatitude(latitude float64) error {
	if latitude < -90 || latitude > 90 {
		return entities.ErrInvalidLongitude
	}
	l.location.Latitude = latitude
	return nil
}

func (l *LocationAggre) EntityType() valueobjects.EntityType {
	return l.location.EntityType
}

func (l *LocationAggre) Status() valueobjects.LocationStatus {
	return l.location.Status
}

func (l *LocationAggre) SetStatus(status valueobjects.LocationStatus) error {
	switch status {
	case valueobjects.LocationAvailable,
		valueobjects.LocationDeleted,
		valueobjects.LocationUnavailable:
		l.location.Status = status
		return nil
	default:
		return entities.ErrUnknownLocationStatus
	}
}

func (l *LocationAggre) CreatedAt() time.Time {
	return l.location.CreatedAt
}

func (l *LocationAggre) UpdatedAt() time.Time {
	return l.location.UpdatedAt
}
