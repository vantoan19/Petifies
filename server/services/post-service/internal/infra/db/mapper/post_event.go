package mapper

import (
	"encoding/json"

	kafkaModels "github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
)

func DbPostEventToOutboxEvent(event *models.PostEvent) (*outbox_repo.Event, error) {
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

func OutboxEventToDbPostEvent(event *outbox_repo.Event) (*models.PostEvent, error) {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return nil, err
	}

	return &models.PostEvent{
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
