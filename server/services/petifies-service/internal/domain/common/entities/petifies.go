package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

type Petifies struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Type        valueobjects.PetifiesType
	Title       string
	Description string
	PetName     string
	Images      []valueobjects.Image
	Status      valueobjects.PetifiesStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	ErrOwnerIDRequired         = status.Errorf(codes.InvalidArgument, "OwnerID is required")
	ErrTitleRequired           = status.Errorf(codes.InvalidArgument, "Title is required")
	ErrTitleExceedsLimit       = status.Errorf(codes.InvalidArgument, "Title cannot exceed 500 characters")
	ErrDescriptionRequired     = status.Errorf(codes.InvalidArgument, "Description is required")
	ErrDescriptionExceedsLimit = status.Errorf(codes.InvalidArgument, "Description cannot exceed 5000 characters")
	ErrPetNameRequired         = status.Errorf(codes.InvalidArgument, "PetName is required")
	ErrPetNameExceedsLimit     = status.Errorf(codes.InvalidArgument, "PetName cannot exceed 50 characters")
	ErrImagesRequired          = status.Errorf(codes.InvalidArgument, "At least one image is required")
	ErrCreatedAtRequired       = status.Errorf(codes.InvalidArgument, "CreatedAt is required")
	ErrUpdatedAtRequired       = status.Errorf(codes.InvalidArgument, "UpdatedAt is required")
	ErrUnknownPetifiesStatus   = status.Errorf(codes.InvalidArgument, "Unknown Petifies status")
	ErrUnknownPetifiesType     = status.Errorf(codes.InvalidArgument, "Unknown Petifies type")
)

func (p *Petifies) Validate() (errs commonutils.MultiError) {
	if p.ID == uuid.Nil {
		errs = append(errs, ErrIDRequired)
	}
	if p.OwnerID == uuid.Nil {
		errs = append(errs, ErrOwnerIDRequired)
	}
	switch p.Type {
	case valueobjects.PetifiesCatAdoption,
		valueobjects.PetifiesCatPlaying,
		valueobjects.PetifiesCatSitting,
		valueobjects.PetifiesDogAdoption,
		valueobjects.PetifiesDogSitting,
		valueobjects.PetifiesDogWalking:
		// do nothing
	default:
		errs = append(errs, ErrUnknownPetifiesType)
	}
	if len(strings.TrimSpace(p.Title)) == 0 {
		errs = append(errs, ErrTitleRequired)
	}
	if len(p.Title) > 500 {
		errs = append(errs, ErrTitleExceedsLimit)
	}
	if len(strings.TrimSpace(p.Description)) == 0 {
		errs = append(errs, ErrDescriptionRequired)
	}
	if len(p.Description) > 5000 {
		errs = append(errs, ErrDescriptionExceedsLimit)
	}
	if len(strings.TrimSpace(p.PetName)) == 0 {
		errs = append(errs, ErrPetNameRequired)
	}
	if len(p.PetName) > 50 {
		errs = append(errs, ErrPetNameExceedsLimit)
	}
	if len(p.Images) == 0 {
		errs = append(errs, ErrImagesRequired)
	}
	for _, image := range p.Images {
		if errs_ := image.Validate(); errs_ != nil {
			errs = append(errs, errs_...)
			break
		}
	}
	switch p.Status {
	case valueobjects.PetifiesUnavailable,
		valueobjects.PetifiesAvailable,
		valueobjects.PetifiesDeleted:
		// do nothing
	default:
		errs = append(errs, ErrUnknownPetifiesStatus)
	}
	if p.CreatedAt.IsZero() {
		errs = append(errs, ErrCreatedAtRequired)
	}
	if p.UpdatedAt.IsZero() {
		errs = append(errs, ErrUpdatedAtRequired)
	}

	return errs
}
