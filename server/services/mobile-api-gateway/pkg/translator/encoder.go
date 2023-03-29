package translator

import (
	"context"

	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-gateway/v1"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
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
