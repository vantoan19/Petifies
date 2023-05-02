package entities

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
)

var (
	ErrUserIDRequired            = status.Errorf(codes.InvalidArgument, "ID must be set")
	ErrPetifiesSessionIDRequired = status.Errorf(codes.InvalidArgument, "PetifiesSessionID must be set")
	ErrProposalEmpty             = status.Errorf(codes.InvalidArgument, "Proposal cannot be empty")
	ErrProposalExceedLimit       = status.Errorf(codes.InvalidArgument, "Proposal exceeds 5000 characters")
	ErrStatusUnknown             = status.Errorf(codes.InvalidArgument, "Status has an unknown value")
)

type PetifiesProposal struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	PetifiesSessionID uuid.UUID
	Proposal          string
	Status            valueobjects.PetifiesProposalStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (p PetifiesProposal) Validate() (errs commonutils.MultiError) {
	if p.ID == uuid.Nil {
		errs = append(errs, ErrIDRequired)
	}
	if p.PetifiesSessionID == uuid.Nil {
		errs = append(errs, ErrPetifiesSessionIDRequired)
	}
	if p.UserID == uuid.Nil {
		errs = append(errs, ErrUserIDRequired)
	}
	if p.Proposal == "" {
		errs = append(errs, ErrProposalEmpty)
	}
	if len(p.Proposal) > 5000 {
		errs = append(errs, ErrProposalExceedLimit)
	}
	switch p.Status {
	case valueobjects.PetifiesProposalStatusWaitingForAcceptance,
		valueobjects.PetifiesProposalStatusAccepted,
		valueobjects.PetifiesProposalStatusRejected,
		valueobjects.PetifiesProposalStatusCancelled,
		valueobjects.PetifiesProposalStatusSessionClosed:
		// do nothing, these are valid values
	default:
		errs = append(errs, ErrStatusUnknown)
	}

	return errs
}
