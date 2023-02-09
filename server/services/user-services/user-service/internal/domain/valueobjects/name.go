package valueobjects

import (
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

func (n Name) GetFirstName() string {
	return n.firstName
}

func (n Name) GetLastName() string {
	return n.lastName
}
