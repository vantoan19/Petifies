package cmd

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

func initializeMongoDatabase() error {
	logger.Info("Start initializeMongoDatabase")

	client, err := dbutils.ConnectToMongoDB(Conf.MongoUrl)
	if err != nil {
		logger.ErrorData("Finished initializeMongoDatabase: FAILED", logging.Data{"error": err.Error()})
		return err
	}
	MongoClient = client

	logger.Info("Finished initializeMongoDatabase: SUCCESSFUL")
	return nil
}

func runMigrations() error {
	logger.Info("Start runMigrations")

	driver, err := mongodb.WithInstance(MongoClient, &mongodb.Config{
		DatabaseName: Conf.DatabaseName,
	})
	if err != nil {
		logger.ErrorData("Finished runMigrations: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"mongodb", driver)
	if err != nil {
		logger.ErrorData("Finished runMigrations: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.ErrorData("Finished runMigrations: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finished runMigrations: SUCCESSFUL")
	return nil
}
