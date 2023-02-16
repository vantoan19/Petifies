package dbutils

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.New("Libs.DBUtils.Postgres")

func ConnectToDB(dbUrl string) (*sql.DB, error) {
	attempt := 0
	logger.InfoData("Start ConnectToDB", logging.Data{"dbUrl": dbUrl})

	for {
		conn, err := openDB(dbUrl)
		if err != nil {
			attempt++
			logger.WarningData("Executing ConnectToDB: Connect to the database fails, attempt again...", logging.Data{"attempt": attempt})
		} else {
			logger.InfoData("Finished ConnectToDB: SUCCESSFUL", logging.Data{"dbUrl": dbUrl})
			return conn, nil
		}

		if attempt > 10 {
			logger.ErrorData("Finished ConnectToDB: FAILED", logging.Data{"dbUrl": dbUrl, "error": err.Error()})
			return nil, err
		}

		logger.Info("Executing ConnectToDB: Wait for 2 seconds before retrying to connect to the database")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dbUrl string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
