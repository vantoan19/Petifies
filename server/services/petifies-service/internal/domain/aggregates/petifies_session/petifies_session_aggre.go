package petifiessessionaggre

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

var (
	ErrExceedProposalLimit          = status.Errorf(codes.AlreadyExists, "a user can only create a proposal")
	ErrDifferentSessionID           = status.Errorf(codes.InvalidArgument, "proposal's sessionID is different than the id of session")
	ErrInvalidNewlyCreatedStatus    = status.Errorf(codes.InvalidArgument, "newly created proposal's status have to be WAITING_FOR_ACCEPTANCE")
	ErrNeedWaitForProposalsStatus   = status.Errorf(codes.InvalidArgument, "the session need to be in WAITING_FOR_PROPOSALS status")
	ErrNeedProposalAcceptanceStatus = status.Errorf(codes.InvalidArgument, "the session need to be in PROPOSAL_ACCEPTED status")
	ErrNeedOnGoingStatus            = status.Errorf(codes.InvalidArgument, "the session need to be in ON_GOING status")
)

type PetifiesSessionAggre struct {
	session   *entities.PetifiesSession
	proposals []uuid.UUID
}

func NewPetitifesSessionAggre(
	id,
	petifiesID uuid.UUID,
	fromTime,
	toTime time.Time,
	status valueobjects.PetifiesSessionStatus,
	createdAt,
	updatedAt time.Time,
) (*PetifiesSessionAggre, error) {
	petifiesSession := &entities.PetifiesSession{
		ID:         id,
		PetifiesID: petifiesID,
		Time:       valueobjects.NewTimeRange(fromTime, toTime),
		Status:     status,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	if errs := petifiesSession.Validate(); errs.Exist() {
		return nil, errs[0]
	}

	return &PetifiesSessionAggre{
		session:   petifiesSession,
		proposals: make([]uuid.UUID, 0),
	}, nil
}

func (p *PetifiesSessionAggre) ToProposalAcceptedStatus() error {
	if p.session.Status != valueobjects.PetifiesSessionStatusWaitingForProposal {
		return ErrNeedWaitForProposalsStatus
	}

	p.session.Status = valueobjects.PetifiesSessionStatusProposalAccepted
	return nil
}

func (p *PetifiesSessionAggre) ToOnGoingStatus() error {
	if p.session.Status != valueobjects.PetifiesSessionStatusProposalAccepted {
		return ErrNeedProposalAcceptanceStatus
	}

	p.session.Status = valueobjects.PetifiesSessionStatusOnGoing
	return nil
}

func (p *PetifiesSessionAggre) ToWaitForProposalsStatus() error {
	if p.session.Status != valueobjects.PetifiesSessionStatusProposalAccepted {
		return ErrNeedProposalAcceptanceStatus
	}

	p.session.Status = valueobjects.PetifiesSessionStatusWaitingForProposal
	return nil
}

func (p *PetifiesSessionAggre) ToEndedStaus() error {
	p.session.Status = valueobjects.PetifiesSessionStatusEnded
	return nil
}

func (p *PetifiesSessionAggre) AddProposalToSession(
	proposal petifiesproposalaggre.PetifiesProposalAggre,
	repo petifiesproposalaggre.PetifiesProposalRepository,
) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	if proposal.GetPetifiesSessionID() != p.session.ID {
		return nil, ErrDifferentSessionID
	}
	if proposal.GetStatus() != valueobjects.PetifiesProposalStatusWaitingForAcceptance {
		return nil, ErrInvalidNewlyCreatedStatus
	}
	if exists, err := repo.ExistsBySessionAndUserID(context.Background(), p.session.ID, proposal.GetID()); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrExceedProposalLimit
	}

	savedProposal, err := repo.Save(context.Background(), proposal)
	if err != nil {
		return nil, err
	}

	p.proposals = append(p.proposals, savedProposal.GetID())
	return savedProposal, nil
}

// ========= Aggregate Root Getters / Setters =========

func (p *PetifiesSessionAggre) GetID() uuid.UUID {
	return p.session.ID
}

func (p *PetifiesSessionAggre) GetPetifiesID() uuid.UUID {
	return p.session.PetifiesID
}

func (p *PetifiesSessionAggre) GetTime() valueobjects.TimeRange {
	return p.session.Time
}

func (p *PetifiesSessionAggre) GetStatus() valueobjects.PetifiesSessionStatus {
	return p.session.Status
}

func (p *PetifiesSessionAggre) GetCreatedAt() time.Time {
	return p.session.CreatedAt
}

func (p *PetifiesSessionAggre) GetUpdatedAt() time.Time {
	return p.session.UpdatedAt
}

func (p *PetifiesSessionAggre) SetTime(time valueobjects.TimeRange) error {
	if errs := time.Validate(); errs.Exist() {
		return errs[0]
	}
	p.session.Time = time
	return nil
}
