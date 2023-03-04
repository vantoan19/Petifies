package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	sqlc "github.com/vantoan19/Petifies/server/services/user-service/internal/infra/db/sqlc"
)

type userEventRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func New(db *sql.DB) (outbox_repo.EventRepository, error) {
	return &userEventRepository{
		db:      db,
		queries: sqlc.New(db),
	}, nil
}

func (ur *userEventRepository) AddEvent(event outbox_repo.Event) (*outbox_repo.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	marshaledPayload, err := json.Marshal(event.Payload)
	if err != nil {
		return nil, err
	}
	userEvent, err := ur.queries.CreateUserEvent(ctx, sqlc.CreateUserEventParams{
		ID:          event.ID,
		Payload:     marshaledPayload,
		OutboxState: sqlc.OutboxState(event.OutboxState),
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return dbEventToOutboxEvent(&userEvent)
}

func (ur *userEventRepository) GetEventsByLockerID(lockerID uuid.UUID) ([]*outbox_repo.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	events := make([]*outbox_repo.Event, 0)

	dbEvents, err := ur.queries.GetUserEventByLockerID(ctx, uuid.NullUUID{UUID: lockerID, Valid: true})
	if err != nil {
		return []*outbox_repo.Event{}, err
	}
	for _, e := range dbEvents {
		e_, err := dbEventToOutboxEvent(&e)
		if err != nil {
			return []*outbox_repo.Event{}, err
		}
		events = append(events, e_)
	}

	return events, nil
}

func (ur *userEventRepository) LockStartedEvents(lockerID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := ur.queries.LockStartedEvents(ctx, sqlc.LockStartedEventsParams{
		LockedBy: uuid.NullUUID{UUID: lockerID, Valid: true},
		LockedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}

func (ur *userEventRepository) UpdateEvent(event outbox_repo.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := ur.queries.UpdateEvent(ctx, sqlc.UpdateEventParams{
		ID:          event.ID,
		OutboxState: sqlc.OutboxState(event.OutboxState),
		LockedBy:    utils.UUIDToNullUUID(event.LockedBy),
		LockedAt:    utils.TimeToNullTime(event.LockedAt),
		Error:       utils.StringToNullString(event.Error),
		CompletedAt: utils.TimeToNullTime(event.CompletedAt),
	})
	if err != nil {
		return err
	}

	return nil
}
func (ur *userEventRepository) UnlockEventsByLockerID(lockerID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := ur.queries.UnlockEventsByLockerID(ctx, uuid.NullUUID{UUID: lockerID, Valid: true})
	if err != nil {
		return err
	}

	return nil
}

func (ur *userEventRepository) UnlockEventsBeforeDatetime(t time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := ur.queries.UnlockEventsBeforeDatetime(ctx, sql.NullTime{Time: t, Valid: true})
	if err != nil {
		return err
	}

	return nil
}

func (ur *userEventRepository) DeleteEventsBeforeDatetime(t time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := ur.queries.DeleteEventsBeforeDatetime(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

func dbEventToOutboxEvent(e *sqlc.UserEvent) (*outbox_repo.Event, error) {
	payload := models.KafkaMessage{}
	err := json.Unmarshal(e.Payload, &payload)
	if err != nil {
		return nil, err
	}

	return &outbox_repo.Event{
		ID:          e.ID,
		Payload:     payload,
		OutboxState: outbox_repo.State(e.OutboxState),
		LockedBy:    utils.NullUUIDToUUID(e.LockedBy),
		LockedAt:    utils.NullTimeToTime(e.LockedAt),
		Error:       utils.NullStringToString(e.Error),
		CompletedAt: utils.NullTimeToTime(e.CompletedAt),
		CreatedAt:   e.CreatedAt,
	}, nil
}
