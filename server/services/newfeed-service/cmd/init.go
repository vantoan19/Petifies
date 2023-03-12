package cmd

import "github.com/vantoan19/Petifies/server/libs/logging-config"

var logger = logging.New("NewFeedService.Cmd")

var initFuncs = []func() error{
	initializeConfig,
	initializeCassandraDatabase,
	runMigrations,
	initRelationshipServiceClient,
}

func Initialize() {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			panic(err.Error())
		}
	}
}
