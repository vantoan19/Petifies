package entities

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media/valueobjects"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidMediaIDErr          = status.Errorf(codes.InvalidArgument, "invalid media ID")
	InvalidMediaTypeErr        = status.Errorf(codes.InvalidArgument, "invalid media type")
	InvalidMediaMetadataErr    = status.Errorf(codes.InvalidArgument, "invalid media metadata")
	InvalidFilenameErr         = status.Errorf(codes.InvalidArgument, "invalid file name")
	InvalidAudioMetadataErr    = status.Errorf(codes.InvalidArgument, "invalid audio metadata")
	InvalidImageMetadataErr    = status.Errorf(codes.InvalidArgument, "invalid image metadata")
	InvalidVideoMetadataErr    = status.Errorf(codes.InvalidArgument, "invalid video metadata")
	InvalidDocumentMetadataErr = status.Errorf(codes.InvalidArgument, "invalid document metadata")
	DiffSizeErr                = status.Errorf(codes.InvalidArgument, "size metadata is different from data size")
)

// Media represents an uploaded media file.
type Media struct {
	ID        uuid.UUID
	Filename  string
	MediaType valueobjects.MediaType
	Metadata  valueobjects.MediaMetadata
	Data      *bytes.Buffer
	CreatedAt time.Time
}

// Validate checks if the media entity is valid.
func (m *Media) Validate() (errs common.MultiError) {
	if m.ID == uuid.Nil {
		errs = append(errs, InvalidMediaIDErr)
	}
	if m.Filename == "" {
		errs = append(errs, InvalidFilenameErr)
	}
	if m.MediaType == "" {
		errs = append(errs, InvalidMediaTypeErr)
	}
	if errsMeta := m.Metadata.Validate(); errsMeta.Exist() {
		errs = append(errs, errsMeta...)
	}
	if m.MediaType == "audio" && (m.Metadata.GetHeight() > 0 || m.Metadata.GetWidth() > 0 || m.Metadata.GetDuration() == 0) {
		errs = append(errs, InvalidAudioMetadataErr)
	}
	if m.MediaType == "image" && (m.Metadata.GetHeight() == 0 || m.Metadata.GetWidth() == 0 || m.Metadata.GetDuration() > 0) {
		errs = append(errs, InvalidImageMetadataErr)
	}
	if m.MediaType == "video" && (m.Metadata.GetHeight() == 0 || m.Metadata.GetWidth() == 0 || m.Metadata.GetDuration() == 0) {
		errs = append(errs, InvalidVideoMetadataErr)
	}
	if m.MediaType == "document" && (m.Metadata.GetHeight() > 0 || m.Metadata.GetWidth() > 0 || m.Metadata.GetDuration() > 0) {
		errs = append(errs, InvalidDocumentMetadataErr)
	}
	if m.Data.Len() != int(m.Metadata.GetSize()) {
		errs = append(errs, DiffSizeErr)
	}

	return errs
}

// GetURI returns the URI of the media file.
func (m *Media) GetURI() string {
	return fmt.Sprintf("/media/%s", m.ID)
}

// UpdateMetadata updates the metadata of the media file.
func (m *Media) UpdateMetadata(metadata valueobjects.MediaMetadata) {
	m.Metadata = metadata
}
