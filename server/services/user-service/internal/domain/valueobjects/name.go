package valueobjects

import (
	"errors"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/validateutils"
)

var validate = validateutils.NewEnglishValidator()

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
		errs = append(errs, errors.New("first name cannot be empty"))
	}
	if n.lastName == "" {
		errs = append(errs, errors.New("last name cannot be empty"))
	}
	if len(n.firstName) > 50 {
		errs = append(errs, errors.New("first name exceeds the maximum length"))
	}
	if len(n.lastName) > 50 {
		errs = append(errs, errors.New("last name exceeds the maximum length"))
	}
	return errs
}

func (n Name) GetFirstName() string {
	return n.firstName
}

func (n Name) GetLastName() string {
	return n.lastName
}
