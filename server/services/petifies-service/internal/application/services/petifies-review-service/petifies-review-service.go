package reviewservice

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"

	eventModels "github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/events"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/publishers"
	revieweventkafka "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/publishers/review_event/kafka"
	petifiesmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies/mongo"
	petifiesproposalmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies_proposal/mongo"
	petifiessessionmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies_session/mongo"
	reviewmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/review/mongo"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

var logger = logging.New("PetifiesService.ReviewSvc")

type reviewService struct {
	mongoClient          *mongo.Client
	petifiesRepo         petifiesaggre.PetifiesRepository
	petifiesSessionRepo  petifiessessionaggre.PetifiesSessionRepository
	petifiesProposalRepo petifiesproposalaggre.PetifiesProposalRepository
	reviewRepo           reviewaggre.ReviewRepository
	reviewEventPublisher publishers.ReviewEventMessagePublisher
}

type ReviewConfiguration func(rs *reviewService) error

type ReviewService interface {
	CreateReview(ctx context.Context, req *models.CreateReviewReq) (*reviewaggre.ReviewAggre, error)
	EditReview(ctx context.Context, req *models.EditReviewReq) (*reviewaggre.ReviewAggre, error)
	GetById(ctx context.Context, id uuid.UUID) (*reviewaggre.ReviewAggre, error)
	ListByIds(ctx context.Context, ids []uuid.UUID) ([]*reviewaggre.ReviewAggre, error)
	ListByPetifiesId(ctx context.Context, petifiesID uuid.UUID, pageSize int, afterId uuid.UUID) ([]*reviewaggre.ReviewAggre, error)
}

func NewReviewService(cfgs ...ReviewConfiguration) (ReviewService, error) {
	rs := &reviewService{}
	for _, cfg := range cfgs {
		if err := cfg(rs); err != nil {
			return nil, err
		}
	}
	return rs, nil
}

func WithMongoPetifiesRepository(client *mongo.Client) ReviewConfiguration {
	return func(rs *reviewService) error {
		repo := petifiesmongo.New(client)
		rs.petifiesRepo = repo
		return nil
	}
}

func WithMongoPetifiesSessionRepository(client *mongo.Client) ReviewConfiguration {
	return func(rs *reviewService) error {
		repo := petifiessessionmongo.New(client)
		rs.petifiesSessionRepo = repo
		return nil
	}
}

func WithMongoPetifiesProposalRepository(client *mongo.Client) ReviewConfiguration {
	return func(rs *reviewService) error {
		repo := petifiesproposalmongo.New(client)
		rs.petifiesProposalRepo = repo
		return nil
	}
}

func WithMongoReviewRepository(client *mongo.Client) ReviewConfiguration {
	return func(rs *reviewService) error {
		repo := reviewmongo.New(client)
		rs.reviewRepo = repo
		return nil
	}
}

func WithKafkaReviewPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) ReviewConfiguration {
	return func(rs *reviewService) error {
		publisher := revieweventkafka.NewReviewEventPublisher(producer, repo)
		rs.reviewEventPublisher = publisher
		return nil
	}
}

func (rs *reviewService) CreateReview(ctx context.Context, req *models.CreateReviewReq) (*reviewaggre.ReviewAggre, error) {
	logger.Info("Start CreateReview")

	petify, err := rs.petifiesRepo.GetByID(ctx, req.PetifiesID)
	if err != nil {
		logger.ErrorData("Finish CreateReview: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	review, err := reviewaggre.NewReviewAggregate(
		uuid.New(),
		req.PetifiesID,
		req.AuthorID,
		valueobjects.NewImage(req.Image.URI, req.Image.Description),
		req.Review,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		logger.ErrorData("Finish CreateReview: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	savedReview, err := petify.AddReview(*review, rs.reviewRepo)
	if err != nil {
		logger.ErrorData("Finish CreateReview: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = rs.reviewEventPublisher.Publish(ctx, eventModels.ReviewEvent(events.NewReviewCreatedEvent(savedReview)))
	if err != nil {
		logger.ErrorData("Finished CreateReview: Failed to publish event", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreateReview: Successful")
	return savedReview, nil
}

func (rs *reviewService) EditReview(ctx context.Context, req *models.EditReviewReq) (*reviewaggre.ReviewAggre, error) {
	logger.Info("Start EditReview")

	review, err := rs.reviewRepo.GetByID(ctx, req.ID)
	if err != nil {
		logger.ErrorData("Finish EditReview: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	review.SetReview(req.Review)
	review.SetImage(valueobjects.NewImage(req.Image.URI, req.Image.Description))

	updatedReview, err := rs.reviewRepo.Update(ctx, *review)
	if err != nil {
		logger.ErrorData("Finish EditReview: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = rs.reviewEventPublisher.Publish(ctx, eventModels.ReviewEvent(events.NewReviewUpdatedEvent(updatedReview)))
	if err != nil {
		logger.ErrorData("Finished CreateReview: Failed to publish event", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish EditReview: Successful")
	return updatedReview, nil
}

func (rs *reviewService) GetById(ctx context.Context, id uuid.UUID) (*reviewaggre.ReviewAggre, error) {
	logger.Info("Start GetById")

	review, err := rs.reviewRepo.GetByID(ctx, id)
	if err != nil {
		logger.ErrorData("Finish GetById: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetById: Successful")
	return review, nil
}

func (rs *reviewService) ListByIds(ctx context.Context, ids []uuid.UUID) ([]*reviewaggre.ReviewAggre, error) {
	logger.Info("Start ListByIds")

	reviews, err := rs.reviewRepo.ListByIds(ctx, ids)
	if err != nil {
		logger.ErrorData("Finish ListByIds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListByIds: Successful")
	return reviews, nil
}

func (rs *reviewService) ListByPetifiesId(ctx context.Context, petifiesID uuid.UUID, pageSize int, afterId uuid.UUID) ([]*reviewaggre.ReviewAggre, error) {
	logger.Info("Start ListByPetifiesId")

	reviews, err := rs.reviewRepo.GetByPetifiesID(ctx, petifiesID, pageSize, afterId)
	if err != nil {
		logger.ErrorData("Finish ListByPetifiesId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListByPetifiesId: Successful")
	return reviews, nil
}
