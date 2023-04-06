package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
)

type PetifiesDeletedEvent models.PetifiesEvent

func NewPetifiesDeletedEvent(petifies *petifiesaggre.PetifiesAggre) PetifiesDeletedEvent {
	return PetifiesDeletedEvent{
		ID:          petifies.GetID(),
		OwnerID:     petifies.GetOwnerID(),
		Type:        string(petifies.GetType()),
		Title:       petifies.GetTitle(),
		Description: petifies.GetDescription(),
		Longitude:   petifies.GetAddress().Coordinates.GetLongitude(),
		Latitude:    petifies.GetAddress().Coordinates.GetLatitude(),
		Status:      string(petifies.GetStatus()),
		EventType:   models.PETIFIES_DELETED,
		CreatedAt:   time.Now(),
	}
}
