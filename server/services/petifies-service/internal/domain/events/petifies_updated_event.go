package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
)

type PetifiesUpdatedEvent models.PetifiesEvent

func NewPetifiesUpdatedEvent(petifies *petifiesaggre.PetifiesAggre) PetifiesUpdatedEvent {
	return PetifiesUpdatedEvent{
		ID:          petifies.GetID(),
		OwnerID:     petifies.GetOwnerID(),
		Type:        string(petifies.GetType()),
		Title:       petifies.GetTitle(),
		Description: petifies.GetDescription(),
		Longitude:   petifies.GetAddress().Coordinates.GetLongitude(),
		Latitude:    petifies.GetAddress().Coordinates.GetLatitude(),
		Status:      string(petifies.GetStatus()),
		EventType:   models.PETIFIES_UPDATED,
		CreatedAt:   time.Now(),
	}
}
