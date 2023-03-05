package outbox_dispatcher

import (
	"time"

	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
)

type UnlockSettings struct {
	LockDuration time.Duration
}

type unlocker interface {
	UnlockLongtimeLockedEvents() error
}

type unlockerImpl struct {
	settings  UnlockSettings
	eventRepo outbox_repo.EventRepository
}

func NewUnlocker(repo outbox_repo.EventRepository, settings UnlockSettings) unlocker {
	return &unlockerImpl{
		settings:  settings,
		eventRepo: repo,
	}
}

func (u *unlockerImpl) UnlockLongtimeLockedEvents() error {
	unlockTime := time.Now().Add(-u.settings.LockDuration)
	return u.eventRepo.UnlockEventsBeforeDatetime(unlockTime)
}
