package models

import (
	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type UserCreatePostReq struct {
	TextContent string
	Visibility  string
	Activity    string
	Images      []models.Image
	Videos      []models.Video
}

type UserCreateCommentReq struct {
	PostID       uuid.UUID
	ParentID     uuid.UUID
	IsParentPost bool
	Content      string
	Image        models.Image
	Video        models.Video
}

type UserEditPostReq struct {
	PostID     uuid.UUID
	Content    string
	Visibility string
	Activity   string
	Images     []models.Image
	Videos     []models.Video
}

type UserEditCommentReq struct {
	CommentID uuid.UUID
	Content   string
	Image     models.Image
	Video     models.Video
}

type UserToggleLoveReq struct {
	TargetID     uuid.UUID
	IsPostTarget bool
}
