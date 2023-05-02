package cmd

import "go.mongodb.org/mongo-driver/mongo"

var (
	MongoClient *mongo.Client
)

var initFuncs = []func() error{
	initializeConfig,
	initializeMongoDatabase,
	runMigrations,
}

func Initialize() {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			panic(err.Error())
		}
	}
}
