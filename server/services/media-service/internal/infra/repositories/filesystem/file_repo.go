package filesystem

import (
	"context"
	"os"
	"path/filepath"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mediaaggre "github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media"
)

var (
	FileNotExistErr = status.Errorf(codes.Internal, "file does not exist")
)

type MediaRepository struct {
	rootDir string
}

func New(rootDir string) (*MediaRepository, error) {
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &MediaRepository{
		rootDir: rootDir,
	}, nil
}

func (m *MediaRepository) Save(ctx context.Context, media *mediaaggre.Media) (string, error) {
	directory := filepath.Join(m.rootDir, media.GetMetadata().GetUploaderID().String())
	if err := os.MkdirAll(directory, 0755); err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}

	filePath := filepath.Join(directory, media.GetFilename())

	file, err := os.Create(filePath)
	if err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}
	defer file.Close()

	_, err = media.GetData().WriteTo(file)
	if err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}

	return filePath, nil
}

func (m *MediaRepository) Remove(ctx context.Context, media *mediaaggre.Media) error {
	filepath := filepath.Join(m.rootDir, media.GetMetadata().GetUploaderID().String(), media.GetFilename())

	if err := os.Remove(filepath); err != nil {
		if os.IsNotExist(err) {
			return FileNotExistErr
		}
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}
