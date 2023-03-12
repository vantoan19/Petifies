package valueobjects

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/common-utils"
)

var (
	InvalidUploaderIDErr = status.Errorf(codes.InvalidArgument, "invalid uploader ID")
	InvalidFileSizeErr   = status.Errorf(codes.InvalidArgument, "invalid file size")
	InvalidWidthErr      = status.Errorf(codes.InvalidArgument, "invalid width")
	InvalidHeightErr     = status.Errorf(codes.InvalidArgument, "invalid height")
	InvalidDurationErr   = status.Errorf(codes.InvalidArgument, "invalid duration")
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
