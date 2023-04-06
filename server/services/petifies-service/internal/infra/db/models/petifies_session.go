package models

import (
	"time"

	"github.com/google/uuid"
)

type TimeRange struct {
	FromTime time.Time `bson:"from_time"`
	ToTime   time.Time `bson:"to_time"`
}

type PetifiesSession struct {
	ID         uuid.UUID `bson:"id"`
	PetifiesID uuid.UUID `bson:"petifies_id"`
	Time       TimeRange `bson:"time"`
	Status     string    `bson:"status"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}
