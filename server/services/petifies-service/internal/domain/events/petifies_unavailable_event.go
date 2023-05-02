package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

type PetifiesUnavailableEvent models.PetifiesEvent

func NewPetifiesUnavailableEvent(petifies *petifiesaggre.PetifiesAggre) PetifiesUnavailableEvent {
	return PetifiesUnavailableEvent{
		ID:          petifies.GetID(),
		OwnerID:     petifies.GetOwnerID(),
		Type:        string(petifies.GetType()),
		Title:       petifies.GetTitle(),
		Description: petifies.GetDescription(),
		Longitude:   petifies.GetAddress().Coordinates.GetLongitude(),
		Latitude:    petifies.GetAddress().Coordinates.GetLatitude(),
		Status:      string(valueobjects.PetifiesUnavailable),
		EventType:   models.PETIFIES_STATUS_CHANGED,
		CreatedAt:   time.Now(),
	}
}
