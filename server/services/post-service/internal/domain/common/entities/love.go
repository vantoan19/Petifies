package entities

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
)

var (
	ErrEmptyPostID                   = status.Errorf(codes.InvalidArgument, "post ID cannot be empty")
	ErrEmptyPostIDAndCommentID       = status.Errorf(codes.InvalidArgument, "at least one of post ID or comment ID is not null")
	ErrBothPostIDAndCommentIDNotNull = status.Errorf(codes.InvalidArgument, "both post ID and comment ID cannot be not null at the same time")
)

type Love struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	CommentID uuid.UUID
	AuthorID  uuid.UUID
	CreatedAt time.Time
}

// Validate validates the Like entity and returns any validation errors as a slice of strings.
func (l *Love) Validate() (errs common.MultiError) {
	if l.ID == uuid.Nil {
		errs = append(errs, ErrEmptyID)
	}
	if l.PostID == uuid.Nil && l.CommentID == uuid.Nil {
		errs = append(errs, ErrEmptyPostIDAndCommentID)
	}
	if l.PostID != uuid.Nil && l.CommentID != uuid.Nil {
		errs = append(errs, ErrBothPostIDAndCommentIDNotNull)
	}
	if l.AuthorID == uuid.Nil {
		errs = append(errs, ErrEmptyAuthorID)
	}
	return errs
}
