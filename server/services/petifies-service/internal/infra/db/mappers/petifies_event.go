package mappers

import (
	"encoding/json"

	kafkaModels "github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/models"
)

func DbPetifiesEventToOutboxEvent(event *models.PetifiesEvent) (*outbox_repo.Event, error) {
	var payload kafkaModels.KafkaMessage
	err := json.Unmarshal([]byte(event.Payload), &payload)
	if err != nil {
		return nil, err
	}

	return &outbox_repo.Event{
		ID:          event.ID,
		Payload:     payload,
		OutboxState: outbox_repo.State(event.OutboxState),
		LockedBy:    event.LockedBy,
		LockedAt:    event.LockedAt,
		Error:       event.Error,
		CompletedAt: event.CompletedAt,
		CreatedAt:   event.CreatedAt,
	}, nil
}

func OutboxEventToDbPetifiesEvent(event *outbox_repo.Event) (*models.PetifiesEvent, error) {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return nil, err
	}

	return &models.PetifiesEvent{
		ID:          event.ID,
		Payload:     string(payload),
		OutboxState: models.OutboxState(event.OutboxState),
		LockedBy:    event.LockedBy,
		LockedAt:    event.LockedAt,
		Error:       event.Error,
		CompletedAt: event.CompletedAt,
		CreatedAt:   event.CreatedAt,
	}, nil
}
