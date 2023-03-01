package dbutils

import (
	"context"
	"time"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoLogger = logging.New("Libs.DBUtils.Mongo")

func ConnectToMongoDB(dbUrl string) (*mongo.Client, error) {
	attempt := 0
	logger.InfoData("Start ConnectToMongoDB", logging.Data{"dbUrl": dbUrl})

	for {
		client, err := openMongoDB(dbUrl)
		if err != nil {
			attempt++
			logger.WarningData("Executing ConnectToMongoDB: Connect to the database fails, attempt again...", logging.Data{"attempt": attempt})
		} else {
			logger.InfoData("Finished ConnectToMongoDB: SUCCESSFUL", logging.Data{"dbUrl": dbUrl})
			return client, nil
		}

		if attempt > 10 {
			logger.ErrorData("Finished ConnectToMongoDB: FAILED", logging.Data{"dbUrl": dbUrl, "error": err.Error()})
			return nil, err
		}

		logger.Info("Executing ConnectToMongoDB: Wait for 2 seconds before retrying to connect to the database")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openMongoDB(dbUrl string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI((options.ServerAPIVersion1))
	opts := options.Client().ApplyURI(dbUrl).SetServerAPIOptions(serverAPI).SetTimeout(time.Second * 2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}
