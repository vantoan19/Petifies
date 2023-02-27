package mediaclient

import (
	"bytes"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	mediaProtoV1 "github.com/vantoan19/Petifies/proto/media-service/v1"
	"github.com/vantoan19/Petifies/server/services/media-service/pkg/models"
)

type mediaClient struct {
	client mediaProtoV1.MediaServiceClient
}

type MediaClient interface {
	CreateUploadFileStream(ctx context.Context) (mediaProtoV1.MediaService_UploadFileClient, error)
	UploadFileMetadata(stream mediaProtoV1.MediaService_UploadFileClient, md *models.FileMetadata) error
	UploadFileChunkData(stream mediaProtoV1.MediaService_UploadFileClient, chunk []byte, len int) error
	UploadFile(ctx context.Context, md *models.FileMetadata, data *bytes.Buffer) (string, int, error)
}

func New(conn *grpc.ClientConn) MediaClient {
	return &mediaClient{
		client: mediaProtoV1.NewMediaServiceClient(conn),
	}
}

func (m *mediaClient) CreateUploadFileStream(ctx context.Context) (mediaProtoV1.MediaService_UploadFileClient, error) {
	stream, err := m.client.UploadFile(ctx)
	return stream, err
}

func (m *mediaClient) UploadFileMetadata(stream mediaProtoV1.MediaService_UploadFileClient, md *models.FileMetadata) error {
	mdReq := &mediaProtoV1.UploadFileRequest{
		Data: &mediaProtoV1.UploadFileRequest_Metadata{
			Metadata: &mediaProtoV1.FileMetadata{
				FileName:   md.FileName,
				MediaType:  md.MediaType,
				UploaderId: md.UploaderId.String(),
				Size:       uint64(md.Size),
				Width:      uint32(md.Width),
				Height:     uint32(md.Height),
				Duration:   durationpb.New(md.Duration),
			},
		},
	}
	return stream.Send(mdReq)
}

func (m *mediaClient) UploadFileChunkData(stream mediaProtoV1.MediaService_UploadFileClient, chunk []byte, len int) error {
	dataReq := &mediaProtoV1.UploadFileRequest{
		Data: &mediaProtoV1.UploadFileRequest_ChunkData{
			ChunkData: chunk[:len],
		},
	}
	return stream.Send(dataReq)
}

func (m *mediaClient) UploadFile(ctx context.Context, md *models.FileMetadata, data *bytes.Buffer) (string, int, error) {
	stream, err := m.CreateUploadFileStream(ctx)
	if err != nil {
		return "", 0, err
	}

	err = m.UploadFileMetadata(stream, md)
	if err != nil {
		return "", 0, err
	}

	err = m.UploadFileChunkData(stream, data.Bytes(), data.Len())
	if err != nil {
		return "", 0, err
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", 0, err
	}

	return resp.Uri, int(resp.Size), nil
}
