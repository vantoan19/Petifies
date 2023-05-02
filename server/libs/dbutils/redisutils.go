package dbutils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var redisLogger = logging.New("Libs.DBUtils.Redis")

func ConnectToRedisDB(redisUrl, user, password string, db int) (*redis.Client, error) {
	attempt := 0
	redisLogger.InfoData("Start ConnectToRedisDB", logging.Data{"redisUrl": redisUrl})

	for {
		client, err := openRedis(redisUrl, user, password, db)
		if err != nil {
			attempt++
			redisLogger.WarningData("Executing ConnectToRedisDB: Connect to the database fails, attempt again...", logging.Data{"attempt": attempt})
		} else {
			redisLogger.InfoData("Finished ConnectToRedisDB: SUCCESSFUL", logging.Data{"redisUrl": redisUrl})
			return client, nil
		}

		if attempt > 10 {
			redisLogger.ErrorData("Finished ConnectToRedisDB: FAILED", logging.Data{"redisUrl": redisUrl, "error": err.Error()})
			return nil, err
		}

		redisLogger.Info("Executing ConnectToRedisDB: Wait for 2 seconds before retrying to connect to the database")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openRedis(redisURL, user, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Username: user,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
