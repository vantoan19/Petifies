package outbox_dispatcher

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
)

// Placeholder in case there are settings needed in the future
type PublishSettings struct{}

type publisher interface {
	PublishEvents() error
}

type publisherImpl struct {
	settings  PublishSettings
	producer  producer.KafkaProducer
	eventRepo outbox_repo.EventRepository
	lockerID  uuid.UUID
}

func NewPublisher(producer producer.KafkaProducer, repo outbox_repo.EventRepository, settings PublishSettings) publisher {
	return &publisherImpl{
		settings:  settings,
		producer:  producer,
		eventRepo: repo,
		lockerID:  uuid.New(),
	}
}

func (p *publisherImpl) PublishEvents() (err error) {
	err = p.eventRepo.LockStartedEvents(p.lockerID)
	defer func() {
		if err_ := p.eventRepo.UnlockEventsByLockerID(p.lockerID); err_ != nil {
			err = err_
		}
	}()
	if err != nil {
		return err
	}

	events, err := p.eventRepo.GetEventsByLockerID(p.lockerID)
	if err != nil {
		return err
	}
	for _, e := range events {
		_, err := p.producer.SendMessage(&e.Payload)
		if err != nil {
			e.LockedBy = nil
			e.LockedAt = nil
			errMsg := err.Error()
			e.Error = &errMsg

			dbErr := p.eventRepo.UpdateEvent(*e)
			if dbErr != nil {
				return dbErr
			}
			return fmt.Errorf("An error occured when sending event: %w", err)
		}

		now := time.Now()
		e.LockedBy = nil
		e.LockedAt = nil
		e.Error = nil
		e.OutboxState = outbox_repo.CompletedState
		e.CompletedAt = &now

		err = p.eventRepo.UpdateEvent(*e)
		if err != nil {
			return err
		}
	}

	return nil
}
