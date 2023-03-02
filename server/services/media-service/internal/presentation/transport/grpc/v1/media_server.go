package v1

import (
	"bytes"
	"io"

	"github.com/google/uuid"
	mediaProtoV1 "github.com/vantoan19/Petifies/proto/media-service/v1"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/media-service/cmd"
	"github.com/vantoan19/Petifies/server/services/media-service/internal/application/services"
	"github.com/vantoan19/Petifies/server/services/media-service/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = logging.New("MediaService.MediaServer")

type mediaServer struct {
	mediaService services.MediaService
}

func NewMediaServer(mediaService services.MediaService) mediaProtoV1.MediaServiceServer {
	return &mediaServer{
		mediaService: mediaService,
	}
}

func (m *mediaServer) UploadFile(stream mediaProtoV1.MediaService_UploadFileServer) error {
	logger.Info("Start UploadFile")

	logger.Info("Executing UploadFile: Reading metadata")
	req, err := stream.Recv()
	if err != nil {
		logger.ErrorData("Finished UploadFile: Failed", logging.Data{"error": err.Error()})
		return err
	}

	uploaderId, err := uuid.Parse(req.GetMetadata().UploaderId)
	if err != nil {
		logger.ErrorData("Finished UploadFile: Failed", logging.Data{"error": err.Error()})
		return err
	}

	md := models.FileMetadata{
		FileName:   req.GetMetadata().FileName,
		MediaType:  req.GetMetadata().MediaType,
		UploaderId: uploaderId,
		Size:       int64(req.GetMetadata().Size),
		Width:      int(req.GetMetadata().Width),
		Height:     int(req.GetMetadata().Height),
		Duration:   req.GetMetadata().Duration.AsDuration(),
	}

	logger.Info("Executing UploadFile: Reading data")
	data := bytes.Buffer{}
	recvSize := 0
	willBeDiscarded := true

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			logger.Info("Executing UploadFile: EOF, done reading data")
			break
		}
		if err != nil {
			logger.ErrorData("Finished UploadFile: Failed", logging.Data{"error": err.Error()})
			return err
		}

		chunk := req.GetChunkData()
		if chunk != nil {
			size := len(chunk)

			recvSize += size
			if recvSize > cmd.Conf.MaxFileSize {
				logger.Error("Finished UploadFile: Failed - file too big")
				return status.Errorf(codes.Canceled, "file too big")
			}

			_, err = data.Write(chunk)
			if err != nil {
				logger.ErrorData("Finished UploadFile: Failed", logging.Data{"error": err.Error()})
				return err
			}
		} else {
			willBeDiscarded = req.GetWillBeDiscarded()
		}
	}

	var resp *mediaProtoV1.UploadFileResponse
	if !willBeDiscarded {
		uri, err := m.mediaService.UploadFile(stream.Context(), &md, &data)
		if err != nil {
			logger.ErrorData("Finished UploadFile: Failed", logging.Data{"error": err.Error()})
			return err
		}
		resp = &mediaProtoV1.UploadFileResponse{
			Uri:  uri,
			Size: uint64(recvSize),
		}
	} else {
		resp = &mediaProtoV1.UploadFileResponse{
			Uri:  "",
			Size: 0,
		}
	}
	err = stream.SendAndClose(resp)
	if err != nil {
		logger.ErrorData("Finished UploadFile: Failed", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finished UploadFile: Success")
	return nil
}
