package cmd

import (
	petifiesclient "github.com/vantoan19/Petifies/server/services/grpc-clients/petifies-client"
	postclient "github.com/vantoan19/Petifies/server/services/grpc-clients/post-client"
	relationshipclient "github.com/vantoan19/Petifies/server/services/grpc-clients/relationship-client"
	userclient "github.com/vantoan19/Petifies/server/services/grpc-clients/user-client"
	feedservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/feed"
	petifiesservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/petifies"
	postservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/post"
	relationshipservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/relationship"
	userservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/user"
)

var PostService postservice.PostService
var RelationshipService relationshipservice.RelationshipService
var UserService userservice.UserService
var FeedService feedservice.FeedService
var PetifiesService petifiesservice.PetifiesService

func initUserService() error {
	logger.Info("Start initUserService")

	service, err := userservice.NewUserService(
		UserServiceConn,
		userservice.WithRedisUserCacheRepository(RedisClient, userclient.New(UserServiceConn)),
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
		relationshipservice.WithRedisRelationshipCacheRepository(RedisClient, relationshipclient.New(UserServiceConn)),
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
		postservice.WithRedisPostCacheRepository(RedisClient, postclient.New(PostServiceConn)),
		postservice.WithRedisCommentCacheRepository(RedisClient, postclient.New(PostServiceConn)),
		postservice.WithRedisLoveCacheRepository(RedisClient, postclient.New(PostServiceConn)),
	)
	if err != nil {
		return err
	}

	logger.Info("Finished initPostService: SUCCESSFUL")
	PostService = service
	return nil
}

func initFeedService() error {
	logger.Info("Start initPostService")

	service, err := feedservice.NewFeedService(
		NewfeedServiceConn,
		PostService,
		feedservice.WithRedisFeedCacheRepository(RedisClient),
	)
	if err != nil {
		return err
	}

	logger.Info("Finished initPostService: SUCCESSFUL")
	FeedService = service
	return nil
}

func initPetifiesSerivce() error {
	logger.Info("Start initPetifiesSerivce")

	service, err := petifiesservice.NewPetifiesService(
		PetifiesServiceConn,
		UserServiceConn,
		LocationServiceConn,
		UserService,
		petifiesservice.WithRedisPetifiesCacheRepository(RedisClient, petifiesclient.New(PetifiesServiceConn)),
		petifiesservice.WithRedisPetifiesProposalCacheRepository(RedisClient, petifiesclient.New(PetifiesServiceConn)),
		petifiesservice.WithRedisPetifiesSessionCacheRepository(RedisClient, petifiesclient.New(PetifiesServiceConn)),
		petifiesservice.WithRedisReviewCacheRepository(RedisClient, petifiesclient.New(PetifiesServiceConn)),
	)
	if err != nil {
		return err
	}

	logger.Info("Finished initPetifiesSerivce: SUCCESSFUL")
	PetifiesService = service
	return nil
}
