package valueobjects

import (
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidLongitude = status.Errorf(codes.InvalidArgument, "Longitude should be between -180 and 180")
	ErrInvalidLatitude  = status.Errorf(codes.InvalidArgument, "Latitude should be between -90 and 90")
)

type Coordinates struct {
	longitude float64
	latitude  float64
}

func NewCoordinates(longitude, latitude float64) Coordinates {
	return Coordinates{
		longitude: longitude,
		latitude:  latitude,
	}
}

func (c Coordinates) Validate() (errs common.MultiError) {
	if c.longitude < -180.0 || c.longitude > 180.0 {
		errs = append(errs, ErrInvalidLongitude)
	}
	if c.latitude < -90.0 || c.latitude > 90.0 {
		errs = append(errs, ErrInvalidLatitude)
	}
	return errs
}

func (c Coordinates) GetLongitude() float64 {
	return c.longitude
}

func (c Coordinates) GetLatitude() float64 {
	return c.latitude
}
