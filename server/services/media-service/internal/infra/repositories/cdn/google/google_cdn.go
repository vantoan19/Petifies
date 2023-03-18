package google

import (
	"context"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/services/media-service/cmd"
	mediaaggre "github.com/vantoan19/Petifies/server/services/media-service/internal/domain/aggregates/media"
)

type MediaRepository struct {
	storageClient *storage.Client
}

func New() (*MediaRepository, error) {
	storageClient, err := storage.NewClient(context.Background(), option.WithCredentialsFile(cmd.Conf.CredentialFile))
	if err != nil {
		return nil, err
	}

	return &MediaRepository{
		storageClient: storageClient,
	}, nil
}

func (m *MediaRepository) Save(ctx context.Context, media *mediaaggre.Media) (string, error) {
	ctx_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	sw := m.storageClient.Bucket(cmd.Conf.BucketName).Object(media.GetID().String() + "-" + media.GetFilename()).NewWriter(ctx_)

	if _, err := media.GetData().WriteTo(sw); err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}

	if err := sw.Close(); err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}

	url, err := url.Parse("/" + cmd.Conf.BucketName + "/" + sw.Attrs().Name)
	if err != nil {
		return "", status.Errorf(codes.Internal, err.Error())
	}

	return cmd.Conf.CDNDomain + url.EscapedPath(), nil
}

func (m *MediaRepository) RemoveByUri(ctx context.Context, uri string) error {
	ctx_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objName := strings.TrimPrefix(uri, cmd.Conf.CDNDomain+"/"+cmd.Conf.BucketName+"/")
	obj := m.storageClient.Bucket(cmd.Conf.BucketName).Object(objName)

	if err := obj.Delete(ctx_); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}
