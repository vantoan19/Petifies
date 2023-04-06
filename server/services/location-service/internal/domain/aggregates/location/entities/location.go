package entities

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location/valueobjects"
)

var (
	ErrEmptyID               = status.Errorf(codes.InvalidArgument, "ID cannot be empty")
	ErrEmptyEntityID         = status.Errorf(codes.InvalidArgument, "EntityID cannot be empty")
	ErrUnknownEntityType     = status.Errorf(codes.InvalidArgument, "Unknown EntityType")
	ErrUnknownLocationStatus = status.Errorf(codes.InvalidArgument, "Unknown LocationStatus")
	ErrInvalidLongitude      = status.Errorf(codes.InvalidArgument, "Longitude should be between -180 and 180")
	ErrInvalidLatitude       = status.Errorf(codes.InvalidArgument, "Latitude should be between -90 and 90")
)

type Location struct {
	ID         uuid.UUID
	Longitude  float64
	Latitude   float64
	Status     valueobjects.LocationStatus
	EntityType valueobjects.EntityType
	EntityID   uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (l *Location) Validate() (errs common.MultiError) {
	if l.ID == uuid.Nil {
		errs = append(errs, ErrEmptyID)
	}
	if l.EntityID == uuid.Nil {
		errs = append(errs, ErrEmptyEntityID)
	}
	if l.EntityType != valueobjects.Petifies {
		errs = append(errs, ErrUnknownEntityType)
	}
	if l.Status != valueobjects.LocationAvailable &&
		l.Status != valueobjects.LocationDeleted &&
		l.Status != valueobjects.LocationUnavailable {
		errs = append(errs, ErrUnknownLocationStatus)
	}
	if l.Longitude < -180 || l.Longitude > 180 {
		errs = append(errs, ErrInvalidLongitude)
	}
	if l.Latitude < -90 || l.Latitude > 90 {
		errs = append(errs, ErrInvalidLatitude)
	}
	return errs
}
