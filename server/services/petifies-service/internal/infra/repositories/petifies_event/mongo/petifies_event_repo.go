package petifieseventmongo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/libs/logging-config"

	"github.com/vantoan19/Petifies/server/services/petifies-service/cmd"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/mappers"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/models"
)

var logger = logging.New("PostService.PetifiesEventRepository")

type petifiesEventRepository struct {
	client                  *mongo.Client
	petifiesEventCollection *mongo.Collection
}

func New(client *mongo.Client) (outbox_repo.EventRepository, error) {
	return &petifiesEventRepository{
		client:                  client,
		petifiesEventCollection: client.Database(cmd.Conf.DatabaseName).Collection("petifies_events"),
	}, nil
}

func (pr *petifiesEventRepository) AddEvent(event outbox_repo.Event) (*outbox_repo.Event, error) {
	logger.Info("Start AddEvent")
	dbEvent, err := mappers.OutboxEventToDbPetifiesEvent(&event)
	if err != nil {
		logger.ErrorData("Finish AddEvent: Failed", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	_, err = pr.petifiesEventCollection.InsertOne(context.Background(), dbEvent)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var petifiesEvent models.PetifiesEvent
	err = pr.petifiesEventCollection.FindOne(context.Background(), bson.D{{Key: "id", Value: event.ID}}).Decode(&petifiesEvent)
	if err != nil {
		logger.ErrorData("Finish AddEvent: Failed", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish AddEvent: Success")
	return mappers.DbPetifiesEventToOutboxEvent(&petifiesEvent)
}

func (pr *petifiesEventRepository) GetEventsByLockerID(lockerID uuid.UUID) ([]*outbox_repo.Event, error) {
	logger.Info("Start GetEventsByLockerID")
	var events []models.PetifiesEvent
	fmt.Println(lockerID)
	cursor, err := pr.petifiesEventCollection.Find(context.Background(), bson.D{{Key: "locked_by", Value: lockerID}})
	if err != nil {
		logger.ErrorData("Finish AddEvent: GetEventsByLockerID", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(context.Background(), &events); err != nil {
		logger.ErrorData("Finish GetEventsByLockerID: Failed", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	fmt.Println(events)

	var results []*outbox_repo.Event
	for _, e := range events {
		e_, err := mappers.DbPetifiesEventToOutboxEvent(&e)
		if err != nil {
			logger.ErrorData("Finish GetEventsByLockerID: Failed", logging.Data{"error": err.Error()})
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		results = append(results, e_)
	}

	logger.Info("Finish GetEventsByLockerID: Success")
	return results, err
}

func (pr *petifiesEventRepository) LockStartedEvents(lockerID uuid.UUID) error {
	logger.Info("Start LockStartedEvents")
	filter := bson.D{{Key: "outbox_state", Value: string(models.OutboxStateSTARTED)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "locked_by", Value: lockerID}, {Key: "locked_at", Value: time.Now()}}}}

	fmt.Println(lockerID)

	res, err := pr.petifiesEventCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		logger.ErrorData("Finish LockStartedEvents: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}
	fmt.Println(res)

	logger.Info("Finish LockStartedEvents: Success")
	return nil
}

func (pr *petifiesEventRepository) UpdateEvent(event outbox_repo.Event) error {
	logger.Info("Start UpdateEvent")

	dbEvent, err := mappers.OutboxEventToDbPetifiesEvent(&event)
	if err != nil {
		logger.ErrorData("Finish AddEvent: UpdateEvent", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}

	_, err = pr.petifiesEventCollection.ReplaceOne(context.Background(), bson.D{{Key: "id", Value: event.ID}}, dbEvent)
	if err != nil {
		logger.ErrorData("Finish UpdateEvent: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish UpdateEvent: Success")
	return nil
}

func (pr *petifiesEventRepository) UnlockEventsByLockerID(lockerID uuid.UUID) error {
	logger.Info("Start UnlockEventsByLockerID")
	filter := bson.D{{Key: "locked_by", Value: lockerID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "locked_by", Value: nil}, {Key: "locked_at", Value: nil}}}}

	_, err := pr.petifiesEventCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		logger.ErrorData("Finish UnlockEventsByLockerID: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish UnlockEventsByLockerID: Success")
	return nil
}

func (pr *petifiesEventRepository) UnlockEventsBeforeDatetime(t time.Time) error {
	logger.Info("Start UnlockEventsBeforeDatetime")
	filter := bson.D{{Key: "locked_at", Value: bson.D{{Key: "$lt", Value: t}}}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "locked_by", Value: nil}, {Key: "locked_at", Value: nil}}}}

	_, err := pr.petifiesEventCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		logger.ErrorData("Finish UnlockEventsBeforeDatetime: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish UnlockEventsBeforeDatetime: Success")
	return nil
}

func (pr *petifiesEventRepository) DeleteEventsBeforeDatetime(t time.Time) error {
	logger.Info("Start DeleteEventsBeforeDatetime")
	filter := bson.D{{Key: "created_at", Value: bson.D{{Key: "$lt", Value: t}}}}

	_, err := pr.petifiesEventCollection.DeleteMany(context.Background(), filter)
	if err != nil {
		logger.ErrorData("Finish DeleteEventsBeforeDatetime: Failed", logging.Data{"error": err.Error()})
		return status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish DeleteEventsBeforeDatetime: Success")
	return nil
}
