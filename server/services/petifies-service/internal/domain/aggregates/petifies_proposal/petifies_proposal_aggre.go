package petifiesproposalaggre

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

var (
	ErrNeedWaitForAcceptanceStatus    = status.Errorf(codes.InvalidArgument, "the proposal needs to be in WAITING_FOR_ACCEPTANCE status")
	ErrInvalidPrecedeStatus           = status.Errorf(codes.InvalidArgument, "the proposal has invalid preceding status")
	ErrStatusAlreadyWaitForAcceptance = status.Errorf(codes.Aborted, "status is already WAITING_FOR_ACCEPTANCE")
	ErrStatusAlreadyAccepted          = status.Errorf(codes.Aborted, "status is already ACCEPTED")
	ErrStatusAlreadyCancelled         = status.Errorf(codes.Aborted, "status is already CANCELLED")
	ErrStatusAlreadyRejected          = status.Errorf(codes.Aborted, "status is already REJECTED")
	ErrStatusAlreadySessionClosed     = status.Errorf(codes.Aborted, "status is already SESSION_CLOSED")
)

type PetifiesProposalAggre struct {
	proposal *entities.PetifiesProposal
}

func NewPetifiesProposalAggregate(
	id,
	userID,
	petifiesSessionID uuid.UUID,
	proposal string,
	status valueobjects.PetifiesProposalStatus,
	createdAt,
	updatedAt time.Time,
) (*PetifiesProposalAggre, error) {
	petifiesProposal := &entities.PetifiesProposal{
		ID:                id,
		UserID:            userID,
		PetifiesSessionID: petifiesSessionID,
		Proposal:          proposal,
		Status:            status,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	if errs := petifiesProposal.Validate(); errs.Exist() {
		return nil, errs[0]
	}

	return &PetifiesProposalAggre{
		proposal: petifiesProposal,
	}, nil
}

func (p *PetifiesProposalAggre) ToAcceptedStatus() error {
	if p.proposal.Status != valueobjects.PetifiesProposalStatusWaitingForAcceptance {
		return ErrNeedWaitForAcceptanceStatus
	}
	if p.proposal.Status == valueobjects.PetifiesProposalStatusAccepted {
		return ErrStatusAlreadyAccepted
	}

	p.proposal.Status = valueobjects.PetifiesProposalStatusAccepted
	return nil
}

func (p *PetifiesProposalAggre) ToCancelledStatus() error {
	switch p.proposal.Status {
	case valueobjects.PetifiesProposalStatusAccepted,
		valueobjects.PetifiesProposalStatusWaitingForAcceptance:
		p.proposal.Status = valueobjects.PetifiesProposalStatusCancelled
		return nil
	case valueobjects.PetifiesProposalStatusCancelled:
		return ErrStatusAlreadyCancelled
	default:
		return ErrInvalidPrecedeStatus
	}
}

func (p *PetifiesProposalAggre) ToWaitingForAcceptanceStatus() error {
	switch p.proposal.Status {
	case valueobjects.PetifiesProposalStatusCancelled,
		valueobjects.PetifiesProposalStatusRejected:
		p.proposal.Status = valueobjects.PetifiesProposalStatusWaitingForAcceptance
		return nil
	case valueobjects.PetifiesProposalStatusWaitingForAcceptance:
		return ErrStatusAlreadyWaitForAcceptance
	default:
		return ErrInvalidPrecedeStatus
	}
}

func (p *PetifiesProposalAggre) ToRejectedStatus() error {
	switch p.proposal.Status {
	case valueobjects.PetifiesProposalStatusWaitingForAcceptance,
		valueobjects.PetifiesProposalStatusAccepted:
		p.proposal.Status = valueobjects.PetifiesProposalStatusRejected
		return nil
	case valueobjects.PetifiesProposalStatusRejected:
		return ErrStatusAlreadyRejected
	default:
		return ErrInvalidPrecedeStatus
	}
}

func (p *PetifiesProposalAggre) ToSessionClosedStatus() error {
	switch p.proposal.Status {
	case valueobjects.PetifiesProposalStatusWaitingForAcceptance:
		p.proposal.Status = valueobjects.PetifiesProposalStatusSessionClosed
		return nil
	case valueobjects.PetifiesProposalStatusSessionClosed:
		return ErrStatusAlreadySessionClosed
	default:
		return ErrInvalidPrecedeStatus
	}
}

// ========= Aggregate Root Getters =========

func (p *PetifiesProposalAggre) GetID() uuid.UUID {
	return p.proposal.ID
}

func (p *PetifiesProposalAggre) GetUserID() uuid.UUID {
	return p.proposal.UserID
}

func (p *PetifiesProposalAggre) GetPetifiesSessionID() uuid.UUID {
	return p.proposal.PetifiesSessionID
}

func (p *PetifiesProposalAggre) GetProposal() string {
	return p.proposal.Proposal
}

func (p *PetifiesProposalAggre) SetProposal(proposal string) error {
	if len(strings.TrimSpace(proposal)) == 0 {
		return entities.ErrReviewRequired
	}
	if len(proposal) > 5000 {
		return entities.ErrReviewExceedsLimit
	}

	p.proposal.Proposal = proposal
	return nil
}

func (p *PetifiesProposalAggre) GetStatus() valueobjects.PetifiesProposalStatus {
	return p.proposal.Status
}

func (p *PetifiesProposalAggre) GetCreatedAt() time.Time {
	return p.proposal.CreatedAt
}

func (p *PetifiesProposalAggre) GetUpdatedAt() time.Time {
	return p.proposal.UpdatedAt
}
