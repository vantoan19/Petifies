package main

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/config"
)

var (
	DB *sql.DB
)

var initFuncs = []func() error{
	config.InitializeConfig,
	initializePostgresDatabase,
	runMigrations,
}

func initialize() {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			panic(err.Error())
		}
	}
}

func initializePostgresDatabase() error {
	logger.Info("Init DB")

	db, err := dbutils.ConnectToDB(config.Conf.PostgresUrl)
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
