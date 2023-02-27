package valueobjects

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/libs/common-utils"
)

var (
	InvalidUploaderIDErr = errors.New("invalid uploader ID")
	InvalidFileSizeErr   = errors.New("invalid file size")
	InvalidWidthErr      = errors.New("invalid width")
	InvalidHeightErr     = errors.New("invalid height")
	InvalidDurationErr   = errors.New("invalid duration")
)

// MediaMetadata contains metadata fields about the uploaded media file.
type MediaMetadata struct {
	uploaderID uuid.UUID
	size       int64
	width      int
	height     int
	duration   time.Duration
}

func NewMediaMetadata(uploaderID uuid.UUID, size int64, width, height int, duration time.Duration) MediaMetadata {
	return MediaMetadata{
		uploaderID: uploaderID,
		size:       size,
		width:      width,
		height:     height,
		duration:   duration,
	}
}

// Validate checks if the metadata is valid.
func (m *MediaMetadata) Validate() (errs common.MultiError) {

	if m.uploaderID == uuid.Nil {
		errs = append(errs, InvalidUploaderIDErr)
	}
	if m.size <= 0 {
		errs = append(errs, InvalidFileSizeErr)
	}
	if m.width < 0 {
		errs = append(errs, InvalidWidthErr)
	}
	if m.height < 0 {
		errs = append(errs, InvalidHeightErr)
	}
	if m.duration < 0 {
		errs = append(errs, InvalidDurationErr)
	}

	return errs
}

func (m MediaMetadata) GetUploaderID() uuid.UUID {
	return m.uploaderID
}

func (m MediaMetadata) GetSize() int64 {
	return m.size
}

func (m MediaMetadata) GetWidth() int {
	return m.width
}

func (m MediaMetadata) GetHeight() int {
	return m.height
}

func (m MediaMetadata) GetDuration() time.Duration {
	return m.duration
}
