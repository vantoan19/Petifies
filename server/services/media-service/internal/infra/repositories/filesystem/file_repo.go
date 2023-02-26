package filesystem

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	mediaaggre "github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media"
)

var (
	FileNotExistErr = errors.New("file does not exist")
)

type MediaRepository struct {
	rootDir string
}

func New(rootDir string) *MediaRepository {
	return &MediaRepository{
		rootDir: rootDir,
	}
}

func (m *MediaRepository) Save(ctx context.Context, media *mediaaggre.Media) (string, error) {
	if err := os.MkdirAll(m.rootDir, 0755); err != nil {
		return "", nil
	}

	filePath := filepath.Join(m.rootDir, media.GetMetadata().GetUploaderID().String(), media.GetFilename())

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = media.GetData().WriteTo(file)
	if err != nil {
		return "", nil
	}

	return filePath, nil
}

func (m *MediaRepository) Remove(ctx context.Context, media *mediaaggre.Media) error {
	filepath := filepath.Join(m.rootDir, media.GetMetadata().GetUploaderID().String(), media.GetFilename())

	if err := os.Remove(filepath); err != nil {
		if os.IsNotExist(err) {
			return FileNotExistErr
		}
		return err
	}

	return nil
}
