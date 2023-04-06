package models

import (
	"time"

	"github.com/google/uuid"
)

type Petifies struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PetName     string    `json:"pet_name"`
	Images      []Image   `json:"images"`
	Status      string    `json:"status"`
	Address     Address   `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PetifiesSession struct {
	ID         uuid.UUID `json:"id"`
	PetifiesID uuid.UUID `json:"petifies_id"`
	FromTime   time.Time `json:"from_time"`
	ToTime     time.Time `json:"to_time"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PetifiesProposal struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	PetifiesSessionID uuid.UUID `json:"petifies_session_id"`
	Proposal          string    `json:"proposal"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Review struct {
	ID         uuid.UUID `json:"id"`
	PetifiesID uuid.UUID `json:"petifies_id"`
	AuthorID   uuid.UUID `json:"author_id"`
	Review     string    `json:"review"`
	Image      Image     `json:"image"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ManyPetifies struct {
	Petifies []*Petifies
}

type ManyPetifiesSessions struct {
	PetifiesSessions []*PetifiesSession
}

type ManyPetifiesProposals struct {
	PetifiesProposals []*PetifiesProposal
}

type ManyReviews struct {
	Reviews []*Review
}
