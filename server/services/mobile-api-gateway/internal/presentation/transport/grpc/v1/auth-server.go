package v1

import (
	"context"
	"io"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-gateway/v1"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	mediaclient "github.com/vantoan19/Petifies/server/services/grpc-clients/media-client"
	mediaModels "github.com/vantoan19/Petifies/server/services/media-service/pkg/models"
	endpoints "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/presentation/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/translator"
	postTranslator "github.com/vantoan19/Petifies/server/services/post-service/pkg/translators"
	userTranslator "github.com/vantoan19/Petifies/server/services/user-service/pkg/translator"
)

var logger = logging.New("MobileAPIGateway.AuthServer")

type gRPCAuthServer struct {
	mediaClient       mediaclient.MediaClient
	getMyInfo         grpctransport.Handler
	userCreatePost    grpctransport.Handler
	userCreateComment grpctransport.Handler
	userEditPost      grpctransport.Handler
	userEditComment   grpctransport.Handler
	removeFileByURI   grpctransport.Handler
}

func NewAuthServer(mediaConn *grpc.ClientConn, userEndpoints endpoints.UserEndpoints, postEndpoints endpoints.PostEndpoints) authProtoV1.AuthGatewayServer {
	mediaClient := mediaclient.New(mediaConn)
	return &gRPCAuthServer{
		mediaClient: mediaClient,
		getMyInfo: grpctransport.NewServer(
			userEndpoints.GetMyInfo,
			decodeGetMyInfoRequest,
			userTranslator.EncodeGetUserResponse,
		),
		userCreatePost: grpctransport.NewServer(
			postEndpoints.CreatePost,
			translator.DecodeUserCreatePostRequest,
			postTranslator.EncodePostResponse,
		),
		userCreateComment: grpctransport.NewServer(
			postEndpoints.CreateComment,
			translator.DecodeUserCreateCommentRequest,
			postTranslator.EncodeCommentResponse,
		),
		userEditPost: grpctransport.NewServer(
			postEndpoints.EditPost,
			translator.DecodeUserEditPostRequest,
			postTranslator.EncodePostResponse,
		),
		userEditComment: grpctransport.NewServer(
			postEndpoints.EditComment,
			translator.DecodeUserEditCommentRequest,
			postTranslator.EncodeCommentResponse,
		),
		removeFileByURI: grpctransport.NewServer(
			makeRemoveFileByURIEndpoint(mediaClient),
			commonutils.CreateClientForwardDecodeRequestFunc[*commonProto.RemoveFileByURIRequest](),
			commonutils.CreateClientForwardEncodeResponseFunc[*commonProto.RemoveFileByURIResponse](),
		),
	}
}

// =============================

func decodeGetMyInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	return nil, nil
}

func (s *gRPCAuthServer) GetMyInfo(ctx context.Context, req *authProtoV1.GetMyInfoRequest) (*commonProto.User, error) {
	_, resp, err := s.getMyInfo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.User), nil
}

func (s *gRPCAuthServer) UserCreatePost(ctx context.Context, req *authProtoV1.UserCreatePostRequest) (*commonProto.Post, error) {
	_, resp, err := s.userCreatePost.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Post), nil
}

func (s *gRPCAuthServer) UserCreateComment(ctx context.Context, req *authProtoV1.UserCreateCommentRequest) (*commonProto.Comment, error) {
	_, resp, err := s.userCreateComment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Comment), nil
}

func (s *gRPCAuthServer) UserEditPost(ctx context.Context, req *authProtoV1.UserEditPostRequest) (*commonProto.Post, error) {
	_, resp, err := s.userEditPost.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Post), nil
}

func (s *gRPCAuthServer) UserEditComment(ctx context.Context, req *authProtoV1.UserEditCommentRequest) (*commonProto.Comment, error) {
	_, resp, err := s.userEditComment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Comment), nil
}

func makeRemoveFileByURIEndpoint(s mediaclient.MediaClient) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commonProto.RemoveFileByURIRequest)
		resp, err := s.RemoveFileByURIForward(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func (s *gRPCAuthServer) RemoveFileByURI(ctx context.Context, req *commonProto.RemoveFileByURIRequest) (*commonProto.RemoveFileByURIResponse, error) {
	_, resp, err := s.removeFileByURI.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.RemoveFileByURIResponse), nil
}

// ============ Stream Endpoints ================

func (s *gRPCAuthServer) UserUploadFile(stream authProtoV1.AuthGateway_UserUploadFileServer) error {
	logger.Info("Start UserUploadFile")

	logger.Info("Executing UserUploadFile: Reading metadata")
	req, err := stream.Recv()
	if err != nil {
		logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}

	uploaderId, err := uuid.Parse(req.GetMetadata().UploaderId)
	if err != nil {
		logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	clientStream, err := s.mediaClient.CreateUploadFileStream(context.Background())
	if err != nil {
		logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}

	md := mediaModels.FileMetadata{
		FileName:   req.GetMetadata().FileName,
		MediaType:  req.GetMetadata().MediaType,
		UploaderId: uploaderId,
		Size:       int64(req.GetMetadata().Size),
		Width:      int(req.GetMetadata().Width),
		Height:     int(req.GetMetadata().Height),
		Duration:   req.GetMetadata().Duration.AsDuration(),
	}
	logger.Info("Executing UserUploadFile: Uploading metadata to MediaService")
	err = s.mediaClient.UploadFileMetadata(clientStream, &md)
	if err != nil {
		return err
	}

	logger.Info("Executing UserUploadFile: Reading data")
	willBeDiscarded := true

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			logger.Info("Executing UserUploadFile: EOF, done reading data")
			break
		}
		if err != nil {
			logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
			return status.Errorf(codes.Internal, err.Error())
		}

		chunk := req.GetChunkData()
		if chunk != nil {
			logger.Info("Executing UserUploadFile: Uploading chunk data to MediaService")
			err = s.mediaClient.UploadFileChunkData(clientStream, chunk, len(chunk))
			if err != nil {
				logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
				return status.Errorf(codes.Internal, err.Error())
			}
		} else {
			willBeDiscarded = req.GetWillBeDiscarded()
			break
		}
	}

	var resp *commonProto.UploadFileResponse
	if !willBeDiscarded {
		logger.Info("Executing UserUploadFile: sending approving signal to MediaService")
		err = s.mediaClient.ApproveFile(clientStream)
		if err != nil {
			logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
			return err
		}
		resp, err = clientStream.CloseAndRecv()
	} else {
		logger.Info("Executing UserUploadFile: sending discard signal to MediaService")
		err = s.mediaClient.DiscardFile(clientStream)
		if err != nil {
			logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
			return err
		}
		resp, err = clientStream.CloseAndRecv()
	}
	err = stream.SendAndClose(resp)
	if err != nil {
		logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish UserUploadFile: Successful")
	return nil
}
