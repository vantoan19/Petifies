package outbox_dispatcher

import (
	"time"

	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
)

type CleanSettings struct {
	EventLifetime time.Duration
}

type cleaner interface {
	DeleteExpiredEvents() error
}

type cleanerImpl struct {
	settings  CleanSettings
	eventRepo outbox_repo.EventRepository
}

func NewCleaner(repo outbox_repo.EventRepository, settings CleanSettings) cleaner {
	return &cleanerImpl{
		settings:  settings,
		eventRepo: repo,
	}
}

func (c *cleanerImpl) DeleteExpiredEvents() error {
	expiryTime := time.Now().Add(-c.settings.EventLifetime)
	return c.eventRepo.DeleteEventsBeforeDatetime(expiryTime)
}
