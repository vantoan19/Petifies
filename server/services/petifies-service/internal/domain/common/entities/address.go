package entities

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

var (
	ErrEmptyAddressLineOne = status.Errorf(codes.InvalidArgument, "AddressLineOne cannot be empty")
	ErrEmptyStreet         = status.Errorf(codes.InvalidArgument, "Street cannot be empty")
	ErrEmptyCity           = status.Errorf(codes.InvalidArgument, "City cannot be empty")
	ErrEmptyRegion         = status.Errorf(codes.InvalidArgument, "Region cannot be empty")
	ErrEmptyCountry        = status.Errorf(codes.InvalidArgument, "Country cannot be empty")
	ErrEmptyPostalCode     = status.Errorf(codes.InvalidArgument, "PostalCode cannot be empty")
	ErrInvalidCoordinates  = status.Errorf(codes.InvalidArgument, "Coordinates are invalid")
)

type Address struct {
	AddressLineOne string
	AddressLineTwo string
	Street         string
	District       string
	City           string
	Region         string
	PostalCode     string
	Country        string
	Coordinates    valueobjects.Coordinates
}

func (a *Address) Validate() (errs commonutils.MultiError) {
	// if a.AddressLineOne == "" {
	// 	errs = append(errs, ErrEmptyAddressLineOne)
	// }
	if a.Street == "" {
		errs = append(errs, ErrEmptyStreet)
	}
	if a.City == "" {
		errs = append(errs, ErrEmptyCity)
	}
	// if a.Region == "" {
	// 	errs = append(errs, ErrEmptyRegion)
	// }
	if a.Country == "" {
		errs = append(errs, ErrEmptyCountry)
	}
	if a.PostalCode == "" {
		errs = append(errs, ErrEmptyPostalCode)
	}
	if errs_ := a.Coordinates.Validate(); errs_.Exist() {
		errs = append(errs, errs_...)
	}
	return errs
}
