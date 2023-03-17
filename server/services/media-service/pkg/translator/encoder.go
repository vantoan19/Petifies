package translator

import (
	"context"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	"github.com/vantoan19/Petifies/server/services/media-service/pkg/models"
)

func EncodeRemoveFileByURIResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(*models.RemoveFileByURIResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &commonProto.RemoveFileByURIResponse{}, nil
}
