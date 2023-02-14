package cmd

import (
	"errors"
	"time"

	common "github.com/vantoan19/Petifies/server/libs/common-utils"
	config "github.com/vantoan19/Petifies/server/libs/config-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.New("UserService.Config")

type Config struct {
	ServerMode string `yaml:"ServerMode,omitempty"`

	WebDomain   string `yaml:"WebDomain,omitempty"`
	GrpcPort    int    `yaml:"GrpcPort,omitempty"`
	PostgresUrl string `yaml:"PostgresUrl,omitempty"`

	TokenSecretKey      string        `yaml:"TokenSecretKey,omitempty"`
	AccessTokenDuration time.Duration `yaml:"AccessTokenDuration,omitempty"`
}

var (
	Conf     Config
	yamlPath = "/app/config/config.yaml"
)

func initializeConfig() error {
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
