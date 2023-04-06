package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
)

type PetifiesCreatedEvent models.PetifiesEvent

func NewPetifiesCreatedEvent(petifies *petifiesaggre.PetifiesAggre) PetifiesCreatedEvent {
	return PetifiesCreatedEvent{
		ID:          petifies.GetID(),
		OwnerID:     petifies.GetOwnerID(),
		Type:        string(petifies.GetType()),
		Title:       petifies.GetTitle(),
		Description: petifies.GetDescription(),
		Longitude:   petifies.GetAddress().Coordinates.GetLongitude(),
		Latitude:    petifies.GetAddress().Coordinates.GetLatitude(),
		Status:      string(petifies.GetStatus()),
		EventType:   models.PETIFIES_CREATED,
		CreatedAt:   time.Now(),
	}
}
