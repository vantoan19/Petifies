package models

import (
	"time"

	"github.com/google/uuid"
)

type LoginResp struct {
	SessionID             uuid.UUID
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	User                  User
}

type VerifyTokenResp struct {
	UserID string
}

type RefreshTokenResp struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
}

type ListUsersByIdsResp struct {
	Users []*User
}
