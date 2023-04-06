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
	ErrIDRequired           = status.Errorf(codes.InvalidArgument, "ID must be set")
	ErrPetifiesIDRequired   = status.Errorf(codes.InvalidArgument, "PetifiesID must be set")
	ErrFromTimeInvalid      = status.Errorf(codes.InvalidArgument, "FromTime must come before ToTime")
	ErrUnknownSessionStatus = status.Errorf(codes.InvalidArgument, "Status has an unknown value")
)

type PetifiesSession struct {
	ID         uuid.UUID
	PetifiesID uuid.UUID
	Time       valueobjects.TimeRange
	Status     valueobjects.PetifiesSessionStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (s PetifiesSession) Validate() (errs commonutils.MultiError) {
	if s.ID == uuid.Nil {
		errs = append(errs, ErrIDRequired)
	}
	if s.PetifiesID == uuid.Nil {
		errs = append(errs, ErrPetifiesIDRequired)
	}
	if errs_ := s.Time.Validate(); errs.Exist() {
		errs = append(errs, errs_...)
	}
	switch s.Status {
	case valueobjects.PetifiesSessionStatusWaitingForProposal,
		valueobjects.PetifiesSessionStatusProposalAccepted,
		valueobjects.PetifiesSessionStatusOnGoing,
		valueobjects.PetifiesSessionStatusEnded:
		// do nothing
	default:
		errs = append(errs, ErrUnknownSessionStatus)
	}

	return errs
}
