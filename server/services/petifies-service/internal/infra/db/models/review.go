package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID         uuid.UUID `bson:"id"`
	PetifiesID uuid.UUID `bson:"petifies_id"`
	AuthorID   uuid.UUID `bson:"author_id"`
	Review     string    `bson:"review"`
	Image      Image     `bson:"image"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}
