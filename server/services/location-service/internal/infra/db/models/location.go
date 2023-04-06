package models

import (
	"time"

	"github.com/google/uuid"
)

type GeoJSON struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

type Location struct {
	ID         uuid.UUID `bson:"id" json:"id"`
	Location   GeoJSON   `bson:"location" json:"location"`
	Status     string    `bson:"status" json:"status"`
	EntityType string    `bson:"entity_type" json:"entity_type"`
	EntityID   uuid.UUID `bson:"entity_id" json:"entity_id"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
