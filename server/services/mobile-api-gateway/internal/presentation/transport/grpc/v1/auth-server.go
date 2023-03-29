package v1

import (
	"context"
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-gateway/v1"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	mediaclient "github.com/vantoan19/Petifies/server/services/grpc-clients/media-client"
	newfeedclient "github.com/vantoan19/Petifies/server/services/grpc-clients/newfeed-client"
	mediaModels "github.com/vantoan19/Petifies/server/services/media-service/pkg/models"
	feedservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/feed"
	postservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/post"
	relationshipservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/relationship"
	endpoints "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/presentation/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/translator"
	userTranslator "github.com/vantoan19/Petifies/server/services/user-service/pkg/translator"
)

var logger = logging.New("MobileAPIGateway.AuthServer")

type gRPCAuthServer struct {
	mediaClient         mediaclient.MediaClient
	newfeedClient       newfeedclient.NewfeedClient
	postService         postservice.PostService
	relationshipService relationshipservice.RelationshipService
	newfeedService      feedservice.FeedService

	getMyInfo           grpctransport.Handler
	userCreatePost      grpctransport.Handler
	userCreateComment   grpctransport.Handler
	userEditPost        grpctransport.Handler
	userEditComment     grpctransport.Handler
	removeFileByURI     grpctransport.Handler
	userToggleLoveReact grpctransport.Handler
}

func NewAuthServer(
	mediaConn *grpc.ClientConn,
	postService postservice.PostService,
	relationshipService relationshipservice.RelationshipService,
	feedService feedservice.FeedService,
	userEndpoints endpoints.UserEndpoints,
	postEndpoints endpoints.PostEndpoints,
) authProtoV1.AuthGatewayServer {
	mediaClient := mediaclient.New(mediaConn)
	return &gRPCAuthServer{
		mediaClient:         mediaClient,
		postService:         postService,
		relationshipService: relationshipService,
		newfeedService:      feedService,
		getMyInfo: grpctransport.NewServer(
			userEndpoints.GetMyInfo,
			decodeGetMyInfoRequest,
			userTranslator.EncodeGetUserResponse,
		),
		userCreatePost: grpctransport.NewServer(
			postEndpoints.CreatePost,
			translator.DecodeUserCreatePostRequest,
			translator.EncodePostWithUserInfo,
		),
		userCreateComment: grpctransport.NewServer(
			postEndpoints.CreateComment,
			translator.DecodeUserCreateCommentRequest,
			translator.EncodeCommentWithUserInfo,
		),
		userEditPost: grpctransport.NewServer(
			postEndpoints.EditPost,
			translator.DecodeUserEditPostRequest,
			translator.EncodePostWithUserInfo,
		),
		userEditComment: grpctransport.NewServer(
			postEndpoints.EditComment,
			translator.DecodeUserEditCommentRequest,
			translator.EncodeCommentWithUserInfo,
		),
		removeFileByURI: grpctransport.NewServer(
			makeRemoveFileByURIEndpoint(mediaClient),
			commonutils.CreateClientForwardDecodeRequestFunc[*commonProto.RemoveFileByURIRequest](),
			commonutils.CreateClientForwardEncodeResponseFunc[*commonProto.RemoveFileByURIResponse](),
		),
		userToggleLoveReact: grpctransport.NewServer(
			postEndpoints.UserToggleLoveReact,
			translator.DecodeUserToggleLoveRequest,
			translator.EncodeUserToggleLoveResponse,
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

func (s *gRPCAuthServer) UserCreatePost(ctx context.Context, req *authProtoV1.UserCreatePostRequest) (*authProtoV1.PostWithUserInfo, error) {
	_, resp, err := s.userCreatePost.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*authProtoV1.PostWithUserInfo), nil
}

func (s *gRPCAuthServer) UserCreateComment(ctx context.Context, req *authProtoV1.UserCreateCommentRequest) (*authProtoV1.CommentWithUserInfo, error) {
	_, resp, err := s.userCreateComment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*authProtoV1.CommentWithUserInfo), nil
}

func (s *gRPCAuthServer) UserEditPost(ctx context.Context, req *authProtoV1.UserEditPostRequest) (*authProtoV1.PostWithUserInfo, error) {
	_, resp, err := s.userEditPost.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*authProtoV1.PostWithUserInfo), nil
}

func (s *gRPCAuthServer) UserEditComment(ctx context.Context, req *authProtoV1.UserEditCommentRequest) (*authProtoV1.CommentWithUserInfo, error) {
	_, resp, err := s.userEditComment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*authProtoV1.CommentWithUserInfo), nil
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

func (s *gRPCAuthServer) UserToggleLoveReact(ctx context.Context, req *authProtoV1.UserToggleLoveRequest) (*authProtoV1.UserToggleLoveResponse, error) {
	_, resp, err := s.userToggleLoveReact.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*authProtoV1.UserToggleLoveResponse), nil
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
		if err := contextError(stream.Context()); err != nil {
			return err
		}

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
		if err != nil {
			logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
			return err
		}
	} else {
		logger.Info("Executing UserUploadFile: sending discard signal to MediaService")
		err = s.mediaClient.DiscardFile(clientStream)
		if err != nil {
			logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
			return err
		}
		resp, err = clientStream.CloseAndRecv()
		if err != nil {
			logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
			return err
		}
	}
	err = stream.SendAndClose(resp)
	if err != nil {
		logger.ErrorData("Finished UserUploadFile: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish UserUploadFile: Successful")
	return nil
}

func (s *gRPCAuthServer) ListNewFeeds(stream authProtoV1.AuthGateway_ListNewFeedsServer) error {
	logger.Info("Start ListNewFeeds")

	userID, err := commonutils.GetUserID(stream.Context())
	if err != nil {
		return err
	}

	for {
		if err := contextError(stream.Context()); err != nil {
			logger.ErrorData("Finished ListNewFeeds: Failed", logging.Data{"error": err.Error()})
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			logger.Info("Executing ListNewFeeds: EOF, done reading data")
			break
		}
		if err != nil {
			logger.ErrorData("Finished ListNewFeeds: Failed", logging.Data{"error": err.Error()})
			return status.Errorf(codes.Internal, err.Error())
		}
		var afterPostID uuid.UUID
		afterPostID, err = uuid.Parse(req.AfterPostId)
		if err != nil {
			afterPostID = uuid.Nil
		}

		logger.Info("Executing ListNewFeeds: Received request to send next post feeds")
		posts, err := s.newfeedService.ListBatchPostFeeds(stream.Context(), userID, afterPostID)
		if err != nil {
			logger.ErrorData("Finished ListNewFeeds: Failed", logging.Data{"error": err.Error()})
			return err
		}

		err = stream.Send(&authProtoV1.ListNewFeedsResponse{
			Posts: commonutils.Map2(posts, func(p *models.PostWithUserInfo) *authProtoV1.PostWithUserInfo {
				return translator.EncodePostWithUserInfoHelper(p)
			}),
		})
	}

	return nil
}

func (s *gRPCAuthServer) UserListCommentsByParentID(stream authProtoV1.AuthGateway_UserListCommentsByParentIDServer) error {
	logger.Info("Start UserListCommentsByParentID")

	for {
		if err := contextError(stream.Context()); err != nil {
			logger.ErrorData("Finished UserListCommentsByParentID: Failed", logging.Data{"error": err.Error()})
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			logger.Info("Executing UserListCommentsByParentID: EOF, done reading data")
			break
		}
		if err != nil {
			logger.ErrorData("Finished UserListCommentsByParentID: Failed", logging.Data{"error": err.Error()})
			return status.Errorf(codes.Internal, err.Error())
		}

		parentID, err := uuid.Parse(req.ParentId)
		if err != nil {
			logger.ErrorData("Finished UserListCommentsByParentID: Failed", logging.Data{"error": err.Error()})
			return status.Errorf(codes.InvalidArgument, err.Error())
		}

		var afterCommentID uuid.UUID
		afterCommentID, err = uuid.Parse(req.AfterCommentId)
		if err != nil {
			afterCommentID = uuid.Nil
		}

		logger.Info("Executing UserListCommentsByParentID: Received request to send next post feeds")
		comments, err := s.postService.ListCommentsWithUserInfosByParentID(stream.Context(), parentID, int(req.PageSize), afterCommentID)
		if err != nil {
			logger.ErrorData("Finished UserListCommentsByParentID: Failed", logging.Data{"error": err.Error()})
			return err
		}

		err = stream.Send(&authProtoV1.UserListCommentsByParentIDResponse{
			Comments: commonutils.Map2(comments, func(c *models.CommentWithUserInfo) *authProtoV1.CommentWithUserInfo {
				return translator.EncodeCommentWithUserInfoHelper(c)
			}),
		})
	}

	return nil
}

func (s *gRPCAuthServer) StreamLoveCount(req *authProtoV1.StreamLoveCountRequest, stream authProtoV1.AuthGateway_StreamLoveCountServer) error {
	logger.Info("Start StreamLoveCount")
	targetID, err := uuid.Parse(req.TargetId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	for {
		if err := contextError(stream.Context()); err != nil {
			logger.Info("Finished StreamLoveCount: Successful - received cancel signal from client or deadline exceeded")
			return nil
		}

		var count int
		if req.IsPostTarget {
			count, err = s.postService.GetPostLoveCount(stream.Context(), targetID)
			if err != nil {
				logger.ErrorData("Finish StreamLoveCount: Failed", logging.Data{"error": err.Error()})
				return err
			}
		} else {
			count, err = s.postService.GetCommentLoveCount(stream.Context(), targetID)
			if err != nil {
				logger.ErrorData("Finish StreamLoveCount: Failed", logging.Data{"error": err.Error()})
				return err
			}
		}

		logger.Info("Executing StreamLoveCount: sending response")
		err = stream.Send(&authProtoV1.StreamLoveCountResponse{
			LoveCount: &wrapperspb.Int32Value{
				Value: int32(count),
			},
		})
		if err != nil {
			logger.ErrorData("Finished StreamLoveCount: Failed", logging.Data{"error": err.Error()})
			return err
		}
		// Send love count every 5 seconds
		time.Sleep(time.Second * 5)
	}
}

func (s *gRPCAuthServer) StreamCommentCount(req *authProtoV1.StreamCommentCountRequest, stream authProtoV1.AuthGateway_StreamCommentCountServer) error {
	logger.Info("Start StreamCommentCount")
	parentID, err := uuid.Parse(req.ParentId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	for {
		if err := contextError(stream.Context()); err != nil {
			logger.Info("Finished StreamCommentCount: Successful - received cancel signal from client or deadline exceeded")
			return nil
		}

		var count int
		if req.IsPostParent {
			count, err = s.postService.GetPostCommentCount(stream.Context(), parentID)
			if err != nil {
				logger.ErrorData("Finish StreamCommentCount: Failed", logging.Data{"error": err.Error()})
				return err
			}
		} else {
			count, err = s.postService.GetCommentSubCommentCount(stream.Context(), parentID)
			if err != nil {
				logger.ErrorData("Finish StreamCommentCount: Failed", logging.Data{"error": err.Error()})
				return err
			}
		}

		logger.Info("Executing StreamCommentCount: sending response")
		err = stream.Send(&authProtoV1.StreamCommentCountResponse{
			CommentCount: &wrapperspb.Int32Value{
				Value: int32(count),
			},
		})
		if err != nil {
			logger.ErrorData("Finish StreamCommentCount: Failed", logging.Data{"error": err.Error()})
			return err
		}
		// Send love count every 5 seconds
		time.Sleep(time.Second * 5)
	}
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Error(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	default:
		return nil
	}
}
