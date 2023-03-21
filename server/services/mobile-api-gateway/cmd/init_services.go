package cmd

import (
	postservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/post"
	relationshipservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/relationship"
	userservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/user"
)

var PostService postservice.PostService
var RelationshipService relationshipservice.RelationshipService
var UserService userservice.UserService

func initUserService() error {
	logger.Info("Start initUserService")

	service, err := userservice.NewUserService(
		UserServiceConn,
		userservice.WithRedisUserCacheRepository(RedisClient),
	)
	if err != nil {
		return err
	}

	logger.Info("Finished initUserService: SUCCESSFUL")
	UserService = service
	return nil
}

func initRelationshipService() error {
	logger.Info("Start initRelationshipService")

	service, err := relationshipservice.NewRelationshipService(
		RelationshipServiceConn,
		UserService,
		relationshipservice.WithRedisRelationshipCacheRepository(RedisClient),
	)
	if err != nil {
		return err
	}

	logger.Info("Finished initRelationshipService: SUCCESSFUL")
	RelationshipService = service
	return nil
}

func initPostService() error {
	logger.Info("Start initPostService")

	service, err := postservice.NewPostService(
		PostServiceConn,
		UserServiceConn,
		UserService,
		postservice.WithRedisPostCacheRepository(RedisClient),
		postservice.WithRedisCommentCacheRepository(RedisClient),
	)
	if err != nil {
		return err
	}

	logger.Info("Finished initPostService: SUCCESSFUL")
	PostService = service
	return nil
}
