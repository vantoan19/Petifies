package services

import (
	"bytes"
	"context"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/media-service/cmd"
	mediaaggre "github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media"
	"github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media/repository"
	"github.com/vantoan19/Petifies/server/services/media-service/internal/infra/repositories/filesystem"
	"github.com/vantoan19/Petifies/server/services/media-service/pkg/models"
)

var logger = logging.New("MediaService.Service")

type MediaConfiguration func(ms *mediaService) error

type mediaService struct {
	mediaRepo repository.MediaRepository
}

type MediaService interface {
	UploadFile(ctx context.Context, md *models.FileMetadata, data *bytes.Buffer) (string, error)
}

func NewMediaService(cfgs ...MediaConfiguration) (MediaService, error) {
	media := &mediaService{}
	for _, cfg := range cfgs {
		if err := cfg(media); err != nil {
			return nil, err
		}
	}
	return media, nil
}

func WithMediaRepository(r repository.MediaRepository) MediaConfiguration {
	return func(ms *mediaService) error {
		ms.mediaRepo = r
		return nil
	}
}

func WithInDiskMediaRepository() MediaConfiguration {
	return func(ms *mediaService) error {
		repo, err := filesystem.New(cmd.Conf.StorageRootDir)
		if err != nil {
			return err
		}

		ms.mediaRepo = repo
		return nil
	}
}

func (m *mediaService) UploadFile(ctx context.Context, md *models.FileMetadata, data *bytes.Buffer) (string, error) {
	logger.Info("Start UploadFile")

	media, err := mediaaggre.New(md, data)
	if err != nil {
		logger.ErrorData("Finish UploadFile: FAILED", logging.Data{"error": err.Error()})
		return "", err
	}

	uri, err := m.mediaRepo.Save(ctx, media)
	if err != nil {
		logger.ErrorData("Finish UploadFile: FAILED", logging.Data{"error": err.Error()})
		return "", err
	}

	logger.Info("Finish UploadFile: SUCCESSFUL")
	return uri, nil
}
