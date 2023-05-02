package cmd

import (
	"github.com/redis/go-redis/v9"

	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var RedisClient *redis.Client

func initializeRedisCache() error {
	logger.Info("Start initializeRedisCache")

	client, err := dbutils.ConnectToRedisDB(Conf.RedisURL, Conf.RedisUser, Conf.RedisPassword, Conf.RedisDatabase)
	if err != nil {
		logger.ErrorData("Finished initializeRedisCache: FAILED", logging.Data{"error": err.Error()})
		return err
	}
	RedisClient = client

	logger.Info("Finished initializeRedisCache: SUCCESSFUL")
	return nil
}
