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

var (
	ErrInvalidPetifiesID  = status.Errorf(codes.InvalidArgument, "invalid Petifies ID")
	ErrInvalidAuthorID    = status.Errorf(codes.InvalidArgument, "invalid author ID")
	ErrInvalidImage       = status.Errorf(codes.InvalidArgument, "invalid image")
	ErrReviewRequired     = status.Errorf(codes.InvalidArgument, "Review is required")
	ErrReviewExceedsLimit = status.Errorf(codes.InvalidArgument, "Review cannot exceed 5000 characters")
)

type Review struct {
	ID         uuid.UUID
	PetifiesID uuid.UUID
	AuthorID   uuid.UUID
	Review     string
	Image      valueobjects.Image
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (r *Review) Validate() (errs commonutils.MultiError) {
	if r.ID == uuid.Nil {
		errs = append(errs, ErrIDRequired)
	}
	if r.PetifiesID == uuid.Nil {
		errs = append(errs, ErrInvalidPetifiesID)
	}
	if r.AuthorID == uuid.Nil {
		errs = append(errs, ErrInvalidAuthorID)
	}
	if len(strings.TrimSpace(r.Review)) == 0 {
		errs = append(errs, ErrReviewRequired)
	}
	if len(r.Review) > 5000 {
		errs = append(errs, ErrReviewExceedsLimit)
	}
	if errs_ := r.Image.Validate(); errs_ != nil {
		errs = append(errs, errs_...)
	}

	return errs
}
