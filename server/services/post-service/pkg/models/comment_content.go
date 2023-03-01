package models

import (
	"github.com/google/uuid"
)

type CommentContent struct {
	PostID       uuid.UUID
	AuthorID     uuid.UUID
	ParentID     uuid.UUID
	IsPostParent bool
	Content      string
	ImageContent Image
	VideoContent Video
}
