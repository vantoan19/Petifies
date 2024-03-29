package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	userAggre "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/infra/db/mapper"
	sqlc "github.com/vantoan19/Petifies/server/services/user-service/internal/infra/db/sqlc"
)

var logger = logging.New("UserService.UserRepository")

var (
	UserNotExistErr     = errors.New("user doesn't exist")
	UserAlreadyExistErr = errors.New("user already existed")
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
	logger.Info("Start GetByUUID")

	tCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.GetUserByID(tCtx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return status.Errorf(codes.NotFound, UserNotExistErr.Error())
			}
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}

		sessions, err := q.GetSessionsForUser(tCtx, userDb.ID)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}

		user, err = mapper.DbModelsToUserAggregate(&userDb, &sessions)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		return nil
	})
	if err != nil {
		logger.ErrorData("Finish GetByUUID: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetByUUID")
	return user, nil
}

func (pr *UserRepository) GetByEmail(ctx context.Context, email string) (*userAggre.User, error) {
	logger.Info("Start GetByEmail")

	tCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.GetUserByEmail(tCtx, email)
		if err != nil {
			if err == sql.ErrNoRows {
				return status.Errorf(codes.NotFound, UserNotExistErr.Error())
			}
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}

		sessions, err := q.GetSessionsForUser(tCtx, userDb.ID)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}

		user, err = mapper.DbModelsToUserAggregate(&userDb, &sessions)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		return nil
	})
	if err != nil {
		logger.ErrorData("Finish GetByEmail: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetByEmail")
	return user, nil
}

func (pr *UserRepository) ListByIds(ctx context.Context, ids []uuid.UUID) ([]*userAggre.User, error) {
	logger.Info("Start ListByIds")

	tCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var users []*userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		usersDb, err := q.ListUsersByIds(ctx, ids)
		if err != nil {
			if err == sql.ErrNoRows {
				return status.Errorf(codes.NotFound, UserNotExistErr.Error())
			}
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}

		for _, u := range usersDb {
			user, err := mapper.DbModelsToUserAggregate(&u, &[]sqlc.Session{})
			if err != nil {
				return status.Errorf(codes.Internal, err.Error())
			}

			users = append(users, user)
		}
		return nil
	})
	if err != nil {
		logger.ErrorData("Finish ListByIds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListByIds")
	return users, nil
}

func (pr *UserRepository) SaveUser(ctx context.Context, user userAggre.User) (*userAggre.User, error) {
	logger.Info("Start SaveUser")

	tCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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
					return status.Errorf(codes.AlreadyExists, UserAlreadyExistErr.Error())
				}
			}
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}

		sessions := user.GetSessions()
		sessionCreateParams := mapper.EntitySessionsToCreateParams(sessions)
		sessionsDb, err := q.BulkCreateSession(tCtx, *sessionCreateParams)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}
		user_, err = mapper.DbModelsToUserAggregate(&userDb, &sessionsDb)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish SaveUser: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish SaveUser")
	return user_, nil
}

func (pr *UserRepository) UpdateUser(ctx context.Context, user *userAggre.User) (*userAggre.User, error) {
	logger.Info("Start UpdateUser")

	tCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}

		sessions := user.GetSessions()
		sessionUpsertParams := mapper.EntitySessionsToUpsertParams(sessions)
		sessionsDb, err := q.BulkUpsertSessions(tCtx, *sessionUpsertParams)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}
		user_, err = mapper.DbModelsToUserAggregate(&userDb, &sessionsDb)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish UpdateUser: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish UpdateUser")
	return user_, nil
}

func (pr *UserRepository) DeleteByUUID(ctx context.Context, id uuid.UUID) (*userAggre.User, error) {
	logger.Info("Start DeleteByUUID")

	tCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var user *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.DeleteUserByID(tCtx, id)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}

		user, err = mapper.DbModelsToUserAggregate(&userDb, &[]sqlc.Session{})
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		return nil
	})
	if err != nil {
		logger.ErrorData("Finish DeleteByUUID: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish DeleteByUUID")
	return user, nil
}

func (pr *UserRepository) DeleteByEmail(ctx context.Context, email string) (*userAggre.User, error) {
	logger.Info("Start DeleteByEmail")

	tCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var user *userAggre.User
	err := pr.execTx(tCtx, func(q *sqlc.Queries) error {
		userDb, err := q.DeleteUserByEmail(tCtx, email)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("db err: %s", err.Error()))
		}

		user, err = mapper.DbModelsToUserAggregate(&userDb, &[]sqlc.Session{})
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish DeleteByEmail: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish DeleteByEmail")
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
