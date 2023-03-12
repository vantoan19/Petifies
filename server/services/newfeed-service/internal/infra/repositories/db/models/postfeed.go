package models

import (
	"time"

	"github.com/gocql/gocql"
)

type PostFeedStatus string

type PostFeed struct {
	UserID    gocql.UUID
	AuthorID  gocql.UUID
	PostID    gocql.UUID
	CreatedAt time.Time
}
