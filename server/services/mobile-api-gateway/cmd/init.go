package cmd

var initFuncs = []func() error{
	initializeConfig,
	initUserServiceClient,
}

func Initialize() {
	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			panic(err.Error())
		}
	}
}
