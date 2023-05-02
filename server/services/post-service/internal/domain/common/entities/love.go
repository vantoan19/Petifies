package entities

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
)

var (
	ErrEmptyPostID   = status.Errorf(codes.InvalidArgument, "post ID cannot be empty")
	ErrEmptyTargetID = status.Errorf(codes.InvalidArgument, "target ID cannot be empty")
)

type Love struct {
	ID           uuid.UUID
	TargetID     uuid.UUID
	IsPostTarget bool
	AuthorID     uuid.UUID
	CreatedAt    time.Time
}

// Validate validates the Like entity and returns any validation errors as a slice of strings.
func (l *Love) Validate() (errs common.MultiError) {
	if l.ID == uuid.Nil {
		errs = append(errs, ErrEmptyID)
	}
	if l.TargetID == uuid.Nil {
		errs = append(errs, ErrEmptyTargetID)
	}
	if l.AuthorID == uuid.Nil {
		errs = append(errs, ErrEmptyAuthorID)
	}
	return errs
}
