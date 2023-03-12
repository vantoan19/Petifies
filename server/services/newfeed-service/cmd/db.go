package cmd

import (
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var DB *gocql.Session

func initializeCassandraDatabase() error {
	logger.Info("Start initializePostgresDatabase")

	db, err := dbutils.ConnectToCassandraDB(Conf.CassandraUrl, Conf.CassandraUser, Conf.CassandraPassword, Conf.Keyspace)
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

	driver, err := cassandra.WithInstance(DB, &cassandra.Config{KeyspaceName: Conf.Keyspace, MultiStatementEnabled: true, MultiStatementMaxSize: 20})
	if err != nil {
		logger.ErrorData("Finished runMigrations: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"cassandra", driver)
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
