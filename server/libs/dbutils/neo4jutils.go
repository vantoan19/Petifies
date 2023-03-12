package dbutils

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var neo4jLogger = logging.New("Libs.DBUtils.Neo4j")

func ConnectToNeo4jDB(dbUrl, dbUsername, dbPassword string) (neo4j.Driver, error) {
	attempt := 0
	neo4jLogger.InfoData("Start ConnectToNeo4jDB", logging.Data{"dbUrl": dbUrl})

	for {
		db, err := openNeo4jDB(dbUrl, dbUsername, dbPassword)
		if err != nil {
			attempt++
			neo4jLogger.WarningData(
				"Executing ConnectToNeo4jDB: Connect to the database fails, attempt again...",
				logging.Data{"attempt": attempt, "error": err.Error()},
			)
		} else {
			neo4jLogger.InfoData("Finished ConnectToNeo4jDB: SUCCESSFUL", logging.Data{"dbUrl": dbUrl})
			return db, nil
		}

		if attempt > 10 {
			neo4jLogger.ErrorData("Finished ConnectToNeo4jDB: FAILED", logging.Data{"dbUrl": dbUrl, "error": err.Error()})
			return nil, err
		}

		neo4jLogger.Info("Executing ConnectToNeo4jDB: Wait for 2 seconds before retrying to connect to the database")
		time.Sleep(5 * time.Second)
		continue
	}
}

func openNeo4jDB(dbUrl, dbUsername, dbPassword string) (neo4j.Driver, error) {
	//nolint
	driver, err := neo4j.NewDriver(dbUrl, neo4j.BasicAuth(dbUsername, dbPassword, ""))
	if err != nil {
		return nil, err
	}

	err = driver.VerifyConnectivity()
	if err != nil {
		return nil, err
	}

	return driver, nil
}
