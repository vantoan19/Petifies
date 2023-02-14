package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	sqlc "github.com/vantoan19/Petifies/server/services/user-service/internal/db"
	userAggre "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/domain/valueobjects"
)

type PostgresRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func New(db *sql.DB) (*PostgresRepository, error) {
	return &PostgresRepository{
		db:      db,
		queries: sqlc.New(db),
	}, nil
}

func FromAggreUserToDBUser(u *userAggre.User) sqlc.User {
	return sqlc.User{
		ID:        u.GetID(),
		Email:     u.GetEmail(),
		Password:  u.GetPassword(),
		FirstName: u.GetName().GetFirstName(),
		LastName:  u.GetName().GetLastName(),
		CreatedAt: u.GetCreatedAt(),
		UpdatedAt: u.GetUpdatedAt(),
	}
}

func FromDBUserToAggreUser(u *sqlc.User) userAggre.User {
	user := userAggre.User{}
	user.SetID(u.ID)
	user.SetEmail(u.Email)
	user.SetPassword(u.Password)
	user.SetName(valueobjects.NewName(u.FirstName, u.LastName))
	user.SetCreatedAt(u.CreatedAt)
	user.SetUpdatedAt(u.UpdatedAt)
	return user
}

func (pr *PostgresRepository) GetByUUID(id uuid.UUID) (userAggre.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := pr.queries.GetUserByID(ctx, id)
	if err != nil {
		return userAggre.User{}, nil
	}
	aggreUser := FromDBUserToAggreUser(&user)

	return aggreUser, nil
}

func (pr *PostgresRepository) GetByEmail(email string) (userAggre.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := pr.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return userAggre.User{}, nil
	}
	aggreUser := FromDBUserToAggreUser(&user)

	return aggreUser, nil
}

func (pr *PostgresRepository) Add(u userAggre.User) (userAggre.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbUser := FromAggreUserToDBUser(&u)

	params := sqlc.CreateUserParams{
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
	}
	addedUsr, err := pr.queries.CreateUser(ctx, params)
	if err != nil {
		return userAggre.User{}, err
	}

	return FromDBUserToAggreUser(&addedUsr), nil
}

func (pr *PostgresRepository) UpdateName(id uuid.UUID, name valueobjects.Name) (userAggre.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	params := sqlc.UpdateUserNameParams{
		ID:        id,
		FirstName: name.GetFirstName(),
		LastName:  name.GetLastName(),
	}
	updatedUser, err := pr.queries.UpdateUserName(ctx, params)
	if err != nil {
		return userAggre.User{}, err
	}

	return FromDBUserToAggreUser(&updatedUser), nil
}

func (pr *PostgresRepository) DeleteByUUID(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pr.queries.DeleteUserByID(ctx, id); err != nil {
		return err
	}

	return nil
}

func (pr *PostgresRepository) DeleteByEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pr.queries.DeleteUserByEmail(ctx, email); err != nil {
		return err
	}

	return nil
}
