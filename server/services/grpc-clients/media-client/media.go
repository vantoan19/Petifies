package mediaclient

import (
	"bytes"
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/vantoan19/Petifies/proto/common"
	mediaProtoV1 "github.com/vantoan19/Petifies/proto/media-service/v1"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/media-service/pkg/models"
)

var logger = logging.New("Clients.MediaClient")

const mediaService = "media_service.v1.MediaService"

type mediaClient struct {
	client                 mediaProtoV1.MediaServiceClient
	removeFileByURIForward endpoint.Endpoint
}

type MediaClient interface {
	CreateUploadFileStream(ctx context.Context) (mediaProtoV1.MediaService_UploadFileClient, error)
	UploadFileMetadata(stream mediaProtoV1.MediaService_UploadFileClient, md *models.FileMetadata) error
	UploadFileChunkData(stream mediaProtoV1.MediaService_UploadFileClient, chunk []byte, len int) error
	UploadFile(ctx context.Context, md *models.FileMetadata, data *bytes.Buffer) (string, int, error)
	DiscardFile(stream mediaProtoV1.MediaService_UploadFileClient) error
	ApproveFile(stream mediaProtoV1.MediaService_UploadFileClient) error
	RemoveFileByURIForward(ctx context.Context, req *common.RemoveFileByURIRequest) (*common.RemoveFileByURIResponse, error)
}

func New(conn *grpc.ClientConn) MediaClient {
	return &mediaClient{
		client: mediaProtoV1.NewMediaServiceClient(conn),
		removeFileByURIForward: grpctransport.NewClient(
			conn,
			mediaService,
			"RemoveFileByURI",
			commonutils.CreateClientForwardEncodeRequestFunc[*common.RemoveFileByURIRequest](),
			commonutils.CreateClientForwardDecodeResponseFunc[*common.RemoveFileByURIResponse](),
			common.RemoveFileByURIResponse{},
		).Endpoint(),
	}
}

func (m *mediaClient) CreateUploadFileStream(ctx context.Context) (mediaProtoV1.MediaService_UploadFileClient, error) {
	logger.Info("Start CreateUploadFileStream")
	stream, err := m.client.UploadFile(ctx)

	logger.Info("Finished CreateUploadFileStream: Successful")
	return stream, err
}

func (m *mediaClient) UploadFileMetadata(stream mediaProtoV1.MediaService_UploadFileClient, md *models.FileMetadata) error {
	logger.Info("Start UploadFileMetadata")
	mdReq := &common.UploadFileRequest{
		Data: &common.UploadFileRequest_Metadata{
			Metadata: &common.FileMetadata{
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

	logger.Info("Finished UploadFileMetadata: Successful")
	return stream.Send(mdReq)
}

func (m *mediaClient) UploadFileChunkData(stream mediaProtoV1.MediaService_UploadFileClient, chunk []byte, len int) error {
	logger.Info("Start UploadFileChunkData")
	dataReq := &common.UploadFileRequest{
		Data: &common.UploadFileRequest_ChunkData{
			ChunkData: chunk[:len],
		},
	}

	logger.Info("Finished UploadFileChunkData: Successful")
	return stream.Send(dataReq)
}

func (m *mediaClient) DiscardFile(stream mediaProtoV1.MediaService_UploadFileClient) error {
	dataReq := &common.UploadFileRequest{
		Data: &common.UploadFileRequest_WillBeDiscarded{
			WillBeDiscarded: true,
		},
	}
	return stream.Send(dataReq)
}

func (m *mediaClient) ApproveFile(stream mediaProtoV1.MediaService_UploadFileClient) error {
	dataReq := &common.UploadFileRequest{
		Data: &common.UploadFileRequest_WillBeDiscarded{
			WillBeDiscarded: false,
		},
	}
	return stream.Send(dataReq)
}

func (m *mediaClient) UploadFile(ctx context.Context, md *models.FileMetadata, data *bytes.Buffer) (string, int, error) {
	logger.Info("Start UploadFile")
	stream, err := m.CreateUploadFileStream(ctx)
	if err != nil {
		logger.ErrorData("Finish UploadFile: Failed", logging.Data{"error": err.Error()})
		return "", 0, err
	}

	err = m.UploadFileMetadata(stream, md)
	if err != nil {
		logger.ErrorData("Finish UploadFile: Failed", logging.Data{"error": err.Error()})
		return "", 0, err
	}

	err = m.UploadFileChunkData(stream, data.Bytes(), data.Len())
	if err != nil {
		logger.ErrorData("Finish UploadFile: Failed", logging.Data{"error": err.Error()})
		return "", 0, err
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		logger.ErrorData("Finish UploadFile: Failed", logging.Data{"error": err.Error()})
		return "", 0, err
	}

	logger.Info("Finished UploadFile: Successful")
	return resp.Uri, int(resp.Size), nil
}

func (m *mediaClient) RemoveFileByURIForward(ctx context.Context, req *common.RemoveFileByURIRequest) (*common.RemoveFileByURIResponse, error) {
	logger.Info("Start RemoveFileByURIForward")

	resp, err := m.removeFileByURIForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished RemoveFileByURIForward: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished RemoveFileByURIForward: SUCCESSFUL")
	return resp.(*common.RemoveFileByURIResponse), nil
}
