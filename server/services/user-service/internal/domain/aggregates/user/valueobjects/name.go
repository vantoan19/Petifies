package valueobjects

import (
	"errors"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
)

var (
	EmptyFirstNameErr        = errors.New("first name cannot be empty")
	EmptyLastNameErr         = errors.New("last name cannot be empty")
	FirstNameExceedLengthErr = errors.New("first name exceeds the maximum length")
	LastNameExceedLengthErr  = errors.New("last name exceeds the maximum length")
)

type Name struct {
	firstName string `validate:"required,omitempty,max=50"`
	lastName  string `validate:"required,omitempty,max=50"`
}

func NewName(firstName, lastName string) Name {
	return Name{
		firstName: firstName,
		lastName:  lastName,
	}
}

func (n Name) Validate() (errs common.MultiError) {
	if n.firstName == "" {
		errs = append(errs, EmptyFirstNameErr)
	}
	if n.lastName == "" {
		errs = append(errs, EmptyLastNameErr)
	}
	if len(n.firstName) > 50 {
		errs = append(errs, FirstNameExceedLengthErr)
	}
	if len(n.lastName) > 50 {
		errs = append(errs, LastNameExceedLengthErr)
	}
	return errs
}

func (n Name) GetFirstName() string {
	return n.firstName
}

func (n Name) GetLastName() string {
	return n.lastName
}
