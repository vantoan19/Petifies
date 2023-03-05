package entities

import (
	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/valueobjects"
)

// Relationship represents a relationship between two users
type Relationship struct {
	ID         uuid.UUID
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
	Type       valueobjects.RelationshipType
}
