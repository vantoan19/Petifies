package cmd

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

func initializePostgresDatabase() error {
	logger.Info("Init DB")

	db, err := dbutils.ConnectToDB(Conf.PostgresUrl)
	if err != nil {
		logger.ErrorData("Failed to Init DB", logging.Data{"error": err.Error()})
		return err
	}
	DB = db

	logger.Info("Init DB successfully")
	return nil
}

func runMigrations() error {
	logger.Info("Start running migrations")

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		logger.ErrorData("Failed to get DB driver", logging.Data{"error": err.Error()})
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		logger.ErrorData("Failed to get migrate", logging.Data{"error": err.Error()})
		return err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.ErrorData("Failed to upgrade the db", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Migrated the DB successfully")
	return nil
}
