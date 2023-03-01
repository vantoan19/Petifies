package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	URL         string `bson:"url"`
	Description string `bson:"description"`
}

type Video struct {
	URL         string `bson:"url"`
	Description string `bson:"description"`
}

type Post struct {
	ID          uuid.UUID `bson:"id"`
	AuthorID    uuid.UUID `bson:"author_id"`
	TextContent string    `bson:"text_content"`
	Images      []Image   `bson:"images"`
	Videos      []Video   `bson:"videos"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
