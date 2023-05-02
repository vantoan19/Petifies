package petifiesproposalservice

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	eventModels "github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/events"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/publishers"
	petifiesproposaleventkafka "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/publishers/petifies_proposal_event/kafka"
	petifiessessioneventkafka "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/publishers/petifies_session_event/kafka"
	petifiesmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies/mongo"
	petifiesproposalmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies_proposal/mongo"
	petifiessessionmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies_session/mongo"
	reviewmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/review/mongo"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

var logger = logging.New("PetifiesService.PetifiesSessionSvc")

type petifiesProposalService struct {
	mongoClient                    *mongo.Client
	petifiesRepo                   petifiesaggre.PetifiesRepository
	petifiesSessionRepo            petifiessessionaggre.PetifiesSessionRepository
	petifiesProposalRepo           petifiesproposalaggre.PetifiesProposalRepository
	reviewRepo                     reviewaggre.ReviewRepository
	petifiesProposalEventPublisher publishers.PetifiesProposalEventMessagePublisher
	petifiesSessionEventPublisher  publishers.PetifiesSessionEventMessagePublisher
}

type PetifiesProposalConfiguration func(ps *petifiesProposalService) error

type PetifesProposalService interface {
	CreatePetifiesProposal(ctx context.Context, req *models.CreatePetifiesProposalReq) (*petifiesproposalaggre.PetifiesProposalAggre, error)
	EditPetifiesProposal(ctx context.Context, req *models.EditPetifiesProposalReq) (*petifiesproposalaggre.PetifiesProposalAggre, error)
	GetPetifiesProposalById(ctx context.Context, id uuid.UUID) (*petifiesproposalaggre.PetifiesProposalAggre, error)
	ListPetifiesProposalsByIds(ctx context.Context, ids []uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error)
	ListPetifiesProposalsBySessionId(ctx context.Context, sessionID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error)
	ListPetifiesProposalsByUserId(ctx context.Context, userId uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error)
	CancelPetifiesProposal(ctx context.Context, req *models.CancelProposalReq) error
}

func NewPetifiesProposalService(cfgs ...PetifiesProposalConfiguration) (PetifesProposalService, error) {
	ps := &petifiesProposalService{}
	for _, cfg := range cfgs {
		if err := cfg(ps); err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func WithMongoPetifiesRepository(client *mongo.Client) PetifiesProposalConfiguration {
	return func(ps *petifiesProposalService) error {
		repo := petifiesmongo.New(client)
		ps.petifiesRepo = repo
		return nil
	}
}

func WithMongoPetifiesSessionRepository(client *mongo.Client) PetifiesProposalConfiguration {
	return func(ps *petifiesProposalService) error {
		repo := petifiessessionmongo.New(client)
		ps.petifiesSessionRepo = repo
		return nil
	}
}

func WithMongoPetifiesProposalRepository(client *mongo.Client) PetifiesProposalConfiguration {
	return func(ps *petifiesProposalService) error {
		repo := petifiesproposalmongo.New(client)
		ps.petifiesProposalRepo = repo
		return nil
	}
}

func WithMongoReviewRepository(client *mongo.Client) PetifiesProposalConfiguration {
	return func(ps *petifiesProposalService) error {
		repo := reviewmongo.New(client)
		ps.reviewRepo = repo
		return nil
	}
}

func WithKafkaPetifiesProposalPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) PetifiesProposalConfiguration {
	return func(ps *petifiesProposalService) error {
		publisher := petifiesproposaleventkafka.NewPetifiesProposalEventPublisher(producer, repo)
		ps.petifiesProposalEventPublisher = publisher
		return nil
	}
}

func WithKafkaPetifiesSessionEventPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) PetifiesProposalConfiguration {
	return func(ps *petifiesProposalService) error {
		publisher := petifiessessioneventkafka.NewPetifiesSessionEventPublisher(producer, repo)
		ps.petifiesSessionEventPublisher = publisher
		return nil
	}
}

func (ps *petifiesProposalService) CreatePetifiesProposal(ctx context.Context, req *models.CreatePetifiesProposalReq) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	logger.Info("Start CreatePetifiesProposal")

	if exists, err := ps.petifiesProposalRepo.ExistsBySessionAndUserID(ctx, req.PetifiesSessionID, req.UserID); err != nil {
		logger.ErrorData("Finish CreatePetifiesProposal: Failed", logging.Data{"error": err.Error()})
		return nil, err
	} else if exists {
		return nil, status.Errorf(codes.AlreadyExists, "user already proposed for this petifies session")
	}

	petifySession, err := ps.petifiesSessionRepo.GetByID(ctx, req.PetifiesSessionID)
	if err != nil {
		logger.ErrorData("Finish CreatePetifiesProposal: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	petifiesProposal, err := petifiesproposalaggre.NewPetifiesProposalAggregate(
		uuid.New(),
		req.UserID,
		req.PetifiesSessionID,
		req.Proposal,
		valueobjects.PetifiesProposalStatusWaitingForAcceptance,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		logger.ErrorData("Finish CreatePetifiesProposal: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	savedProposal, err := petifySession.AddProposalToSession(*petifiesProposal, ps.petifiesProposalRepo)
	if err != nil {
		logger.ErrorData("Finish CreatePetifiesProposal: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = ps.petifiesProposalEventPublisher.Publish(ctx, eventModels.PetifiesProposalEvent(events.NewPetifiesProposalCreatedEvent(savedProposal)))
	if err != nil {
		logger.ErrorData("Finished CreatePetifiesProposal: Failed to publish event", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreatePetifiesProposal: Successful")
	return savedProposal, nil
}

func (ps *petifiesProposalService) EditPetifiesProposal(ctx context.Context, req *models.EditPetifiesProposalReq) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	logger.Info("Start EditPetifiesProposal")

	proposal, err := ps.petifiesProposalRepo.GetByID(ctx, req.ID)
	if err != nil {
		logger.ErrorData("Finish EditPetifiesProposal: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	proposal.SetProposal(req.Proposal)

	updatedProposal, err := ps.petifiesProposalRepo.Update(ctx, *proposal)
	if err != nil {
		logger.ErrorData("Finish EditPetifiesProposal: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = ps.petifiesProposalEventPublisher.Publish(ctx, eventModels.PetifiesProposalEvent(events.NewPetifiesProposalUpdatedEvent(updatedProposal)))
	if err != nil {
		logger.ErrorData("Finished EditPetifiesProposal: Failed to publish event", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish EditPetifiesProposal: Successful")
	return updatedProposal, nil
}

func (ps *petifiesProposalService) GetPetifiesProposalById(ctx context.Context, id uuid.UUID) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	logger.Info("Start GetPetifiesProposalById")

	proposal, err := ps.petifiesProposalRepo.GetByID(ctx, id)
	if err != nil {
		logger.ErrorData("Finish GetPetifiesProposalById: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetPetifiesProposalById: Successful")
	return proposal, nil
}

func (ps *petifiesProposalService) ListPetifiesProposalsByIds(ctx context.Context, ids []uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error) {
	logger.Info("Start ListPetifiesProposalsByIds")

	proposals, err := ps.petifiesProposalRepo.ListByIDs(ctx, ids)
	if err != nil {
		logger.ErrorData("Finish ListPetifiesProposalsByIds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListPetifiesProposalsByIds: Successful")
	return proposals, nil
}

func (ps *petifiesProposalService) ListPetifiesProposalsBySessionId(ctx context.Context, sessionID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error) {
	logger.Info("Start ListPetifiesProposalsBySessionId")

	proposals, err := ps.petifiesProposalRepo.GetBySessionID(ctx, sessionID, pageSize, afterID)
	if err != nil {
		logger.ErrorData("Finish ListPetifiesProposalsBySessionId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListPetifiesProposalsBySessionId: Successful")
	return proposals, nil
}

func (ps *petifiesProposalService) ListPetifiesProposalsByUserId(ctx context.Context, userId uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error) {
	logger.Info("Start ListPetifiesProposalsByUserId")

	proposals, err := ps.petifiesProposalRepo.GetByUserID(ctx, userId, pageSize, afterID)
	if err != nil {
		logger.ErrorData("Finish ListPetifiesProposalsByUserId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListPetifiesProposalsByUserId: Successful")
	return proposals, nil
}

func (ps *petifiesProposalService) CancelPetifiesProposal(ctx context.Context, req *models.CancelProposalReq) error {
	logger.Info("Start CancelPetifiesProposal")
	err := dbutils.ExecWithSession(ctx, ps.mongoClient, func(ssCtx mongo.SessionContext) error {
		proposal, err := ps.petifiesProposalRepo.GetByIDWithSession(ssCtx, req.ProposalId)
		if err != nil {
			return err
		}
		if proposal.GetUserID() != req.UserId {
			return status.Errorf(codes.PermissionDenied, "user does not have permission to cancel the proposal")
		}

		session, err := ps.petifiesSessionRepo.GetByIDWithSession(ssCtx, proposal.GetPetifiesSessionID())
		if err != nil {
			return err
		}

		sessionStatusChanged := false
		var sessionNextStatus valueobjects.PetifiesSessionStatus
		if proposal.GetStatus() == valueobjects.PetifiesProposalStatusAccepted {
			err = session.ToWaitForProposalsStatus()
			if err != nil {
				return err
			}
			sessionStatusChanged = true
			sessionNextStatus = valueobjects.PetifiesSessionStatusWaitingForProposal
		}
		err = proposal.ToCancelledStatus()
		if err != nil {
			return err
		}

		updatedProposal, err := ps.petifiesProposalRepo.Update(ssCtx, *proposal)
		if err != nil {
			return err
		}
		if updatedProposal.GetStatus() != valueobjects.PetifiesProposalStatusRejected {
			return status.Errorf(codes.Internal, "failed to change proposal's status")
		}

		err = ps.petifiesProposalEventPublisher.Publish(ssCtx, eventModels.PetifiesProposalEvent(
			events.NewPetifiesProposalAcceptedEvent(updatedProposal)))
		if err != nil {
			return err
		}

		if sessionStatusChanged {
			updatedSession, err := ps.petifiesSessionRepo.Update(ssCtx, *session)
			if err != nil {
				return err
			}
			if updatedSession.GetStatus() != sessionNextStatus {
				return status.Errorf(codes.Internal, "failed to change session's status")
			}

			err = ps.petifiesSessionEventPublisher.Publish(ssCtx, eventModels.PetifiesSessionEvent(
				events.NewPetifiesSessionWaitingForProposalsEvent(updatedSession)))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish CancelPetifiesProposal: Failed", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finish CancelPetifiesProposal: Successful")
	return nil
}
