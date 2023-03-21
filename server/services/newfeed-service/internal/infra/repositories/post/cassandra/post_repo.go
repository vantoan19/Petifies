package cassandra

import (
	"context"
	"errors"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	postfeedaggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/post-feed"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/db/mapper"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/db/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = logging.New("NewFeedService.PostFeedRepository")

type PostRepository struct {
	session *gocql.Session
}

func NewCassandraPostRepository(session *gocql.Session) (*PostRepository, error) {
	return &PostRepository{
		session: session,
	}, nil
}

func (pr *PostRepository) GetByUserID(ctx context.Context, userID uuid.UUID, pageSize int, beforeTime time.Time) ([]*postfeedaggre.PostFeedAggre, error) {
	logger.Info("Start GetByUserID")

	var result []*postfeedaggre.PostFeedAggre
	query := `SELECT user_id, author_id, post_id, created_at 
              FROM postfeeds_by_user_id 
              WHERE user_id=? AND created_at<?
              LIMIT ?`
	iter := pr.session.Query(query, userID.String(), beforeTime, pageSize).Iter()
	var postFeedModel models.PostFeed
	for iter.Scan(&postFeedModel.UserID, &postFeedModel.AuthorID, &postFeedModel.PostID, &postFeedModel.CreatedAt) {
		postfeed, err := mapper.DbPostFeedToPostFeedAggregate(&postFeedModel)
		if err != nil {
			logger.ErrorData("Finish GetByUserID: FAILED", logging.Data{"error": err.Error()})
			return []*postfeedaggre.PostFeedAggre{}, err
		}
		result = append(result, postfeed)
	}
	if err := iter.Close(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	logger.Info("Finish GetByUserID: SUCCESSFUL")
	return result, nil
}

func (pr *PostRepository) ExistsPostFeed(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (bool, error) {
	logger.Info("Start ExistsPostFeed")

	var count int
	query := `SELECT COUNT(*) FROM postfeeds_by_user_id_and_post_id WHERE user_id=? AND post_id=?`
	if err := pr.session.Query(query, userID.String(), postID.String()).Scan(&count); err != nil {
		logger.ErrorData("Finish ExistsPostFeed: FAILED", logging.Data{"error": err.Error()})
		return false, status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish ExistsPostFeed: SUCCESSFUL")
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (pr *PostRepository) Save(ctx context.Context, post postfeedaggre.PostFeedAggre) (*postfeedaggre.PostFeedAggre, error) {
	logger.Info("Start Save")

	postFeedModel, err := mapper.PostFeedAggregateToPostFeedDb(&post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	query := `INSERT INTO postfeeds_by_user_id(user_id, author_id, post_id, created_at) VALUES(?,?,?,?)`
	if err := pr.session.Query(query, postFeedModel.UserID, postFeedModel.AuthorID, postFeedModel.PostID, postFeedModel.CreatedAt).Exec(); err != nil {
		logger.ErrorData("Finish Save: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	query = `INSERT INTO postfeeds_by_user_id_and_post_id(user_id, author_id, post_id, created_at) VALUES(?,?,?,?)`
	if err := pr.session.Query(query, postFeedModel.UserID, postFeedModel.AuthorID, postFeedModel.PostID, postFeedModel.CreatedAt).Exec(); err != nil {
		logger.ErrorData("Finish Save: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish Save: SUCCESSFUL")
	return mapper.DbPostFeedToPostFeedAggregate(postFeedModel)
}

func (pr *PostRepository) Update(ctx context.Context, post postfeedaggre.PostFeedAggre) (*postfeedaggre.PostFeedAggre, error) {
	return nil, errors.New("Not implemented")
}

func (pr *PostRepository) DeleteByUserAndPostID(ctx context.Context, userID uuid.UUID, postID uuid.UUID) (*postfeedaggre.PostFeedAggre, error) {
	return nil, errors.New("Not implemented")
}
