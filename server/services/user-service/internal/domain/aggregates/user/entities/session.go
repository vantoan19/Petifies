package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
)

type Session struct {
	// Identifier of the Session
	ID uuid.UUID `validate:"required"`
	// The owner of the session
	UserID uuid.UUID `validate:"required"`
	// Refresh token
	RefreshToken string `validate:"required"`
	// IP of the client establishing the session
	ClientIP string `validate:"required"`
	// Is the session disabled
	IsDisabled bool `validate:"required"`
	// Timestamp when the session is expired
	ExpiresAt time.Time `validate:"required"`
	// Timestamp when the session is created
	CreatedAt time.Time `validate:"required"`
}

func (s *Session) Validate() (errs common.MultiError) {
	errs = append(errs, validate.Struct(s)...)
	return errs
}

func (s *Session) HasExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
