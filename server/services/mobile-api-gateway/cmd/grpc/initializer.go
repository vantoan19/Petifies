package main

import "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/config"

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
