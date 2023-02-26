package entities

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media/valueobjects"
)

var (
	InvalidMediaIDErr          = errors.New("invalid media ID")
	InvalidMediaTypeErr        = errors.New("invalid media type")
	InvalidMediaMetadataErr    = errors.New("invalid media metadata")
	InvalidFilenameErr         = errors.New("invalid file name")
	InvalidAudioMetadataErr    = errors.New("invalid audio metadata")
	InvalidImageMetadataErr    = errors.New("invalid image metadata")
	InvalidVideoMetadataErr    = errors.New("invalid video metadata")
	InvalidDocumentMetadataErr = errors.New("invalid document metadata")
	DiffSizeErr                = errors.New("size metadata is different from data size")
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
