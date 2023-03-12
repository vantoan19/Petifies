package entities

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
)

var (
	ErrEmptyComment           = status.Errorf(codes.InvalidArgument, "comment content cannot be empty")
	ErrBothImageAndVideoExist = status.Errorf(codes.InvalidArgument, "a comment cannot hold both image and video content")
)

type Comment struct {
	ID           uuid.UUID
	PostID       uuid.UUID
	AuthorID     uuid.UUID
	ParentID     uuid.UUID
	IsPostParent bool
	Content      valueobjects.TextContent
	ImageContent valueobjects.ImageContent
	VideoContent valueobjects.VideoContent
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Validate validates the Comment entity and returns any validation errors as a slice of strings.
func (c *Comment) Validate() (errs common.MultiError) {
	if c.ID == uuid.Nil {
		errs = append(errs, ErrEmptyID)
	}

	if c.PostID == uuid.Nil {
		errs = append(errs, ErrEmptyPostID)
	}

	if c.AuthorID == uuid.Nil {
		errs = append(errs, ErrEmptyAuthorID)
	}

	if c.ParentID == uuid.Nil {
		errs = append(errs, ErrEmptyAuthorID)
	}

	if c.Content.IsEmpty() && c.ImageContent.IsEmpty() && c.VideoContent.IsEmpty() {
		errs = append(errs, ErrEmptyComment)
	}

	if !c.ImageContent.IsEmpty() && !c.VideoContent.IsEmpty() {
		errs = append(errs, ErrBothImageAndVideoExist)
	}

	return errs
}
