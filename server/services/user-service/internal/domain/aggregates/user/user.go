package useragg

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/validateutils"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user/entities"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user/valueobjects"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/utils"
)

var validate = validateutils.GetEnglishValidatorInstance()

type User struct {
	// Root of the user aggregate
	user     *entities.User
	sessions []*entities.Session
}

func New(email, password, firstName, lastName string, isActivated bool) (User, common.MultiError) {
	hashedPW, err := utils.HashPassword(password)
	if err != nil {
		return User{}, []error{err}
	}
	user := entities.User{
		ID:          uuid.New(),
		Email:       email,
		Password:    hashedPW,
		Name:        valueobjects.NewName(firstName, lastName),
		IsActivated: isActivated,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if errs := user.Validate(); errs.Exist() {
		return User{}, errs
	}
	return User{
		user:     &user,
		sessions: make([]*entities.Session, 0),
	}, nil
}

func (u *User) SetUserEntity(ue *entities.User) common.MultiError {
	if errs := ue.Validate(); errs.Exist() {
		return errs
	}
	u.user = ue
	return nil
}

func (u *User) GetID() uuid.UUID {
	return u.user.ID
}

func (u *User) SetID(id uuid.UUID) error {
	if err := validate.Var(id, "required"); err != nil {
		return err
	}
	u.user.ID = id
	return nil
}

func (u *User) GetEmail() string {
	return u.user.Email
}

func (u *User) SetEmail(email string) error {
	if err := validate.Var(email, "required,email,max=300"); err != nil {
		return err
	}
	u.user.Email = email
	return nil
}

func (u *User) GetName() valueobjects.Name {
	return u.user.Name
}

func (u *User) SetName(name valueobjects.Name) error {
	if errs := name.Validate(); errs.Exist() {
		return errors.New(errs.Error())
	}
	u.user.Name = name
	return nil
}

func (u *User) GetPassword() string {
	return u.user.Password
}

func (u *User) SetPassword(password string) error {
	if err := validate.Var(password, "required,max=300"); err != nil {
		return err
	}
	u.user.Password = password
	return nil
}

func (u *User) GetIsActivated() bool {
	return u.user.IsActivated
}

func (u *User) SetIsActivated(isActivated bool) error {
	u.user.IsActivated = isActivated
	return nil
}

func (u *User) GetCreatedAt() time.Time {
	return u.user.CreatedAt
}

func (u *User) SetCreatedAt(t time.Time) error {
	if err := validate.Var(t, "required"); err != nil {
		return err
	}
	u.user.CreatedAt = t
	return nil
}

func (u *User) GetUpdatedAt() time.Time {
	return u.user.UpdatedAt
}

func (u *User) SetUpdatedAt(t time.Time) error {
	if err := validate.Var(t, "required"); err != nil {
		return err
	}
	u.user.UpdatedAt = t
	return nil
}

func (u *User) GetSessions() []*entities.Session {
	return u.sessions
}

func (u *User) GetSessionById(id uuid.UUID) *entities.Session {
	idx := slices.IndexFunc(u.sessions,
		func(s *entities.Session) bool {
			return s.ID.String() == id.String()
		})

	if idx != -1 {
		return u.sessions[idx]
	}
	return nil
}

func (u *User) AddSession(s *entities.Session) error {
	if u.GetSessionById(s.ID) != nil {
		return errors.New("")
	}

	u.sessions = append(u.sessions, s)
	return nil
}
