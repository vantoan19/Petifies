package models

import (
	"time"

	"github.com/google/uuid"
	petifiesModel "github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
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

type UserCreatePetifiesReq struct {
	Type        string
	Title       string
	Description string
	PetName     string
	Images      []petifiesModel.Image
	Address     petifiesModel.Address
}

type UserCreatePetifiesSessionReq struct {
	PetifiesId uuid.UUID
	FromTime   time.Time
	ToTime     time.Time
}

type UserCreatePetifiesProposalReq struct {
	PetifiesSessionId uuid.UUID
	Proposal          string
}

type UserCreateReviewReq struct {
	PetifiesId uuid.UUID
	Review     string
	Image      petifiesModel.Image
}

type ListNearByPetifiesReq struct {
	Type      string
	Longitude float64
	Latitude  float64
	Radius    float64
	PageSize  int32
	Offset    int
}

type ListPetifiesByUserIdReq struct {
	UserId   uuid.UUID
	PageSize int
	AfterId  uuid.UUID
}

type ListSessionsByPetifiesIdReq struct {
	PetifiesId uuid.UUID
	PageSize   int
	AfterId    uuid.UUID
}

type ListProposalsBySessionIdReq struct {
	SessionId uuid.UUID
	PageSize  int
	AfterId   uuid.UUID
}

type ListProposalsByUserIdReq struct {
	UserId   uuid.UUID
	PageSize int
	AfterId  uuid.UUID
}

type ListReviewsByPetifiesIdReq struct {
	PetifiesId uuid.UUID
	PageSize   int
	AfterId    uuid.UUID
}

type ListReviewsByUserIdReq struct {
	UserId   uuid.UUID
	PageSize int
	AfterId  uuid.UUID
}
