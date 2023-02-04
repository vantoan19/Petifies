package main

import "github.com/vantoan19/Petifies/server/services/user-services/auth-service/internal/config"

var initFuncs = []func() error{
	config.InitializeConfig,
}

func initialize() {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			panic(err.Error())
		}
	}
}
