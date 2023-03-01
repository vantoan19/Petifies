package mapper

import (
	"errors"

	useragg "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user/entities"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user/valueobjects"
	db "github.com/vantoan19/Petifies/server/services/user-service/internal/infra/db/sqlc"
)

func DbUserToEntity(u *db.User) entities.User {
	return entities.User{
		ID:          u.ID,
		Email:       u.Email,
		Password:    u.Password,
		Name:        valueobjects.NewName(u.FirstName, u.LastName),
		IsActivated: u.IsActivated,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func EntityUserToDb(u entities.User) *db.User {
	return &db.User{
		ID:          u.ID,
		Email:       u.Email,
		Password:    u.Password,
		FirstName:   u.Name.GetFirstName(),
		LastName:    u.Name.GetLastName(),
		IsActivated: u.IsActivated,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func DbModelsToUserAggregate(u *db.User, ss *[]db.Session) (*useragg.User, error) {
	user := &useragg.User{}

	if errs := user.SetUserEntity(DbUserToEntity(u)); errs.Exist() {
		return nil, errors.New(errs.Error())
	}

	for _, s := range *ss {
		if err := user.AddSession(DbSessionToEntity(&s)); err != nil {
			return nil, err
		}
	}

	return user, nil
}
