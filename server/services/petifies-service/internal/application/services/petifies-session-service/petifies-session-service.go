package petifiessessionservice

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

type petifiesSessionService struct {
	mongoClient                    *mongo.Client
	petifiesRepo                   petifiesaggre.PetifiesRepository
	petifiesSessionRepo            petifiessessionaggre.PetifiesSessionRepository
	petifiesProposalRepo           petifiesproposalaggre.PetifiesProposalRepository
	reviewRepo                     reviewaggre.ReviewRepository
	petifiesSessionEventPublisher  publishers.PetifiesSessionEventMessagePublisher
	petifiesProposalEventPublisher publishers.PetifiesProposalEventMessagePublisher
}

type PetifiesSessionConfiguration func(ps *petifiesSessionService) error

type PetifesSessionService interface {
	CreatePetifiesSession(ctx context.Context, req *models.CreatePetifiesSessionReq) (*petifiessessionaggre.PetifiesSessionAggre, error)
	EditPetifiesSession(ctx context.Context, req *models.EditPetifiesSessionReq) (*petifiessessionaggre.PetifiesSessionAggre, error)
	GetById(ctx context.Context, id uuid.UUID) (*petifiessessionaggre.PetifiesSessionAggre, error)
	ListByIds(ctx context.Context, ids []uuid.UUID) ([]*petifiessessionaggre.PetifiesSessionAggre, error)
	ListByPetifiesId(ctx context.Context, petifiesId uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiessessionaggre.PetifiesSessionAggre, error)
	AcceptProposal(ctx context.Context, req *models.AcceptProposalReq) error
	RejectProposal(ctx context.Context, req *models.RejectProposalReq) error
}

func NewPetifiesSessionService(cfgs ...PetifiesSessionConfiguration) (PetifesSessionService, error) {
	ps := &petifiesSessionService{}
	for _, cfg := range cfgs {
		if err := cfg(ps); err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func WithMongoPetifiesRepository(client *mongo.Client) PetifiesSessionConfiguration {
	return func(ps *petifiesSessionService) error {
		repo := petifiesmongo.New(client)
		ps.petifiesRepo = repo
		return nil
	}
}

func WithMongoPetifiesSessionRepository(client *mongo.Client) PetifiesSessionConfiguration {
	return func(ps *petifiesSessionService) error {
		repo := petifiessessionmongo.New(client)
		ps.petifiesSessionRepo = repo
		return nil
	}
}

func WithMongoPetifiesProposalRepository(client *mongo.Client) PetifiesSessionConfiguration {
	return func(ps *petifiesSessionService) error {
		repo := petifiesproposalmongo.New(client)
		ps.petifiesProposalRepo = repo
		return nil
	}
}

func WithMongoReviewRepository(client *mongo.Client) PetifiesSessionConfiguration {
	return func(ps *petifiesSessionService) error {
		repo := reviewmongo.New(client)
		ps.reviewRepo = repo
		return nil
	}
}

func WithKafkaPetifiesSessionEventPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) PetifiesSessionConfiguration {
	return func(ps *petifiesSessionService) error {
		publisher := petifiessessioneventkafka.NewPetifiesSessionEventPublisher(producer, repo)
		ps.petifiesSessionEventPublisher = publisher
		return nil
	}
}

func WithKafkaPetifiesProposalEventPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) PetifiesSessionConfiguration {
	return func(ps *petifiesSessionService) error {
		publisher := petifiesproposaleventkafka.NewPetifiesProposalEventPublisher(producer, repo)
		ps.petifiesProposalEventPublisher = publisher
		return nil
	}
}

func (ps *petifiesSessionService) CreatePetifiesSession(ctx context.Context, req *models.CreatePetifiesSessionReq) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	logger.Info("Start CreatePetifiesSession")

	petify, err := ps.petifiesRepo.GetByID(ctx, req.PetifiesID)
	if err != nil {
		logger.ErrorData("Finish CreatePetifiesSession: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	petifiesSession, err := petifiessessionaggre.NewPetitifesSessionAggre(
		uuid.New(),
		req.PetifiesID,
		req.FromTime,
		req.ToTime,
		valueobjects.PetifiesSessionStatusWaitingForProposal,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		logger.ErrorData("Finish CreatePetifiesSession: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	savedSession, err := petify.AddSessionToPetifies(*petifiesSession, ps.petifiesSessionRepo)
	if err != nil {
		logger.ErrorData("Finish CreatePetifiesSession: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = ps.petifiesSessionEventPublisher.Publish(ctx, eventModels.PetifiesSessionEvent(events.NewPetifiesSessionCreatedEvent(savedSession)))
	if err != nil {
		logger.ErrorData("Finished CreatePetifiesSession: Failed to publish event", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreatePetifiesSession: Successful")
	return savedSession, nil
}

func (ps *petifiesSessionService) EditPetifiesSession(ctx context.Context, req *models.EditPetifiesSessionReq) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	logger.Info("Start EditPetifiesSession")

	petifySession, err := ps.petifiesSessionRepo.GetByID(ctx, req.ID)
	if err != nil {
		logger.ErrorData("Finish EditPetifiesSession: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	petifySession.SetTime(valueobjects.NewTimeRange(req.FromTime, req.ToTime))

	updatedSession, err := ps.petifiesSessionRepo.Update(ctx, *petifySession)
	if err != nil {
		logger.ErrorData("Finish EditPetifiesSession: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = ps.petifiesSessionEventPublisher.Publish(ctx, eventModels.PetifiesSessionEvent(events.NewPetifiesSessionUpdatedEvent(updatedSession)))
	if err != nil {
		logger.ErrorData("Finished CreatePetifiesSession: Failed to publish event", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish EditPetifiesSession: Successful")
	return updatedSession, nil
}

func (ps *petifiesSessionService) GetById(ctx context.Context, id uuid.UUID) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	logger.Info("Start GetById")

	session, err := ps.petifiesSessionRepo.GetByID(ctx, id)
	if err != nil {
		logger.ErrorData("Finish GetById: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetById: Successful")
	return session, nil
}

func (ps *petifiesSessionService) ListByIds(ctx context.Context, ids []uuid.UUID) ([]*petifiessessionaggre.PetifiesSessionAggre, error) {
	logger.Info("Start ListByIds")

	sessions, err := ps.petifiesSessionRepo.ListByIds(ctx, ids)
	if err != nil {
		logger.ErrorData("Finish ListByIds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListByIds: Successful")
	return sessions, nil
}
func (ps *petifiesSessionService) ListByPetifiesId(ctx context.Context, petifiesId uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiessessionaggre.PetifiesSessionAggre, error) {
	logger.Info("Start ListByPetifiesId")

	sessions, err := ps.petifiesSessionRepo.GetByPetifiesID(ctx, petifiesId, pageSize, afterID)
	if err != nil {
		logger.ErrorData("Finish ListByPetifiesId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListByPetifiesId: Successful")
	return sessions, nil
}

func (ps *petifiesSessionService) AcceptProposal(ctx context.Context, req *models.AcceptProposalReq) error {
	logger.Info("Start AcceptProposal")
	err := dbutils.ExecWithSession(ctx, ps.mongoClient, func(ssCtx mongo.SessionContext) error {
		proposal, err := ps.petifiesProposalRepo.GetByIDWithSession(ssCtx, req.ProposalId)
		if err != nil {
			return err
		}
		if proposal.GetPetifiesSessionID() != req.SessionId {
			return status.Errorf(codes.PermissionDenied, "proposal does not belong to the session")
		}

		session, err := ps.petifiesSessionRepo.GetByIDWithSession(ssCtx, req.SessionId)
		if err != nil {
			return nil
		}
		petify, err := ps.petifiesRepo.GetByIDWithSession(ssCtx, session.GetPetifiesID())
		if err != nil {
			return nil
		}

		if petify.GetOwnerID() != req.UserId {
			return status.Errorf(codes.PermissionDenied, "user does not have the permission to accept proposal")
		}

		err = proposal.ToAcceptedStatus()
		if err != nil {
			return err
		}
		err = session.ToProposalAcceptedStatus()
		if err != nil {
			return err
		}

		updatedProposal, err := ps.petifiesProposalRepo.Update(ssCtx, *proposal)
		if err != nil {
			return err
		}
		updatedSession, err := ps.petifiesSessionRepo.Update(ssCtx, *session)
		if err != nil {
			return err
		}
		if updatedProposal.GetStatus() != valueobjects.PetifiesProposalStatusAccepted {
			return status.Errorf(codes.Internal, "failed to change proposal's status")
		}
		if updatedSession.GetStatus() != valueobjects.PetifiesSessionStatusProposalAccepted {
			return status.Errorf(codes.Internal, "failed to change session's status")
		}

		err = ps.petifiesSessionEventPublisher.Publish(ssCtx, eventModels.PetifiesSessionEvent(
			events.NewPetifiesSessionProposalAcceptedEvent(updatedSession)))
		if err != nil {
			return err
		}
		err = ps.petifiesProposalEventPublisher.Publish(ssCtx, eventModels.PetifiesProposalEvent(
			events.NewPetifiesProposalAcceptedEvent(updatedProposal)))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish AcceptProposal: Failed", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finish AcceptProposal: Successful")
	return nil
}

func (ps *petifiesSessionService) RejectProposal(ctx context.Context, req *models.RejectProposalReq) error {
	logger.Info("Start RejectProposal")
	err := dbutils.ExecWithSession(ctx, ps.mongoClient, func(ssCtx mongo.SessionContext) error {
		proposal, err := ps.petifiesProposalRepo.GetByIDWithSession(ssCtx, req.ProposalId)
		if err != nil {
			return err
		}
		if proposal.GetPetifiesSessionID() != req.SessionId {
			return status.Errorf(codes.PermissionDenied, "proposal does not belong to the session")
		}

		session, err := ps.petifiesSessionRepo.GetByIDWithSession(ssCtx, req.SessionId)
		if err != nil {
			return err
		}
		petify, err := ps.petifiesRepo.GetByIDWithSession(ssCtx, session.GetPetifiesID())
		if err != nil {
			return err
		}

		if petify.GetOwnerID() != req.UserId {
			return status.Errorf(codes.PermissionDenied, "user does not have the permission to reject proposal")
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
		err = proposal.ToRejectedStatus()
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
		logger.ErrorData("Finish RejectProposal: Failed", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finish RejectProposal: Successful")
	return nil
}
