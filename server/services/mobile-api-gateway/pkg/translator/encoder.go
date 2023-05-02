package translator

import (
	"context"

	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-gateway/v1"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
	petifiesModels "github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
	petifiesTranslator "github.com/vantoan19/Petifies/server/services/petifies-service/pkg/translators"
	postModels "github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func EncodePostWithUserInfo(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.PostWithUserInfo)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return EncodePostWithUserInfoHelper(resp), nil
}

func EncodeCommentWithUserInfo(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.CommentWithUserInfo)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return EncodeCommentWithUserInfoHelper(resp), nil
}

func EncodeLoveWithUserInfo(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.LoveWithUserInfo)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return EncodeLoveWithUserInfoHelper(resp), nil
}

func EncodeUserToggleLoveResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.UserToggleLoveResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &authProtoV1.UserToggleLoveResponse{
		HasReacted: &wrapperspb.BoolValue{
			Value: resp.HasReacted,
		},
	}, nil
}

func EncodePetifiesWithUserInfoResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.PetifiesWithUserInfo)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return EncodePetifiesWithUserInfoHelper(resp), nil
}

func EncodePetifiesProposalWithUserInfoResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.PetifiesProposalWithUserInfo)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return EncodePetifiesProposalWithUserInfoHelper(resp), nil
}

func EncodePetifiesSessionResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.PetifiesSession)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return EncodePetifiesSessionHelper(resp), nil
}

func EncodeReviewWithUserInfoResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ReviewWithUserInfo)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return EncodeReviewWithUserInfoHelper(resp), nil
}

func EncodeListNearByPetifiesResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListNearByPetifiesResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &authProtoV1.ListNearByPetifiesResponse{
		Petifies: commonutils.Map2(resp.Petifies, func(p *models.PetifiesWithUserInfo) *authProtoV1.PetifiesWithUserInfo {
			return EncodePetifiesWithUserInfoHelper(p)
		}),
	}, nil
}

func EncodeListPetifiesByUserIdResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListPetifiesByUserIdResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &authProtoV1.ListPetifiesByUserIdResponse{
		Petifies: commonutils.Map2(resp.Petifies, func(p *models.PetifiesWithUserInfo) *authProtoV1.PetifiesWithUserInfo {
			return EncodePetifiesWithUserInfoHelper(p)
		}),
	}, nil
}

func EncodeListSessionsByPetifiesIdResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListSessionsByPetifiesIdResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &authProtoV1.ListSessionsByPetifiesIdResponse{
		Sessions: commonutils.Map2(resp.Sessions, func(p *models.PetifiesSession) *authProtoV1.UserPetifiesSession {
			return EncodePetifiesSessionHelper(p)
		}),
	}, nil
}

func EncodeListProposalsBySessionIdResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListProposalsBySessionIdResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &authProtoV1.ListProposalsBySessionIdResponse{
		Proposals: commonutils.Map2(resp.Proposals, func(p *models.PetifiesProposalWithUserInfo) *authProtoV1.PetifiesProposalWithUserInfo {
			return EncodePetifiesProposalWithUserInfoHelper(p)
		}),
	}, nil
}

func EncodeListProposalsByUserIdResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListProposalsByUserIdResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &authProtoV1.ListProposalsByUserIdResponse{
		Proposals: commonutils.Map2(resp.Proposals, func(p *models.PetifiesProposalWithUserInfo) *authProtoV1.PetifiesProposalWithUserInfo {
			return EncodePetifiesProposalWithUserInfoHelper(p)
		}),
	}, nil
}

func EncodeListReviewsByPetifiesIdResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListReviewsByPetifiesIdResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &authProtoV1.ListReviewsByPetifiesIdResponse{
		Reviews: commonutils.Map2(resp.Reviews, func(p *models.ReviewWithUserInfo) *authProtoV1.ReviewWithUserInfo {
			return EncodeReviewWithUserInfoHelper(p)
		}),
	}, nil
}

func EncodeListReviewsByUserIdResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ListReviewsByUserIdResp)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &authProtoV1.ListReviewsByUserIdResponse{
		Reviews: commonutils.Map2(resp.Reviews, func(p *models.ReviewWithUserInfo) *authProtoV1.ReviewWithUserInfo {
			return EncodeReviewWithUserInfoHelper(p)
		}),
	}, nil
}

// ============= Helpers ================

func EncodeBasicUserInfoHelper(info *models.BasicUserInfo) *commonProto.BasicUser {
	return &commonProto.BasicUser{
		Id:         info.ID.String(),
		Email:      info.Email,
		UserAvatar: info.UserAvatar,
		FirstName:  info.FirstName,
		LastName:   info.LastName,
	}
}

func EncodePostWithUserInfoHelper(post *models.PostWithUserInfo) *authProtoV1.PostWithUserInfo {
	return &authProtoV1.PostWithUserInfo{
		Id:      post.ID.String(),
		Author:  EncodeBasicUserInfoHelper(&post.Author),
		Content: post.Content,
		Images: commonutils.Map2(post.Images, func(i postModels.Image) *commonProto.Image {
			return &commonProto.Image{Uri: i.URL, Description: i.Description}
		}),
		Videos: commonutils.Map2(post.Videos, func(v postModels.Video) *commonProto.Video {
			return &commonProto.Video{Uri: v.URL, Description: v.Description}
		}),
		LoveCount:    int32(post.LoveCount),
		CommentCount: int32(post.CommentCount),
		Visibility:   post.Visibility,
		Activity:     post.Activity,
		HasReacted:   post.HasReacted,
		CreatedAt:    timestamppb.New(post.CreatedAt),
		UpdatedAt:    timestamppb.New(post.UpdatedAt),
	}
}

func EncodeCommentWithUserInfoHelper(comment *models.CommentWithUserInfo) *authProtoV1.CommentWithUserInfo {
	return &authProtoV1.CommentWithUserInfo{
		Id:              comment.ID.String(),
		Author:          EncodeBasicUserInfoHelper(&comment.Author),
		PostId:          comment.PostID.String(),
		ParentId:        comment.ParentID.String(),
		IsPostParent:    comment.IsPostParent,
		Content:         comment.Content,
		Image:           &commonProto.Image{Uri: comment.Image.URL, Description: comment.Image.Description},
		Video:           &commonProto.Video{Uri: comment.Video.URL, Description: comment.Video.Description},
		LoveCount:       int32(comment.LoveCount),
		SubcommentCount: int32(comment.SubcommentCount),
		HasReacted:      comment.HasReacted,
		CreatedAt:       timestamppb.New(comment.CreatedAt),
		UpdatedAt:       timestamppb.New(comment.UpdatedAt),
	}
}

func EncodeLoveWithUserInfoHelper(love *models.LoveWithUserInfo) *authProtoV1.LoveWithUserInfo {
	return &authProtoV1.LoveWithUserInfo{
		Id:           love.ID.String(),
		TargetId:     love.TargetID.String(),
		IsPostTarget: love.IsPostTarget,
		Author:       EncodeBasicUserInfoHelper(&love.Author),
		CreatedAt:    timestamppb.New(love.CreatedAt),
	}
}

func EncodePetifiesWithUserInfoHelper(petifies *models.PetifiesWithUserInfo) *authProtoV1.PetifiesWithUserInfo {
	return &authProtoV1.PetifiesWithUserInfo{
		Id:          petifies.Id.String(),
		Owner:       EncodeBasicUserInfoHelper(&petifies.Owner),
		Type:        petifiesTranslator.GetPetifiesType(petifies.Type),
		Title:       petifies.Title,
		Description: petifies.Description,
		PetName:     petifies.PetName,
		Images: commonutils.Map2(petifies.Images, func(i petifiesModels.Image) *commonProto.Image {
			return &commonProto.Image{Uri: i.URI, Description: i.Description}
		}),
		Status: petifiesTranslator.GetPetifiesStatus(petifies.Status),
		Address: &commonProto.Address{
			AddressLineOne: petifies.Address.AddressLineOne,
			AddressLineTwo: petifies.Address.AddressLineTwo,
			Street:         petifies.Address.Street,
			District:       petifies.Address.District,
			City:           petifies.Address.City,
			Region:         petifies.Address.Region,
			PostalCode:     petifies.Address.PostalCode,
			Country:        petifies.Address.Country,
			Longitude:      petifies.Address.Longitude,
			Latitude:       petifies.Address.Latitude,
		},
		CreatedAt: timestamppb.New(petifies.CreatedAt),
		UpdatedAt: timestamppb.New(petifies.UpdatedAt),
	}
}

func EncodePetifiesSessionHelper(session *models.PetifiesSession) *authProtoV1.UserPetifiesSession {
	return &authProtoV1.UserPetifiesSession{
		Id:         session.Id.String(),
		PetifiesId: session.PetifiesId.String(),
		FromTime:   timestamppb.New(session.FromTime),
		ToTime:     timestamppb.New(session.ToTime),
		Status:     petifiesTranslator.GetPetifiesSessionStatus(session.Status),
		CreatedAt:  timestamppb.New(session.CreatedAt),
		UpdatedAt:  timestamppb.New(session.UpdatedAt),
	}
}

func EncodePetifiesProposalWithUserInfoHelper(proposal *models.PetifiesProposalWithUserInfo) *authProtoV1.PetifiesProposalWithUserInfo {
	return &authProtoV1.PetifiesProposalWithUserInfo{
		Id:                proposal.Id.String(),
		User:              EncodeBasicUserInfoHelper(&proposal.User),
		PetifiesSessionId: proposal.PetifiesSessionId.String(),
		Proposal:          proposal.Proposal,
		Status:            petifiesTranslator.GetPetifiesProposalStatus(proposal.Status),
		CreatedAt:         timestamppb.New(proposal.CreatedAt),
		UpdatedAt:         timestamppb.New(proposal.UpdatedAt),
	}
}

func EncodeReviewWithUserInfoHelper(review *models.ReviewWithUserInfo) *authProtoV1.ReviewWithUserInfo {
	return &authProtoV1.ReviewWithUserInfo{
		Id:         review.Id.String(),
		Author:     EncodeBasicUserInfoHelper(&review.Author),
		PetifiesId: review.PetifiesId.String(),
		Review:     review.Review,
		Image:      &commonProto.Image{Uri: review.Image.URI, Description: review.Image.Description},
		CreatedAt:  timestamppb.New(review.CreatedAt),
		UpdatedAt:  timestamppb.New(review.UpdatedAt),
	}
}
