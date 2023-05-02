package petifiesservice

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"

	eventModels "github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/events"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/publishers"
	petifieseventkafka "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/publishers/petifies_event/kafka"
	petifiesmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies/mongo"
	petifiesproposalmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies_proposal/mongo"
	petifiessessionmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/petifies_session/mongo"
	reviewmongo "github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/repositories/review/mongo"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

var logger = logging.New("PetifiesService.PetifiesSvc")

type petifiesService struct {
	mongoClient            *mongo.Client
	petifiesRepo           petifiesaggre.PetifiesRepository
	petifiesSessionRepo    petifiessessionaggre.PetifiesSessionRepository
	petifiesProposalRepo   petifiesproposalaggre.PetifiesProposalRepository
	reviewRepo             reviewaggre.ReviewRepository
	petifiesEventPublisher publishers.PetifiesEventMessagePublisher
}

type PetifiesConfiguration func(ps *petifiesService) error

type PetifesService interface {
	CreatePetify(ctx context.Context, req *models.CreatePetifiesReq) (*petifiesaggre.PetifiesAggre, error)
	EditPetify(ctx context.Context, req *models.EditPetifiesReq) (*petifiesaggre.PetifiesAggre, error)
	GetById(ctx context.Context, id uuid.UUID) (*petifiesaggre.PetifiesAggre, error)
	ListByIds(ctx context.Context, ids []uuid.UUID) ([]*petifiesaggre.PetifiesAggre, error)
	ListByOwnerId(ctx context.Context, ownerID uuid.UUID, pageSize int, afterId uuid.UUID) ([]*petifiesaggre.PetifiesAggre, error)
}

func NewPetifiesService(cfgs ...PetifiesConfiguration) (PetifesService, error) {
	ps := &petifiesService{}
	for _, cfg := range cfgs {
		if err := cfg(ps); err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func WithMongoPetifiesRepository(client *mongo.Client) PetifiesConfiguration {
	return func(ps *petifiesService) error {
		repo := petifiesmongo.New(client)
		ps.petifiesRepo = repo
		return nil
	}
}

func WithMongoPetifiesSessionRepository(client *mongo.Client) PetifiesConfiguration {
	return func(ps *petifiesService) error {
		repo := petifiessessionmongo.New(client)
		ps.petifiesSessionRepo = repo
		return nil
	}
}

func WithMongoPetifiesProposalRepository(client *mongo.Client) PetifiesConfiguration {
	return func(ps *petifiesService) error {
		repo := petifiesproposalmongo.New(client)
		ps.petifiesProposalRepo = repo
		return nil
	}
}

func WithMongoReviewRepository(client *mongo.Client) PetifiesConfiguration {
	return func(ps *petifiesService) error {
		repo := reviewmongo.New(client)
		ps.reviewRepo = repo
		return nil
	}
}

func WithKafkaPetifiesEventPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) PetifiesConfiguration {
	return func(ps *petifiesService) error {
		publisher := petifieseventkafka.NewPetifiesEventPublisher(producer, repo)
		ps.petifiesEventPublisher = publisher
		return nil
	}
}

func (ps *petifiesService) CreatePetify(ctx context.Context, req *models.CreatePetifiesReq) (*petifiesaggre.PetifiesAggre, error) {
	logger.Info("Start CreatePetify")

	petify, err := petifiesaggre.NewPetifiesAggregate(
		uuid.New(),
		req.OwnerID,
		valueobjects.PetifiesType(req.Type),
		req.Title,
		req.Description,
		req.PetName,
		commonutils.Map2(req.Images, func(i models.Image) valueobjects.Image { return valueobjects.NewImage(i.URI, i.Description) }),
		valueobjects.PetifiesUnavailable,
		time.Now(),
		time.Now(),
		entities.Address{
			AddressLineOne: req.Address.AddressLineOne,
			AddressLineTwo: req.Address.AddressLineTwo,
			Street:         req.Address.Street,
			District:       req.Address.District,
			City:           req.Address.City,
			Region:         req.Address.Region,
			PostalCode:     req.Address.PostalCode,
			Country:        req.Address.Country,
			Coordinates:    valueobjects.NewCoordinates(req.Address.Longitude, req.Address.Latitude),
		},
	)
	if err != nil {
		logger.ErrorData("Finish CreatePetify: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	savedPetify, err := ps.petifiesRepo.Save(ctx, *petify)
	if err != nil {
		logger.ErrorData("Finish CreatePetify: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = ps.petifiesEventPublisher.Publish(ctx, eventModels.PetifiesEvent(events.NewPetifiesCreatedEvent(savedPetify)))
	if err != nil {
		logger.ErrorData("Finished CreatePetify: Failed to publish event", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreatePetify: Successful")
	return savedPetify, nil
}

func (ps *petifiesService) EditPetify(ctx context.Context, req *models.EditPetifiesReq) (*petifiesaggre.PetifiesAggre, error) {
	logger.Info("Start EditPetify")

	petify, err := ps.petifiesRepo.GetByID(ctx, req.ID)
	if err != nil {
		logger.ErrorData("Finish EditPetify: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	petify.SetTitle(req.Title)
	petify.SetDescription(req.Description)
	petify.SetPetName(req.PetName)
	petify.SetImages(commonutils.Map2(req.Images, func(i models.Image) valueobjects.Image { return valueobjects.NewImage(i.URI, i.Description) }))
	petify.SetAddress(entities.Address{
		AddressLineOne: req.Address.AddressLineOne,
		AddressLineTwo: req.Address.AddressLineTwo,
		Street:         req.Address.Street,
		District:       req.Address.District,
		City:           req.Address.City,
		Region:         req.Address.Region,
		PostalCode:     req.Address.PostalCode,
		Country:        req.Address.Country,
		Coordinates:    valueobjects.NewCoordinates(req.Address.Longitude, req.Address.Latitude),
	})

	updatedPetify, err := ps.petifiesRepo.Update(ctx, *petify)
	if err != nil {
		logger.ErrorData("Finish EditPetify: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = ps.petifiesEventPublisher.Publish(ctx, eventModels.PetifiesEvent(events.NewPetifiesUpdatedEvent(updatedPetify)))
	if err != nil {
		logger.ErrorData("Finished EditPetify: Failed to publish event", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish EditPetify: Successful")
	return updatedPetify, nil
}

func (ps *petifiesService) GetById(ctx context.Context, id uuid.UUID) (*petifiesaggre.PetifiesAggre, error) {
	logger.Info("Start GetById")

	petify, err := ps.petifiesRepo.GetByID(ctx, id)
	if err != nil {
		logger.ErrorData("Finish GetById: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetById: Successful")
	return petify, nil
}

func (ps *petifiesService) ListByIds(ctx context.Context, ids []uuid.UUID) ([]*petifiesaggre.PetifiesAggre, error) {
	logger.Info("Start ListByIds")

	petifies, err := ps.petifiesRepo.ListByIDs(ctx, ids)
	if err != nil {
		logger.ErrorData("Finish ListByIds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListByIds: Successful")
	return petifies, nil
}

func (ps *petifiesService) ListByOwnerId(ctx context.Context, ownerID uuid.UUID, pageSize int, afterId uuid.UUID) ([]*petifiesaggre.PetifiesAggre, error) {
	logger.Info("Start ListByOwnerId")

	petifies, err := ps.petifiesRepo.GetByUserID(ctx, ownerID, pageSize, afterId)
	if err != nil {
		logger.ErrorData("Finish ListByOwnerId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListByOwnerId: Successful")
	return petifies, nil
}
