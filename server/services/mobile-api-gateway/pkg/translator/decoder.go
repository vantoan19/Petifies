package translator

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-gateway/v1"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
	petifiesModels "github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
	postModels "github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var (
	MustBeProtoReqErr     = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr    = status.Error(codes.InvalidArgument, "must be proto response")
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
)

func DecodeUserCreatePostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserCreatePostRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.UserCreatePostReq{
		TextContent: req.Content,
		Visibility:  req.Visibility,
		Activity:    req.Activity,
		Images: utils.Map2(req.Images, func(i *commonProto.Image) postModels.Image {
			return postModels.Image{URL: i.Uri, Description: i.Description}
		}),
		Videos: utils.Map2(req.Videos, func(v *commonProto.Video) postModels.Video {
			return postModels.Video{URL: v.Uri, Description: v.Description}
		}),
	}, nil
}

func DecodeUserCreateCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserCreateCommentRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}
	parentID, err := uuid.Parse(req.ParentId)
	if err != nil {
		return nil, err
	}
	var image postModels.Image
	if req.Image != nil {
		image = postModels.Image{URL: req.Image.Uri, Description: req.Image.Description}
	}
	var video postModels.Video
	if req.Video != nil {
		video = postModels.Video{URL: req.Video.Uri, Description: req.Video.Description}
	}

	return &models.UserCreateCommentReq{
		PostID:       postID,
		ParentID:     parentID,
		IsParentPost: req.IsPostParent,
		Content:      req.Content,
		Image:        image,
		Video:        video,
	}, nil
}

func DecodeUserEditPostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserEditPostRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	postID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	return &models.UserEditPostReq{
		PostID:     postID,
		Content:    req.Content,
		Visibility: req.Visibility,
		Activity:   req.Activity,
		Images: utils.Map2(req.Images, func(i *commonProto.Image) postModels.Image {
			return postModels.Image{URL: i.Uri, Description: i.Description}
		}),
		Videos: utils.Map2(req.Videos, func(v *commonProto.Video) postModels.Video {
			return postModels.Video{URL: v.Uri, Description: v.Description}
		}),
	}, nil
}

func DecodeUserEditCommentRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserEditCommentRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	commentID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	var image postModels.Image
	if req.Image != nil {
		image = postModels.Image{URL: req.Image.Uri, Description: req.Image.Description}
	}
	var video postModels.Video
	if req.Video != nil {
		video = postModels.Video{URL: req.Video.Uri, Description: req.Video.Description}
	}

	return &models.UserEditCommentReq{
		CommentID: commentID,
		Content:   req.Content,
		Image:     image,
		Video:     video,
	}, nil
}

func DecodeUserToggleLoveRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserToggleLoveRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	targetID, err := uuid.Parse(req.TargetId)
	if err != nil {
		return nil, err
	}

	return &models.UserToggleLoveReq{
		TargetID:     targetID,
		IsPostTarget: req.IsPostTarget,
	}, nil
}

func DecodeUserCreatePetifiesRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserCreatePetifiesRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.UserCreatePetifiesReq{
		Type:        req.Type.String(),
		Title:       req.Title,
		Description: req.Description,
		PetName:     req.PetName,
		Images: utils.Map2(req.Images, func(i *commonProto.Image) petifiesModels.Image {
			return petifiesModels.Image{URI: i.Uri, Description: i.Description}
		}),
		Address: petifiesModels.Address{
			AddressLineOne: req.Address.AddressLineOne,
			AddressLineTwo: req.Address.AddressLineTwo,
			Street:         req.Address.Street,
			District:       req.Address.District,
			City:           req.Address.City,
			Region:         req.Address.Region,
			PostalCode:     req.Address.PostalCode,
			Country:        req.Address.Country,
			Longitude:      req.Address.Longitude,
			Latitude:       req.Address.Latitude,
		},
	}, nil
}

func DecodeUserCreatePetifiesSessionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserCreatePetifiesSessionRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	petifiesId, err := uuid.Parse(req.PetifiesId)
	if err != nil {
		return nil, err
	}

	return &models.UserCreatePetifiesSessionReq{
		PetifiesId: petifiesId,
		FromTime:   req.FromTime.AsTime(),
		ToTime:     req.ToTime.AsTime(),
	}, nil
}

func DecodeUserCreatePetifiesProposalRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserCreatePetifiesProposalRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	petifiesSessionId, err := uuid.Parse(req.PetifiesSessionId)
	if err != nil {
		return nil, err
	}

	return &models.UserCreatePetifiesProposalReq{
		PetifiesSessionId: petifiesSessionId,
		Proposal:          req.Proposal,
	}, nil
}

func DecodeUserCreateReviewRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.UserCreateReviewRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	petifiesId, err := uuid.Parse(req.PetifiesId)
	if err != nil {
		return nil, err
	}
	var image petifiesModels.Image
	if req.Image != nil {
		image = petifiesModels.Image{URI: req.Image.Uri, Description: req.Image.Description}
	}

	return &models.UserCreateReviewReq{
		PetifiesId: petifiesId,
		Review:     req.Review,
		Image:      image,
	}, nil
}

func DecodeListNearByPetifiesRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.ListNearByPetifiesRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &models.ListNearByPetifiesReq{
		Type:      req.Type.String(),
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
		Radius:    req.Radius,
		PageSize:  req.PageSize,
		Offset:    int(req.Offset),
	}, nil
}

func DecodeListPetifiesByUserIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.ListPetifiesByUserIdRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	afterId, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterId = uuid.Nil
	}

	return &models.ListPetifiesByUserIdReq{
		UserId:   userId,
		PageSize: int(req.PageSize),
		AfterId:  afterId,
	}, nil
}

func DecodeListSessionsByPetifiesIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.ListSessionsByPetifiesIdRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	petifiesId, err := uuid.Parse(req.PetifiesId)
	if err != nil {
		return nil, err
	}
	afterId, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterId = uuid.Nil
	}

	return &models.ListSessionsByPetifiesIdReq{
		PetifiesId: petifiesId,
		PageSize:   int(req.PageSize),
		AfterId:    afterId,
	}, nil
}

func DecodeListProposalsBySessionIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.ListProposalsBySessionIdRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	sessionId, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, err
	}
	afterId, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterId = uuid.Nil
	}

	return &models.ListProposalsBySessionIdReq{
		SessionId: sessionId,
		PageSize:  int(req.PageSize),
		AfterId:   afterId,
	}, nil
}

func DecodeListProposalsByUserIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.ListProposalsByUserIdRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	afterId, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterId = uuid.Nil
	}

	return &models.ListProposalsByUserIdReq{
		UserId:   userId,
		PageSize: int(req.PageSize),
		AfterId:  afterId,
	}, nil
}

func DecodeListReviewsByPetifiesIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.ListReviewsByPetifiesIdRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	petifiesId, err := uuid.Parse(req.PetifiesId)
	if err != nil {
		return nil, err
	}
	afterId, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterId = uuid.Nil
	}

	return &models.ListReviewsByPetifiesIdReq{
		PetifiesId: petifiesId,
		PageSize:   int(req.PageSize),
		AfterId:    afterId,
	}, nil
}

func DecodeListReviewsByUserIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*authProtoV1.ListReviewsByUserIdRequest)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	afterId, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterId = uuid.Nil
	}

	return &models.ListReviewsByUserIdReq{
		UserId:   userId,
		PageSize: int(req.PageSize),
		AfterId:  afterId,
	}, nil
}
