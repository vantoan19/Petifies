package models

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	AddressLineOne string  `json:"address_line_one"`
	AddressLineTwo string  `json:"address_line_two"`
	Street         string  `json:"street"`
	District       string  `json:"district"`
	City           string  `json:"city"`
	Region         string  `json:"region"`
	PostalCode     string  `json:"postal_code"`
	Country        string  `json:"country"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
}

type Image struct {
	URI         string `json:"uri"`
	Description string `json:"description"`
}

type CreatePetifiesReq struct {
	OwnerID     uuid.UUID
	Type        string
	Title       string
	Description string
	PetName     string
	Images      []Image
	Address     Address
}

type CreatePetifiesSessionReq struct {
	PetifiesID uuid.UUID
	FromTime   time.Time
	ToTime     time.Time
}

type CreatePetifiesProposalReq struct {
	UserID            uuid.UUID
	PetifiesSessionID uuid.UUID
	Proposal          string
}

type CreateReviewReq struct {
	PetifiesID uuid.UUID
	AuthorID   uuid.UUID
	Review     string
	Image      Image
}

type EditPetifiesReq struct {
	ID uuid.UUID
	// Type        string
	Title       string
	Description string
	PetName     string
	Images      []Image
	Address     Address
}

type EditPetifiesSessionReq struct {
	ID       uuid.UUID
	FromTime time.Time
	ToTime   time.Time
}

type EditPetifiesProposalReq struct {
	ID       uuid.UUID
	Proposal string
}

type EditReviewReq struct {
	ID     uuid.UUID
	Review string
	Image  Image
}

type GetPetifiesByIdReq struct {
	ID uuid.UUID
}

type ListPetifiesByIdsReq struct {
	PetifiesIDs []uuid.UUID
}

type ListPetifiesByOwnerIdReq struct {
	OwnerID  uuid.UUID
	PageSize int
	AfterID  uuid.UUID
}

type GetSessionByIdReq struct {
	ID uuid.UUID
}

type ListSessionsByPetifiesIdReq struct {
	PetifiesID uuid.UUID
	PageSize   int
	AfterID    uuid.UUID
}

type ListSessionsByIdsReq struct {
	PetifiesSessionIDs []uuid.UUID
}

type GetProposalByIdReq struct {
	ID uuid.UUID
}

type ListProposalsBySessionIdReq struct {
	PetifiesSessionID uuid.UUID
	PageSize          int
	AfterID           uuid.UUID
}

type ListProposalsByIdsReq struct {
	PetifiesProposalIDs []uuid.UUID
}

type GetReviewByIdReq struct {
	ID uuid.UUID
}

type ListReviewsByPetifiesIdReq struct {
	PetifiesID uuid.UUID
	PageSize   int
	AfterID    uuid.UUID
}

type ListReviewsByIdsReq struct {
	ReviewIDs []uuid.UUID
}
