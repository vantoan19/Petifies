package cmd

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

func initializeNeo4jDatabase() error {
	logger.Info("Start initializeNeo4jDatabase")

	db, err := dbutils.ConnectToNeo4jDB(Conf.Neo4jUrl, Conf.Neo4jUser, Conf.Neo4jPassword)
	if err != nil {
		logger.ErrorData("Finished initializeNeo4jDatabase: FAILED", logging.Data{"error": err.Error()})
		return err
	}
	DB = db

	logger.Info("Finished initializeNeo4jDatabase: SUCCESSFUL")
	return nil
}

func runMigrations() error {
	logger.Info("Start runMigrations")

	session := DB.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	// Run your migration queries
	_, err := session.Run("CREATE CONSTRAINT IF NOT EXISTS FOR (u:User) REQUIRE u.id IS UNIQUE", nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = session.Run("CREATE CONSTRAINT IF NOT EXISTS FOR (u:User) REQUIRE u.email IS UNIQUE", nil)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Finished runMigrations: SUCCESSFUL")
	return nil
}
