package petifiesproposaleventkafka

import (
	"context"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/infrastructure/outbox/utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/petifies-service/cmd"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/publishers"
)

var logger = logging.New("PetifiesService.PetifiesProposalEventPublisher")

type PetifiesProposalEventPublisher struct {
	producer  *producer.KafkaProducer
	eventRepo outbox_repo.EventRepository
}

func NewPetifiesProposalEventPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) publishers.PetifiesProposalEventMessagePublisher {
	return &PetifiesProposalEventPublisher{
		producer:  producer,
		eventRepo: repo,
	}
}

func (p *PetifiesProposalEventPublisher) Publish(ctx context.Context, event models.PetifiesProposalEvent) error {
	logger.Info("Start Publish")

	value, err := event.Serialize()
	if err != nil {
		return err
	}
	payload := models.KafkaMessage{
		Topic:     cmd.Conf.PetifiesProposalEventTopic,
		Partition: 0,
		Offset:    0,
		Key:       []byte("PetifiesProposal"),
		Value:     value,
	}
	outboxEvent := outbox_repo.NewEventWithPayload(payload)

	outboxEvent_, err := p.eventRepo.AddEvent(outboxEvent)
	if err != nil {
		logger.ErrorData("Finish Publish: FAILED", logging.Data{"error": err.Error()})
		return err
	}
	utils.SendOutboxMessageImmediately(*outboxEvent_, payload, p.eventRepo, *p.producer, *logger)

	logger.Info("Finish Publish: SUCCESSFUL")
	return nil
}
