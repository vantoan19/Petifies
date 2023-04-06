package reviewmongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/petifies-service/cmd"
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/mappers"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/models"
)

var logger = logging.New("PetifiesService.MongoReviewRepository")

var (
	ErrReviewNotExist = status.Errorf(codes.NotFound, "review does not exist")
)

type reviewMongoRepository struct {
	client           *mongo.Client
	reviewCollection *mongo.Collection
}

func New(client *mongo.Client) reviewaggre.ReviewRepository {
	return &reviewMongoRepository{
		client:           client,
		reviewCollection: client.Database(cmd.Conf.DatabaseName).Collection("reviews"),
	}
}

func (rr *reviewMongoRepository) GetByPetifiesID(ctx context.Context, petifiesID uuid.UUID, pageSize int, afterId uuid.UUID) ([]*reviewaggre.ReviewAggre, error) {
	var reviews []*reviewaggre.ReviewAggre

	err := dbutils.ExecWithSession(ctx, rr.client, func(ssCtx mongo.SessionContext) error {
		reviews_, err := rr.GetByPetifiesIDWithSession(ssCtx, petifiesID, pageSize, afterId)
		if err != nil {
			return err
		}
		reviews = reviews_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (rr *reviewMongoRepository) GetByID(ctx context.Context, id uuid.UUID) (*reviewaggre.ReviewAggre, error) {
	var review *reviewaggre.ReviewAggre

	err := dbutils.ExecWithSession(ctx, rr.client, func(ssCtx mongo.SessionContext) error {
		review_, err := rr.GetByIDWithSession(ssCtx, id)
		if err != nil {
			return err
		}
		review = review_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (rr *reviewMongoRepository) ListByIds(ctx context.Context, ids []uuid.UUID) ([]*reviewaggre.ReviewAggre, error) {
	var reviews []*reviewaggre.ReviewAggre

	err := dbutils.ExecWithSession(ctx, rr.client, func(ssCtx mongo.SessionContext) error {
		reviews_, err := rr.ListByIdsWithSession(ssCtx, ids)
		if err != nil {
			return err
		}
		reviews = reviews_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (rr *reviewMongoRepository) Save(ctx context.Context, review reviewaggre.ReviewAggre) (*reviewaggre.ReviewAggre, error) {
	var savedReview *reviewaggre.ReviewAggre

	err := dbutils.ExecWithSession(ctx, rr.client, func(ssCtx mongo.SessionContext) error {
		review_, err := rr.SaveWithSession(ssCtx, review)
		if err != nil {
			return err
		}
		savedReview = review_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return savedReview, nil
}
func (rr *reviewMongoRepository) Update(ctx context.Context, review reviewaggre.ReviewAggre) (*reviewaggre.ReviewAggre, error) {
	var updatedReview *reviewaggre.ReviewAggre

	err := dbutils.ExecWithSession(ctx, rr.client, func(ssCtx mongo.SessionContext) error {
		review_, err := rr.UpdateWithSession(ssCtx, review)
		if err != nil {
			return err
		}
		updatedReview = review_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedReview, nil
}

func (rr *reviewMongoRepository) GetByPetifiesIDWithSession(ctx context.Context, petifiesID uuid.UUID, pageSize int, afterId uuid.UUID) ([]*reviewaggre.ReviewAggre, error) {
	var createdAt time.Time

	if afterId == uuid.Nil {
		createdAt = time.Now()
	} else {
		petifies, err := rr.GetByIDWithSession(ctx, afterId)
		if err != nil {
			return nil, err
		}
		createdAt = petifies.GetCreatedAt()
	}

	var result []*reviewaggre.ReviewAggre
	var reviews []models.Review
	filter := bson.D{{Key: "petifies_id", Value: petifiesID}, {Key: "created_at", Value: bson.D{{Key: "$lt", Value: createdAt}}}}
	opts := options.Find().SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := rr.reviewCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for _, r := range reviews {
		review, err := mappers.DbModelToReviewAggregate(&r)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		result = append(result, review)
	}

	return result, nil
}

func (rr *reviewMongoRepository) GetByIDWithSession(ctx context.Context, id uuid.UUID) (*reviewaggre.ReviewAggre, error) {
	var reviewDB models.Review
	err := rr.reviewCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&reviewDB)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrReviewNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mappers.DbModelToReviewAggregate(&reviewDB)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rr *reviewMongoRepository) ListByIdsWithSession(ctx context.Context, ids []uuid.UUID) ([]*reviewaggre.ReviewAggre, error) {
	var results []*reviewaggre.ReviewAggre

	cursor, err := rr.reviewCollection.Find(ctx, bson.D{{Key: "id", Value: bson.D{{Key: "$in", Value: ids}}}})
	if err != nil {
		return nil, err
	}
	var reviewDB []models.Review
	if err := cursor.All(ctx, reviewDB); err != nil {
		return nil, err
	}

	for _, r := range reviewDB {
		review, err := mappers.DbModelToReviewAggregate(&r)
		if err != nil {
			return nil, err
		}
		results = append(results, review)
	}

	return results, nil
}

func (rr *reviewMongoRepository) SaveWithSession(ctx context.Context, review reviewaggre.ReviewAggre) (*reviewaggre.ReviewAggre, error) {
	reviewDB := mappers.AggregateReviewToDbReview(&review)
	_, err := rr.reviewCollection.InsertOne(ctx, reviewDB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	savedReview, err := rr.GetByIDWithSession(ctx, review.GetID())
	if err != nil {
		return nil, err
	}
	return savedReview, nil
}

func (rr *reviewMongoRepository) UpdateWithSession(ctx context.Context, review reviewaggre.ReviewAggre) (*reviewaggre.ReviewAggre, error) {
	reviewDB := mappers.AggregateReviewToDbReview(&review)
	_, err := rr.reviewCollection.ReplaceOne(ctx, bson.D{{Key: "id", Value: review.GetID()}}, reviewDB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	savedReview, err := rr.GetByIDWithSession(ctx, review.GetID())
	if err != nil {
		return nil, err
	}
	return savedReview, nil
}
