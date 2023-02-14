package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/validateutils"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/domain/valueobjects"
)

var validate = validateutils.NewEnglishValidator()

type User struct {
	// Identifier of the user
	ID uuid.UUID `validate:"required"`
	// Email of the user
	Email string `validate:"required,email,max=300"`
	// Username
	Password string `validate:"required,max=300"`
	// Name
	Name valueobjects.Name `validate:"required"`
	// Timestamp when the user is created
	CreatedAt time.Time `validate:"required"`
	// Timestamp when the user is updated
	UpdatedAt time.Time `validate:"required"`
}

func (u *User) Validate() (errs common.MultiError) {
	errs = append(errs, u.Name.Validate()...)
	errs = append(errs, validate.Struct(u)...)
	return
}
