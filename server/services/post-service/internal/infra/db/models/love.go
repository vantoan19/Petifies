package models

import (
	"time"

	"github.com/google/uuid"
)

type Love struct {
	ID        uuid.UUID `bson:"id"`
	PostID    uuid.UUID `bson:"post_id"`
	CommentID uuid.UUID `bson:"comment_id"`
	AuthorID  uuid.UUID `bson:"author_id"`
	CreatedAt time.Time `bson:"created_at"`
}
