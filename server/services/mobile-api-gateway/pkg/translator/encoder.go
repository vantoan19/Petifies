package translator

import (
	"context"

	authProtoV1 "github.com/vantoan19/Petifies/proto/auth-gateway/v1"
	commonProto "github.com/vantoan19/Petifies/proto/common"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
	postModels "github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EncodePostWithUserInfo(_ context.Context, response interface{}) (interface{}, error) {
	req, ok := response.(*models.PostWithUserInfo)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return _encodePostWithUserInfo(req), nil
}

func EncodeCommentWithUserInfo(_ context.Context, response interface{}) (interface{}, error) {
	req, ok := response.(*models.CommentWithUserInfo)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return _encodeCommentWithUserInfo(req), nil
}

func _encodeBasicUserInfo(info *models.BasicUserInfo) *commonProto.BasicUser {
	return &commonProto.BasicUser{
		Id:         info.ID.String(),
		Email:      info.Email,
		UserAvatar: info.UserAvatar,
		FirstName:  info.FirstName,
		LastName:   info.LastName,
	}
}

func _encodePostWithUserInfo(post *models.PostWithUserInfo) *authProtoV1.PostWithUserInfo {
	return &authProtoV1.PostWithUserInfo{
		Id:      post.ID.String(),
		Author:  _encodeBasicUserInfo(&post.Author),
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
		CreatedAt:    timestamppb.New(post.CreatedAt),
		UpdatedAt:    timestamppb.New(post.UpdatedAt),
	}
}

func _encodeCommentWithUserInfo(comment *models.CommentWithUserInfo) *authProtoV1.CommentWithUserInfo {
	return &authProtoV1.CommentWithUserInfo{
		Id:              comment.ID.String(),
		Author:          _encodeBasicUserInfo(&comment.Author),
		PostId:          comment.PostID.String(),
		ParentId:        comment.ParentID.String(),
		IsPostParent:    comment.IsPostParent,
		Content:         comment.Content,
		Image:           &commonProto.Image{Uri: comment.Image.URL, Description: comment.Image.Description},
		Video:           &commonProto.Video{Uri: comment.Video.URL, Description: comment.Video.Description},
		LoveCount:       int32(comment.LoveCount),
		SubcommentCount: int32(comment.SubcommentCount),
		CreatedAt:       timestamppb.New(comment.CreatedAt),
		UpdatedAt:       timestamppb.New(comment.UpdatedAt),
	}
}
