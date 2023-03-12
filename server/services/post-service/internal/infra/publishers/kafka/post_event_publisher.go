package kafka

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/infrastructure/outbox/utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/post-service/cmd"
)

var logger = logging.New("PostService.PostEventPublisher")

type PostEventPublisher struct {
	producer  *producer.KafkaProducer
	eventRepo outbox_repo.EventRepository
}

func NewPostEventPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) *PostEventPublisher {
	return &PostEventPublisher{
		producer:  producer,
		eventRepo: repo,
	}
}

func (p *PostEventPublisher) Publish(ctx context.Context, event models.PostEvent) error {
	logger.Info("Start Publish")

	value, err := event.Serialize()
	if err != nil {
		return err
	}
	payload := models.KafkaMessage{
		Topic:     cmd.Conf.PostEventTopic,
		Partition: 0,
		Offset:    0,
		Key:       []byte("post"),
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

	outboxEvent_, err := p.eventRepo.AddEvent(outboxEvent)
	if err != nil {
		logger.ErrorData("Finish Publish: FAILED", logging.Data{"error": err.Error()})
		return err
	}
	utils.SendOutboxMessageImmediately(*outboxEvent_, payload, p.eventRepo, *p.producer, *logger)

	logger.Info("Finish Publish: SUCCESSFUL")
	return nil
}
