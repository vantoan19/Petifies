package entities

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

var (
	ErrEmptyId      = errors.New("user ID cannot be empty")
	ErrEmptyEmail   = errors.New("user email cannot be empty")
	ErrInvalidEmail = errors.New("invalid email format")
)

// User represents a user node in the graph
type User struct {
	ID    uuid.UUID
	Email string
}

// Validate checks if a User entity is valid
func (u *User) Validate() error {
	if u.ID == uuid.Nil {
		return ErrEmptyId
	}
	if u.Email == "" {
		return ErrEmptyEmail
	}

	// validate email format
	matched, err := regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", u.Email)
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidEmail
	}

	return nil
}
