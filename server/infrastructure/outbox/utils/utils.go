package utils

import (
	"time"

	"github.com/google/uuid"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

func SendOutboxMessageImmediately(
	outboxEvent outbox_repo.Event,
	payload models.KafkaMessage,
	repo outbox_repo.EventRepository,
	producer producer.KafkaProducer,
	logger logging.Logger,
) {
	// Lock event and publish immediately
	lockerID := uuid.New()
	now := time.Now()
	outboxEvent.LockedBy = &lockerID
	outboxEvent.LockedAt = &now
	err := repo.UpdateEvent(outboxEvent)
	if err != nil {
		logger.WarningData("Executing Publish: error at setting lock, publishing event later by outbox", logging.Data{"error": err.Error()})
	}

	_, err = producer.SendMessage(&payload)
	if err != nil {
		outboxEvent.LockedBy = nil
		outboxEvent.LockedAt = nil
		errMsg := err.Error()
		outboxEvent.Error = &errMsg
		dbErr := repo.UpdateEvent(outboxEvent)
		if dbErr != nil {
			logger.WarningData("Executing Publish: error at updating event, publishing event later by outbox", logging.Data{"error": dbErr.Error()})
		}
		logger.WarningData("Executing Publish: error at publishing msg, publishing event later by outbox", logging.Data{"error": err.Error()})
	} else {
		outboxEvent.LockedBy = nil
		outboxEvent.LockedAt = nil
		outboxEvent.OutboxState = outbox_repo.CompletedState
		outboxEvent.CompletedAt = &now
		dbErr := repo.UpdateEvent(outboxEvent)
		if dbErr != nil {
			logger.WarningData("Executing Publish: error at updating event", logging.Data{"error": dbErr.Error()})
		}
	}
}
