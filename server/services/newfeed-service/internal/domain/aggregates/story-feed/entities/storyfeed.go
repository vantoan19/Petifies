package entities

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
)

var (
	InvalidUserIdErr   = status.Errorf(codes.InvalidArgument, "invalid user id")
	InvalidAuthorIdErr = status.Errorf(codes.InvalidArgument, "invalid author id")
	InvalidStoryIdErr  = status.Errorf(codes.InvalidArgument, "invalid story id")
)

type StoryFeed struct {
	UserID    uuid.UUID
	AuthorID  uuid.UUID
	StoryID   uuid.UUID
	CreatedAt time.Time
}

// Validate checks if the media entity is valid.
func (s *StoryFeed) Validate() (errs common.MultiError) {
	if s.UserID == uuid.Nil {
		errs = append(errs, InvalidUserIdErr)
	}
	if s.AuthorID == uuid.Nil {
		errs = append(errs, InvalidAuthorIdErr)
	}
	if s.StoryID == uuid.Nil {
		errs = append(errs, InvalidStoryIdErr)
	}

	return errs
}
