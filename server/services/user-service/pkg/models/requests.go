package models

import "github.com/google/uuid"

type CreateUserReq struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type LoginReq struct {
	Email    string
	Password string
}

type VerifyTokenReq struct {
	Token string
}

type RefreshTokenReq struct {
	RefreshToken string
}

type GetUserReq struct {
	ID uuid.UUID
}
