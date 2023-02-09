package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/vantoan19/Petifies/server/libs/validateutils"
	"github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/domain/entities"
	"github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/domain/valueobjects"
)

var validate = validateutils.NewEnglishValidator()

type User struct {
	// Root of the user aggregate
	user *entities.User
}

func New(email, password, firstName, lastName string) (User, error) {
	user := entities.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		Name:      valueobjects.NewName(firstName, lastName),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := validate.Struct(user); err != nil {
		return User{}, err
	}

	return User{
		user: &user,
	}, nil
}

func (u *User) GetID() uuid.UUID {
	return u.user.ID
}

func (u *User) SetID(id uuid.UUID) {
	if u.user == nil {
		u.user = &entities.User{}
	}
	u.user.ID = id
}

func (u *User) GetEmail() string {
	return u.user.Email
}

func (u *User) SetEmail(email string) {
	if u.user == nil {
		u.user = &entities.User{}
	}
	u.user.Email = email
}

func (u *User) GetName() valueobjects.Name {
	return u.user.Name
}

func (u *User) SetName(name valueobjects.Name) {
	if u.user == nil {
		u.user = &entities.User{}
	}
	u.user.Name = name
}

func (u *User) GetPassword() string {
	return u.user.Password
}

func (u *User) SetPassword(password string) {
	if u.user == nil {
		u.user = &entities.User{}
	}
	u.user.Password = password
}

func (u *User) GetCreatedAt() time.Time {
	return u.user.CreatedAt
}

func (u *User) SetCreatedAt(t time.Time) {
	if u.user == nil {
		u.user = &entities.User{}
	}
	u.user.CreatedAt = t
}

func (u *User) GetUpdatedAt() time.Time {
	return u.user.UpdatedAt
}

func (u *User) SetUpdatedAt(t time.Time) {
	if u.user == nil {
		u.user = &entities.User{}
	}
	u.user.UpdatedAt = t
}
