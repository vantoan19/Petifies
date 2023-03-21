package cmd

import "github.com/vantoan19/Petifies/server/libs/logging-config"

var logger = logging.New("MobileGateWay.Cmd")

var initFuncs = []func() error{
	initializeConfig,
	initializeRedisCache,
	initUserServiceClient,
	initPostServiceClient,
	initMediaServiceClient,
	initRelationshipServiceClient,
	initNewfeedServiceClient,
	initUserService,
	initRelationshipService,
	initPostService,
}

func Initialize() {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			panic(err.Error())
		}
	}
}
