package cmd

import (
	"database/sql"
)

var (
	DB *sql.DB
)

var initFuncs = []func() error{
	initializeConfig,
	initializePostgresDatabase,
	runMigrations,
	initUserProducer,
}

func Initialize() {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			panic(err.Error())
		}
	}
}
