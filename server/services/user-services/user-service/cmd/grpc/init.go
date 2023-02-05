package main

import (
	"github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/db"
	"github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/config"
) 

var initFuncs = []func() error{
	config.InitializeConfig,
	db.InitializePostgresDatabase,
}

func initialize() {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			panic(err.Error())
		}
	}
}
