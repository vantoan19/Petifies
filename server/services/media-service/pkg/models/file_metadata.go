package models

import (
	"time"

	"github.com/google/uuid"
)

type FileMetadata struct {
	FileName   string
	MediaType  string
	UploaderId uuid.UUID
	Size       int64
	Width      int
	Height     int
	Duration   time.Duration
}
