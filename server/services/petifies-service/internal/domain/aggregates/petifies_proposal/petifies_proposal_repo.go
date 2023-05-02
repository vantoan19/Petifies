package petifiesproposalaggre

import (
	"context"

	"github.com/google/uuid"
)

type PetifiesProposalRepository interface {
	GetBySessionAndUserID(ctx context.Context, sessionID, userID uuid.UUID) (*PetifiesProposalAggre, error)
	GetByID(ctx context.Context, id uuid.UUID) (*PetifiesProposalAggre, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*PetifiesProposalAggre, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*PetifiesProposalAggre, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*PetifiesProposalAggre, error)
	ExistsBySessionAndUserID(ctx context.Context, sessionID, userID uuid.UUID) (bool, error)
	Save(ctx context.Context, proposal PetifiesProposalAggre) (*PetifiesProposalAggre, error)
	Update(ctx context.Context, proposal PetifiesProposalAggre) (*PetifiesProposalAggre, error)

	GetBySessionAndUserIDWithSession(ctx context.Context, sessionID, userID uuid.UUID) (*PetifiesProposalAggre, error)
	GetByIDWithSession(ctx context.Context, id uuid.UUID) (*PetifiesProposalAggre, error)
	GetBySessionIDWithSession(ctx context.Context, sessionID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*PetifiesProposalAggre, error)
	GetByUserIDWithSession(ctx context.Context, userID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*PetifiesProposalAggre, error)
	ListByIDsWithSession(ctx context.Context, ids []uuid.UUID) ([]*PetifiesProposalAggre, error)
	ExistsBySessionAndUserIDWithSession(ctx context.Context, sessionID, userID uuid.UUID) (bool, error)
	SaveWithSession(ctx context.Context, proposal PetifiesProposalAggre) (*PetifiesProposalAggre, error)
	UpdateWithSession(ctx context.Context, proposal PetifiesProposalAggre) (*PetifiesProposalAggre, error)
}
