package mediaaggre

import (
	"bytes"
	"fmt"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media/entities"
	"github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media/valueobjects"
)

// Media represents an aggregate for media.
type Media struct {
	media *entities.Media
}

// NewMedia creates a new media aggregate.
func NewMedia(media *entities.Media) (*Media, error) {
	if err := media.Validate(); err != nil {
		return nil, fmt.Errorf("invalid media: %w", err)
	}
	return &Media{
		media: media,
	}, nil
}

// UpdateMetadata updates the metadata of the media entity.
func (ma *Media) UpdateMetadata(metadata valueobjects.MediaMetadata) error {
	updatedMedia := *ma.media
	updatedMedia.Metadata = metadata

	if err := updatedMedia.Validate(); err != nil {
		return fmt.Errorf("invalid updated media: %w", err)
	}

	ma.media = &updatedMedia
	return nil
}

func (m *Media) GetID() uuid.UUID {
	return m.media.ID
}

func (m *Media) GetType() valueobjects.MediaType {
	return m.media.MediaType
}

func (m *Media) GetMetadata() valueobjects.MediaMetadata {
	return m.media.Metadata
}

func (m *Media) GetFilename() string {
	return m.media.Filename
}

func (m *Media) GetData() *bytes.Buffer {
	return m.media.Data
}

func (m *Media) SetID(id uuid.UUID) error {
	if id == uuid.Nil {
		return entities.InvalidMediaIDErr
	}
	m.media.ID = id
	return nil
}

func (m *Media) SetType(mediaType valueobjects.MediaType) error {
	m.media.MediaType = mediaType
	return nil
}

func (m *Media) SetMetadata(metadata valueobjects.MediaMetadata) error {
	if err := metadata.Validate(); err != nil {
		return err
	}
	m.media.Metadata = metadata
	return nil
}

func (m *Media) SetFilename(filename string) error {
	if filename == "" {
		return entities.InvalidFilenameErr
	}
	m.media.Filename = filename
	return nil
}

func (m *Media) SetData(data *bytes.Buffer) error {
	if data.Len() != int(m.media.Metadata.GetSize()) {
		return entities.DiffSizeErr
	}
	m.media.Data = data
	return nil
}
