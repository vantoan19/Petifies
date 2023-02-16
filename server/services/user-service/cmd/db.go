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
	logger.Info("Start initializePostgresDatabase")

	db, err := dbutils.ConnectToDB(Conf.PostgresUrl)
	if err != nil {
		logger.ErrorData("Finished initializePostgresDatabase: FAILED", logging.Data{"error": err.Error()})
		return err
	}
	DB = db

	logger.Info("Finished initializePostgresDatabase: SUCCESSFUL")
	return nil
}

func runMigrations() error {
	logger.Info("Start runMigrations")

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		logger.ErrorData("Finished runMigrations: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
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
