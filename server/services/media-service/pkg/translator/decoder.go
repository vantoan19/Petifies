package translator

import (
	"context"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	"github.com/vantoan19/Petifies/server/services/media-service/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	MustBeProtoReqErr     = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr    = status.Error(codes.InvalidArgument, "must be proto response")
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
)

func DecodeRemoveFileByURIRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*commonProto.RemoveFileByURIRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.RemoveFileByURIReq{
		URI: req.Uri,
	}, nil
}
