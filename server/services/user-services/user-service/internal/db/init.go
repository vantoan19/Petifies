package db

import (
	"database/sql"

	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/config"
)

var DB *sql.DB

func InitializePostgresDatabase() error {
	db, err := dbutils.ConnectToDB(config.Conf.PostgresUrl)
	if err != nil {
		return err
	}

	DB = db
	return nil
}
