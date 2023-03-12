package entities

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	common "github.com/vantoan19/Petifies/server/libs/common-utils"
)

var (
	InvalidUserIdErr   = status.Errorf(codes.InvalidArgument, "invalid user id")
	InvalidAuthorIdErr = status.Errorf(codes.InvalidArgument, "invalid author id")
	InvalidPetifyIdErr = status.Errorf(codes.InvalidArgument, "invalid petify id")
)

type PetifyFeed struct {
	UserID    uuid.UUID
	AuthorID  uuid.UUID
	PetifyID  uuid.UUID
	CreatedAt time.Time
}

// Validate checks if the media entity is valid.
func (p *PetifyFeed) Validate() (errs common.MultiError) {
	if p.UserID == uuid.Nil {
		errs = append(errs, InvalidUserIdErr)
	}
	if p.AuthorID == uuid.Nil {
		errs = append(errs, InvalidAuthorIdErr)
	}
	if p.PetifyID == uuid.Nil {
		errs = append(errs, InvalidPetifyIdErr)
	}

	return errs
}
