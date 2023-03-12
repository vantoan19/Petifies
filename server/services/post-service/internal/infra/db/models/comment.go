package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID           uuid.UUID `bson:"id"`
	PostID       uuid.UUID `bson:"post_id"`
	AuthorID     uuid.UUID `bson:"author_id"`
	ParentID     uuid.UUID `bson:"parent_id"`
	IsPostParent bool      `bson:"is_post_parent"`
	Content      string    `bson:"content"`
	ImageContent Image     `bson:"image"`
	VideoContent Video     `bson:"video"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}
