package petifiesaggre

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

var (
	ErrInvalidSessionTime         = status.Errorf(codes.InvalidArgument, "time of a session is not allowed to intersect with the existing ones")
	ErrExpiredTimeRange           = status.Errorf(codes.InvalidArgument, "time of a newly created session have to be in the future")
	ErrMismatchPetifiesID         = status.Errorf(codes.InvalidArgument, "petifiesID doesn't match with the id of petifies")
	ErrNeedOnGoingStatus          = status.Errorf(codes.InvalidArgument, "the session need to be in ON_GOING status")
	ErrNeedWaitForProposalsStatus = status.Errorf(codes.InvalidArgument, "the session need to be in WAITING_FOR_PROPOSALS status when being created")
)

type PetifiesAggre struct {
	petifies *entities.Petifies
	address  *entities.Address
	sessions []uuid.UUID
	reviews  []uuid.UUID
}

func NewPetifiesAggregate(
	id,
	ownerID uuid.UUID,
	petifiesType valueobjects.PetifiesType,
	title,
	description,
	petName string,
	images []valueobjects.Image,
	status valueobjects.PetifiesStatus,
	createdAt,
	updatedAt time.Time,
	address entities.Address,
) (*PetifiesAggre, error) {
	petifies := &entities.Petifies{
		ID:          id,
		OwnerID:     ownerID,
		Type:        petifiesType,
		Title:       title,
		Description: description,
		PetName:     petName,
		Images:      images,
		Status:      status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	fmt.Println(petifies)
	fmt.Println(address)

	if errs := petifies.Validate(); errs.Exist() {
		return nil, errs[0]
	}
	if errs := address.Validate(); errs.Exist() {
		return nil, errs[0]
	}

	return &PetifiesAggre{
		petifies: petifies,
		address:  &address,
		sessions: make([]uuid.UUID, 0),
		reviews:  make([]uuid.UUID, 0),
	}, nil
}

func (p *PetifiesAggre) AddSessionToPetifies(
	session petifiessessionaggre.PetifiesSessionAggre,
	repo petifiessessionaggre.PetifiesSessionRepository,
) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	if session.GetTime().GetFromTime().Before(time.Now()) {
		return nil, ErrExpiredTimeRange
	}
	if session.GetStatus() != valueobjects.PetifiesSessionStatusWaitingForProposal {
		return nil, ErrNeedWaitForProposalsStatus
	}
	if session.GetPetifiesID() != p.petifies.ID {
		return nil, ErrMismatchPetifiesID
	}

	existingSessions, err := repo.GetByPetifiesID(context.Background(), p.petifies.ID, 0, uuid.Nil)
	if err != nil {
		return nil, err
	}
	for _, s := range existingSessions {
		if s.GetTime().Intersects(session.GetTime()) {
			return nil, ErrInvalidSessionTime
		}
	}

	savedSession, err := repo.Save(context.Background(), session)
	if err != nil {
		return nil, err
	}

	p.sessions = append(p.sessions, savedSession.GetID())
	return savedSession, nil
}

func (p *PetifiesAggre) AddReview(
	review reviewaggre.ReviewAggre,
	repo reviewaggre.ReviewRepository,
) (*reviewaggre.ReviewAggre, error) {
	if review.GetPetifiesID() != p.petifies.ID {
		return nil, ErrMismatchPetifiesID
	}

	savedReview, err := repo.Save(context.Background(), review)
	if err != nil {
		return nil, err
	}

	p.reviews = append(p.reviews, savedReview.GetID())
	return savedReview, nil
}

// =========== Root Aggregate Getters ===========

func (p *PetifiesAggre) GetID() uuid.UUID {
	return p.petifies.ID
}

func (p *PetifiesAggre) GetOwnerID() uuid.UUID {
	return p.petifies.OwnerID
}

func (p *PetifiesAggre) GetType() valueobjects.PetifiesType {
	return p.petifies.Type
}

func (p *PetifiesAggre) GetTitle() string {
	return p.petifies.Title
}

func (p *PetifiesAggre) SetTitle(title string) error {
	if len(strings.TrimSpace(title)) == 0 {
		return entities.ErrTitleRequired
	}
	if len(title) > 500 {
		return entities.ErrTitleExceedsLimit
	}
	p.petifies.Title = title
	return nil
}

func (p *PetifiesAggre) GetDescription() string {
	return p.petifies.Description
}

func (p *PetifiesAggre) SetDescription(description string) error {
	if len(strings.TrimSpace(description)) == 0 {
		return entities.ErrDescriptionRequired
	}
	if len(description) > 5000 {
		return entities.ErrDescriptionExceedsLimit
	}
	p.petifies.Description = description
	return nil
}

func (p *PetifiesAggre) GetPetName() string {
	return p.petifies.PetName
}

func (p *PetifiesAggre) SetPetName(petName string) error {
	if len(strings.TrimSpace(petName)) == 0 {
		return entities.ErrPetNameRequired
	}
	if len(petName) > 5000 {
		return entities.ErrPetNameExceedsLimit
	}
	p.petifies.PetName = petName
	return nil
}

func (p *PetifiesAggre) GetImages() []valueobjects.Image {
	return p.petifies.Images
}

func (p *PetifiesAggre) SetImages(images []valueobjects.Image) error {
	for _, image := range images {
		if errs := image.Validate(); errs != nil {
			return errs[0]
		}
	}
	p.petifies.Images = images
	return nil
}

func (p *PetifiesAggre) GetStatus() valueobjects.PetifiesStatus {
	return p.petifies.Status
}

func (p *PetifiesAggre) SetStatus(status valueobjects.PetifiesStatus) error {
	switch p.petifies.Status {
	case valueobjects.PetifiesDeleted,
		valueobjects.PetifiesOnASession,
		valueobjects.PetifiesUnavailable:
		p.petifies.Status = status
		return nil
	default:
		return entities.ErrUnknownPetifiesStatus
	}
}

func (p *PetifiesAggre) GetCreatedAt() time.Time {
	return p.petifies.CreatedAt
}

func (p *PetifiesAggre) GetUpdatedAt() time.Time {
	return p.petifies.UpdatedAt
}

func (p *PetifiesAggre) GetAddress() entities.Address {
	return *p.address
}

func (p *PetifiesAggre) SetAddress(address entities.Address) error {
	if errs := address.Validate(); errs.Exist() {
		return errs[0]
	}
	p.address = &address
	return nil
}
