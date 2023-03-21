package models

import (
	"time"

	"github.com/google/uuid"
)

type Love struct {
	ID           uuid.UUID `bson:"id"`
	TargetID     uuid.UUID `bson:"target_id"`
	IsPostTarget bool      `bson:"is_post_target"`
	AuthorID     uuid.UUID `bson:"author_id"`
	CreatedAt    time.Time `bson:"created_at"`
}
