package config

import (
	"errors"

	common "github.com/vantoan19/Petifies/server/libs/common-utils"
	config "github.com/vantoan19/Petifies/server/libs/config-utils"
	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.New("AuthService.Config")

type Config struct {
	ServerMode string `yaml:"ServerMode,omitempty"`

	WebDomain string `yaml:"WebDomain,omitempty"`
	GrpcPort  int    `yaml:"GrpcPort,omitempty"`
}

var (
	Conf     Config
	yamlPath = "/app/config/config.yaml"
)

func InitializeConfig() error {
	logger.Info("Initializing config for mobile api gateway")

	if common.IsDevEnv() {
		err := config.LoadFromYaml(yamlPath, &Conf)
		if err != nil {
			logger.Error("Failed to initialize the config")
			return errors.New("Unable to read configuration")
		}
	} else {
		// Load from k8s for production
	}

	return nil
}
