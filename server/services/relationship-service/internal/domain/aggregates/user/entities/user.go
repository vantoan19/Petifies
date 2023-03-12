package entities

import (
	"regexp"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrEmptyId      = status.Errorf(codes.InvalidArgument, "user ID cannot be empty")
	ErrEmptyEmail   = status.Errorf(codes.InvalidArgument, "user email cannot be empty")
	ErrInvalidEmail = status.Errorf(codes.InvalidArgument, "invalid email format")
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
