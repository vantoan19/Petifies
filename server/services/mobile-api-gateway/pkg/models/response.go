package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type BasicUserInfo struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	UserAvatar string    `json:"user_avatar"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
}

type PostWithUserInfo struct {
	ID           uuid.UUID      `json:"id"`
	Author       BasicUserInfo  `json:"author"`
	Content      string         `json:"content"`
	Images       []models.Image `json:"images"`
	Videos       []models.Video `json:"videos"`
	LoveCount    int            `json:"love_count"`
	CommentCount int            `json:"comment_count"`
	Visibility   string         `json:"visibility"`
	Activity     string         `json:"activity"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CommentWithUserInfo struct {
	ID              uuid.UUID     `json:"id"`
	Author          BasicUserInfo `json:"author"`
	PostID          uuid.UUID     `json:"post_id"`
	ParentID        uuid.UUID     `json:"parent_id"`
	IsPostParent    bool          `json:"is_post_parent"`
	Content         string        `json:"content"`
	Image           models.Image  `json:"image"`
	Video           models.Video  `json:"video"`
	LoveCount       int           `json:"love_count"`
	SubcommentCount int           `json:"subcomment_count"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}
