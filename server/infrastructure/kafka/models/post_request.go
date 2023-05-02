package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PostStatus string

const (
	POST_CREATED PostStatus = "POST_CREATED"
	POST_UPDATED PostStatus = "POST_UPDATED"
	POST_DELETED PostStatus = "POST_DELETED"
)

type PostEvent struct {
	ID        uuid.UUID  `json:"id"`
	AuthorID  uuid.UUID  `json:"author_id"`
	CreatedAt time.Time  `json:"created_at"`
	Status    PostStatus `json:"status"`
}

func (u *PostEvent) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

func (u *PostEvent) Deserialize(data []byte) error {
	return json.Unmarshal(data, u)
}
