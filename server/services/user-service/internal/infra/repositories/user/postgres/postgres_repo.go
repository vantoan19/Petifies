package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	userAggre "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/infra/db/mapper"
	sqlc "github.com/vantoan19/Petifies/server/services/user-service/internal/infra/db/sqlc"
)

type UserRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func New(db *sql.DB) (*UserRepository, error) {
	return &UserRepository{
		db:      db,
		queries: sqlc.New(db),
	}, nil
}

func (pr *UserRepository) GetByUUID(ctx context.Context, id uuid.UUID) (*userAggre.User, error) {
	tCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.GetUserByID(tCtx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("user doesn't exist")
			}
			return err
		}

		sessions, err := q.GetSessionsForUser(tCtx, userDb.ID)
		if err != nil {
			return err
		}

		user, err = mapper.DbModelsToUserAggregate(&userDb, &sessions)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pr *UserRepository) GetByEmail(ctx context.Context, email string) (*userAggre.User, error) {
	tCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var user *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.GetUserByEmail(tCtx, email)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("user doesn't exist")
			}
			return err
		}

		sessions, err := q.GetSessionsForUser(tCtx, userDb.ID)
		if err != nil {
			return err
		}

		user, err = mapper.DbModelsToUserAggregate(&userDb, &sessions)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pr *UserRepository) SaveUser(ctx context.Context, user userAggre.User) (*userAggre.User, error) {
	tCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user_ *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.CreateUser(tCtx, sqlc.CreateUserParams{
			ID:          user.GetID(),
			Email:       user.GetEmail(),
			Password:    user.GetPassword(),
			FirstName:   user.GetName().GetFirstName(),
			LastName:    user.GetName().GetLastName(),
			IsActivated: user.GetIsActivated(),
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					return errors.New("user already exists")
				}
			}
			return err
		}

		sessions := user.GetSessions()
		sessionCreateParams := mapper.EntitySessionsToCreateParams(sessions)
		sessionsDb, err := q.BulkCreateSession(tCtx, *sessionCreateParams)
		if err != nil {
			return err
		}
		user_, err = mapper.DbModelsToUserAggregate(&userDb, &sessionsDb)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user_, nil
}

func (pr *UserRepository) UpdateUser(ctx context.Context, user *userAggre.User) (*userAggre.User, error) {
	tCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user_ *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.UpdateUser(tCtx, sqlc.UpdateUserParams{
			ID:          user.GetID(),
			Email:       user.GetEmail(),
			Password:    user.GetPassword(),
			FirstName:   user.GetName().GetFirstName(),
			LastName:    user.GetName().GetLastName(),
			IsActivated: user.GetIsActivated(),
			UpdatedAt:   time.Now(),
		})
		if err != nil {
			return err
		}

		sessions := user.GetSessions()
		sessionUpsertParams := mapper.EntitySessionsToUpsertParams(sessions)
		sessionsDb, err := q.BulkUpsertSessions(tCtx, *sessionUpsertParams)
		if err != nil {
			return err
		}
		user_, err = mapper.DbModelsToUserAggregate(&userDb, &sessionsDb)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user_, nil
}

func (pr *UserRepository) DeleteByUUID(ctx context.Context, id uuid.UUID) (*userAggre.User, error) {
	tCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.DeleteUserByID(tCtx, id)
		if err != nil {
			return err
		}

		user, err = mapper.DbModelsToUserAggregate(&userDb, &[]sqlc.Session{})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pr *UserRepository) DeleteByEmail(ctx context.Context, email string) (*userAggre.User, error) {
	tCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.DeleteUserByEmail(tCtx, email)
		if err != nil {
			return err
		}

		user, err = mapper.DbModelsToUserAggregate(&userDb, &[]sqlc.Session{})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pr *UserRepository) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)
	if err = fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
