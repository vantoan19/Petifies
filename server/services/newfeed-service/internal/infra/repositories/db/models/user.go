package models

import (
	"github.com/gocql/gocql"
)

type UserStatus string

const (
	ACTIVE_STATUS   UserStatus = "ACTIVE"
	DEACTIVE_STATUS UserStatus = "DEACTIVE"
	DELETED_STATUS  UserStatus = "DELETED"
)

type User struct {
	ID    gocql.UUID `json:"id"`
	Email string     `json:"email"`
}
