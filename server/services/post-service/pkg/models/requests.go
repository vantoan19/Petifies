package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	URL         string
	Description string
}

type Video struct {
	URL         string
	Description string
}

type Love struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	CommentID uuid.UUID
	AuthorID  uuid.UUID
	CreatedAt time.Time
}

type CreatePostReq struct {
	AuthorID    uuid.UUID
	TextContent string
	Images      []Image
	Videos      []Video
}

type CreateCommentReq struct {
	PostID       uuid.UUID
	AuthorID     uuid.UUID
	ParentID     uuid.UUID
	IsPostParent bool
	Content      string
	ImageContent Image
	VideoContent Video
}
