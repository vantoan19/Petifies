package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	Images    []Image
	Videos    []Video
	Loves     []Love
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	ID              uuid.UUID
	PostID          uuid.UUID
	AuthorID        uuid.UUID
	ParentID        uuid.UUID
	IsPostParent    bool
	Content         string
	Image           Image
	Video           Video
	Loves           []Love
	SubcommentCount int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Love struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	CommentID uuid.UUID
	AuthorID  uuid.UUID
	CreatedAt time.Time
}

type ListCommentsResp struct {
	Comments []*Comment
}

type ListPostsResp struct {
	Posts []*Post
}
