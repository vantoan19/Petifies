package models

import (
	"time"

	"github.com/google/uuid"
	petifiesModel "github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
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
	HasReacted   bool           `json:"has_reacted"`
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
	HasReacted      bool          `json:"has_reacted"`
	SubcommentCount int           `json:"subcomment_count"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type LoveWithUserInfo struct {
	ID           uuid.UUID     `json:"id"`
	TargetID     uuid.UUID     `json:"target_id"`
	IsPostTarget bool          `json:"is_post_target"`
	Author       BasicUserInfo `json:"author"`
	CreatedAt    time.Time     `json:"created_at"`
}

type UserToggleLoveResp struct {
	HasReacted bool
}

type PetifiesWithUserInfo struct {
	Id          uuid.UUID             `json:"id"`
	Owner       BasicUserInfo         `json:"owner"`
	Type        string                `json:"type"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	PetName     string                `json:"pet_name"`
	Images      []petifiesModel.Image `json:"images"`
	Status      string                `json:"status"`
	Address     petifiesModel.Address `json:"address"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type PetifiesSession struct {
	Id         uuid.UUID `json:"id"`
	PetifiesId uuid.UUID `json:"petifies_id"`
	FromTime   time.Time `json:"from_time"`
	ToTime     time.Time `json:"to_time"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PetifiesProposalWithUserInfo struct {
	Id                uuid.UUID     `json:"id"`
	User              BasicUserInfo `json:"user"`
	PetifiesSessionId uuid.UUID     `json:"petifies_session_id"`
	Proposal          string        `json:"proposal"`
	Status            string        `json:"status"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

type ReviewWithUserInfo struct {
	Id         uuid.UUID           `json:"id"`
	PetifiesId uuid.UUID           `json:"petifies_id"`
	Author     BasicUserInfo       `json:"author"`
	Review     string              `json:"review"`
	Image      petifiesModel.Image `json:"image"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

type ListNearByPetifiesResp struct {
	Petifies []*PetifiesWithUserInfo
}

type ListPetifiesByUserIdResp struct {
	Petifies []*PetifiesWithUserInfo
}

type ListSessionsByPetifiesIdResp struct {
	Sessions []*PetifiesSession
}

type ListProposalsBySessionIdResp struct {
	Proposals []*PetifiesProposalWithUserInfo
}

type ListProposalsByUserIdResp struct {
	Proposals []*PetifiesProposalWithUserInfo
}

type ListReviewsByPetifiesIdResp struct {
	Reviews []*ReviewWithUserInfo
}

type ListReviewsByUserIdResp struct {
	Reviews []*ReviewWithUserInfo
}
