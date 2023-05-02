package locationmongo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/location-service/cmd"
	locationaggre "github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location/valueobjects"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/infra/db/mappers"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/infra/db/models"
)

var logger = logging.New("LocationService.MongoLocationRepository")

var (
	ErrLocationNotExist = status.Errorf(codes.NotFound, "location does not exist")
	wc                  = writeconcern.New(writeconcern.WMajority())
	rc                  = readconcern.Snapshot()
	transOpts           = options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
)

type locationRepository struct {
	client             *mongo.Client
	locationCollection *mongo.Collection
}

func New(client *mongo.Client) locationaggre.LocationRepository {
	return &locationRepository{
		client:             client,
		locationCollection: client.Database(cmd.Conf.DatabaseName).Collection("locations"),
	}
}

func (lr *locationRepository) FindNearbyLocationsByEntityType(
	ctx context.Context,
	longitude,
	latitude float64,
	maxDistance float64,
	entityType valueobjects.EntityType,
	pageSize int,
	offset int,
) ([]*locationaggre.LocationAggre, error) {
	logger.Info("Start FindNearbyLocationsByEntityType")
	var locations []*locationaggre.LocationAggre

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		locations_, err := lr.FindNearbyLocationsByEntityTypeWithSession(ssCtx, longitude, latitude, maxDistance, entityType, pageSize, offset)
		if err != nil {
			return err
		}
		locations = locations_

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish FindNearbyLocationsByEntityType: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish FindNearbyLocationsByEntityType: Successful")
	return locations, nil
}

func (lr *locationRepository) GetByID(ctx context.Context, id uuid.UUID) (*locationaggre.LocationAggre, error) {
	logger.Info("Start GetByID")
	var location *locationaggre.LocationAggre

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		location_, err := lr.GetByIDWithSesion(ssCtx, id)
		if err != nil {
			return err
		}
		location = location_

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish GetByID: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetByID: Successful")
	return location, nil
}

func (lr *locationRepository) GetByEntityID(ctx context.Context, entityID uuid.UUID) (*locationaggre.LocationAggre, error) {
	logger.Info("Start GetByEntityID")
	var location *locationaggre.LocationAggre

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		location_, err := lr.GetByEntityIDWithSession(ssCtx, entityID)
		if err != nil {
			return err
		}
		location = location_

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish GetByEntityID: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetByEntityID: Successful")
	return location, nil
}

func (lr *locationRepository) Save(ctx context.Context, location locationaggre.LocationAggre) (*locationaggre.LocationAggre, error) {
	logger.Info("Start Save")
	var savedLocation *locationaggre.LocationAggre

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		location_, err := lr.SaveWithSession(ssCtx, location)
		if err != nil {
			return err
		}
		savedLocation = location_

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish Save: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish Save: Successful")
	return savedLocation, nil
}

func (lr *locationRepository) Update(ctx context.Context, location locationaggre.LocationAggre) (*locationaggre.LocationAggre, error) {
	logger.Info("Start Update")
	var updatedLocation *locationaggre.LocationAggre

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		location_, err := lr.UpdateLocationWithSession(ssCtx, location)
		if err != nil {
			return err
		}
		updatedLocation = location_

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish Update: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish Update: Successful")
	return updatedLocation, nil
}

func (lr *locationRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	logger.Info("Start DeleteByID")

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		err := lr.DeleteByIDWithSession(ssCtx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish DeleteByID: Failed", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finish DeleteByID: Successful")
	return nil
}

func (lr *locationRepository) DeleteByEntityID(ctx context.Context, entityID uuid.UUID) error {
	logger.Info("Start DeleteByEntityID")

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		err := lr.DeleteByEntityIDWithSession(ssCtx, entityID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish DeleteByEntityID: Failed", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finish DeleteByEntityID: Successful")
	return nil
}

func (lr *locationRepository) FindNearbyLocationsByEntityTypeWithSession(
	ctx context.Context,
	longitude,
	latitude float64,
	maxDistance float64,
	entityType valueobjects.EntityType,
	pageSize int,
	offset int,
) ([]*locationaggre.LocationAggre, error) {
	var results []*locationaggre.LocationAggre

	center := bson.D{{Key: "type", Value: "Point"}, {Key: "coordinates", Value: []float64{longitude, latitude}}}
	filter := bson.D{
		{
			Key: "location",
			Value: bson.D{
				{
					Key: "$nearSphere",
					Value: bson.D{
						{
							Key:   "$geometry",
							Value: center,
						},
						{
							Key:   "$maxDistance",
							Value: maxDistance,
						},
					},
				},
			},
		},
		{
			Key:   "entity_type",
			Value: string(entityType),
		},
	}
	opts := options.Find().SetLimit(int64(pageSize)).SetSkip(int64(offset))

	var locations []models.Location
	cursor, err := lr.locationCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for _, l := range locations {
		location, err := mappers.DbModelToLocationAggregate(&l)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		results = append(results, location)
	}

	return results, nil
}

func (lr *locationRepository) GetByIDWithSesion(ctx context.Context, id uuid.UUID) (*locationaggre.LocationAggre, error) {
	var location models.Location
	err := lr.locationCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&location)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLocationNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mappers.DbModelToLocationAggregate(&location)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lr *locationRepository) GetByEntityIDWithSession(ctx context.Context, entityID uuid.UUID) (*locationaggre.LocationAggre, error) {
	var location models.Location
	err := lr.locationCollection.FindOne(ctx, bson.D{{Key: "entity_id", Value: entityID}}).Decode(&location)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLocationNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mappers.DbModelToLocationAggregate(&location)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lr *locationRepository) SaveWithSession(ctx context.Context, location locationaggre.LocationAggre) (*locationaggre.LocationAggre, error) {
	locationModel := mappers.AggregateLocationToDbLocation(&location)
	_, err := lr.locationCollection.InsertOne(ctx, locationModel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	location_, err := lr.GetByIDWithSesion(ctx, location.ID())
	if err != nil {
		return nil, err
	}
	return location_, nil
}

func (lr *locationRepository) UpdateLocationWithSession(ctx context.Context, location locationaggre.LocationAggre) (*locationaggre.LocationAggre, error) {
	locationModel := mappers.AggregateLocationToDbLocation(&location)
	locationModel.UpdatedAt = time.Now()

	_, err := lr.locationCollection.ReplaceOne(ctx, bson.D{{Key: "id", Value: location.ID()}}, locationModel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	location_, err := lr.GetByEntityIDWithSession(ctx, location.ID())
	if err != nil {
		return nil, err
	}
	return location_, nil
}

func (lr *locationRepository) DeleteByIDWithSession(ctx context.Context, id uuid.UUID) error {
	_, err := lr.locationCollection.DeleteOne(ctx, bson.D{{Key: "id", Value: id}})
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

func (lr *locationRepository) DeleteByEntityIDWithSession(ctx context.Context, entityID uuid.UUID) error {
	_, err := lr.locationCollection.DeleteOne(ctx, bson.D{{Key: "entity_id", Value: entityID}})
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

func (lr *locationRepository) execSession(ctx context.Context, fn func(ssCtx mongo.SessionContext) error) error {
	session, err := lr.client.StartSession()
	defer session.EndSession(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	err = session.StartTransaction(transOpts)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	if err = fn(mongo.NewSessionContext(ctx, session)); err != nil {
		if abErr := session.AbortTransaction(ctx); abErr != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("session err: %v, abort err: %v", err, abErr))
		}
		if err == ErrLocationNotExist {
			return err
		}
		return status.Errorf(codes.Internal, err.Error())
	}

	return session.CommitTransaction(ctx)
}
