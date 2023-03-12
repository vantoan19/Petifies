package cmd

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

var (
	DB neo4j.Driver
)

var initFuncs = []func() error{
	initializeConfig,
	initializeNeo4jDatabase,
	runMigrations,
}

func Initialize() {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			panic(err.Error())
		}
	}
}
