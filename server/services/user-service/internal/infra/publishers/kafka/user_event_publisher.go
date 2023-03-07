package kafka

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/user-service/cmd"
)

var logger = logging.New("UserService.UserEventPublisher")

type UserEventPublisher struct {
	producer  *producer.KafkaProducer
	eventRepo outbox_repo.EventRepository
}

func NewUserEventPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) *UserEventPublisher {
	return &UserEventPublisher{
		producer:  producer,
		eventRepo: repo,
	}
}

func (u *UserEventPublisher) Publish(ctx context.Context, event models.UserEvent) error {
	value, err := event.Serialize()
	if err != nil {
		return nil
	}
	payload := models.KafkaMessage{
		Topic:     cmd.Conf.UserRequestTopic,
		Partition: 0,
		Offset:    0,
		Key:       []byte("user"),
		Value:     value,
	}
	outboxEvent := outbox_repo.Event{
		ID:          uuid.New(),
		Payload:     payload,
		OutboxState: outbox_repo.StartedState,
		LockedBy:    nil,
		LockedAt:    nil,
		Error:       nil,
		CompletedAt: nil,
		CreatedAt:   time.Now(),
	}

	outboxEvent_, err := u.eventRepo.AddEvent(outboxEvent)
	if err != nil {
		return err
	}

	// Lock event and publish immediately
	lockerID := uuid.New()
	now := time.Now()
	outboxEvent_.LockedBy = &lockerID
	outboxEvent_.LockedAt = &now
	err = u.eventRepo.UpdateEvent(*outboxEvent_)
	if err != nil {
		logger.WarningData("Executing Publish: error at setting lock, publishing event later by outbox", logging.Data{"error": err.Error()})
		return nil
	}

	_, err = u.producer.SendMessage(&payload)
	if err != nil {
		outboxEvent_.LockedBy = nil
		outboxEvent_.LockedAt = nil
		errMsg := err.Error()
		outboxEvent_.Error = &errMsg
		dbErr := u.eventRepo.UpdateEvent(*outboxEvent_)
		if dbErr != nil {
			logger.WarningData("Executing Publish: error at updating event, publishing event later by outbox", logging.Data{"error": dbErr.Error()})
			return nil
		}
		logger.WarningData("Executing Publish: error at publishing msg, publishing event later by outbox", logging.Data{"error": err.Error()})
	} else {
		outboxEvent_.LockedBy = nil
		outboxEvent_.LockedAt = nil
		outboxEvent_.OutboxState = outbox_repo.CompletedState
		outboxEvent_.CompletedAt = &now
		dbErr := u.eventRepo.UpdateEvent(*outboxEvent_)
		if dbErr != nil {
			logger.WarningData("Executing Publish: error at updating event", logging.Data{"error": dbErr.Error()})
		}
	}

	return nil
}
