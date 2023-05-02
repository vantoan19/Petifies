package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

type PetifiesOnASessionEvent models.PetifiesEvent

func NewPetifiesAvailableEvent(petifies *petifiesaggre.PetifiesAggre) PetifiesOnASessionEvent {
	return PetifiesOnASessionEvent{
		ID:          petifies.GetID(),
		OwnerID:     petifies.GetOwnerID(),
		Type:        string(petifies.GetType()),
		Title:       petifies.GetTitle(),
		Description: petifies.GetDescription(),
		Longitude:   petifies.GetAddress().Coordinates.GetLongitude(),
		Latitude:    petifies.GetAddress().Coordinates.GetLatitude(),
		Status:      string(valueobjects.PetifiesAvailable),
		EventType:   models.PETIFIES_STATUS_CHANGED,
		CreatedAt:   time.Now(),
	}
}
