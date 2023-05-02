package listener

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	locationaggre "github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location/valueobjects"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/domain/listeners"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = logging.New("LocationService.KafkaPetifiesEventListener")

type petifiesEventListener struct {
	locationRepo locationaggre.LocationRepository
}

func NewKafkaPetifiesEventListener(locationRepo locationaggre.LocationRepository) listeners.PetifiesEventListener {
	return &petifiesEventListener{
		locationRepo: locationRepo,
	}
}

func (pl *petifiesEventListener) Receive(ctx context.Context, event models.PetifiesEvent) error {
	logger.Info("Start Receive")

	switch event.EventType {
	case models.PETIFIES_CREATED:
		location, err := locationaggre.NewLocation(
			uuid.New(),
			event.ID,
			event.Longitude,
			event.Latitude,
			convertPetifiesTypeToLocationEntityType(event.Type),
			valueobjects.LocationUnavailable,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return err
		}

		_, err = pl.locationRepo.Save(ctx, *location)
		if err != nil {
			return err
		}
		logger.Info("Finish Receive: Successful")
		return nil

	case models.PETIFIES_UPDATED:
		location, err := pl.locationRepo.GetByEntityID(ctx, event.ID)
		if err != nil {
			return err
		}

		location.SetLongitude(event.Longitude)
		location.SetLatitude(event.Latitude)

		_, err = pl.locationRepo.Update(ctx, *location)
		if err != nil {
			return err
		}
		logger.Info("Finish Receive: Successful")
		return nil

	case models.PETIFIES_STATUS_CHANGED:
		location, err := pl.locationRepo.GetByEntityID(ctx, event.ID)
		if err != nil {
			return err
		}

		location.SetStatus(valueobjects.LocationStatus(event.Status))

		_, err = pl.locationRepo.Update(ctx, *location)
		if err != nil {
			return err
		}
		logger.Info("Finish Receive: Successful")
		return nil
	case models.PETIFIES_DELETED:
		logger.Info("Finish Receive: Successful")
		return nil
	default:
		logger.Info("Finish Receive: Unknown PetifiesEvent")
		return status.Errorf(codes.Unknown, "Unknown PetifiesEvent")
	}
}

func convertPetifiesTypeToLocationEntityType(t string) valueobjects.EntityType {
	switch t {
	case "PETIFIES_TYPE_DOG_WALKING":
		return valueobjects.PetifiesDogWalking
	case "PETIFIES_TYPE_CAT_PLAYING":
		return valueobjects.PetifiesCatPlaying
	case "PETIFIES_TYPE_DOG_SITTING":
		return valueobjects.PetifiesDogSitting
	case "PETIFIES_TYPE_CAT_SITTING":
		return valueobjects.PetifiesCatSitting
	case "PETIFIES_TYPE_DOG_ADOPTION":
		return valueobjects.PetifiesDogAdoption
	case "PETIFIES_TYPE_CAT_ADOPTION":
		return valueobjects.PetifiesCatAdoption
	default:
		return valueobjects.UnknownType
	}
}
