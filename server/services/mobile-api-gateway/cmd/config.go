package cmd

import (
	common "github.com/vantoan19/Petifies/server/libs/common-utils"
	config "github.com/vantoan19/Petifies/server/libs/config-utils"
)

type Config struct {
	ServerMode string `yaml:"ServerMode,omitempty"`

	WebDomain string `yaml:"WebDomain,omitempty"`
	GrpcPort  int    `yaml:"GrpcPort,omitempty"`

	TLSKeyPath  string `yaml:"TLSKeyPath,omitempty"`
	TLSCertPath string `yaml:"TLSCertPath,omitempty"`

	RedisURL      string `yaml:"RedisURL"`
	RedisUser     string `yaml:"RedisUser"`
	RedisPassword string `yaml:"RedisPassword"`
	RedisDatabase int    `yaml:"RedisDatabase"`

	UserServiceHost         string `yaml:"UserServiceHost"`
	PostServiceHost         string `yaml:"PostServiceHost"`
	MediaServiceHost        string `yaml:"MediaServiceHost"`
	RelationshipServiceHost string `yaml:"RelationshipServiceHost"`
	NewfeedServiceHost      string `yaml:"NewfeedServiceHost"`
	PetifiesServiceHost     string `yaml:"PetifiesServiceHost"`
	LocationServiceHost     string `yaml:"LocationServiceHost"`
}

var (
	Conf     Config
	yamlPath = "/app/config/config.yaml"
)

func initializeConfig() error {
	logger.Info("Start initializeConfig")

	if common.IsDevEnv() {
		logger.Info("Executing initializeConfig: DEV environment")
		err := config.LoadFromYaml(yamlPath, &Conf)
		if err != nil {
			logger.Error("Finished initializeConfig: FAILED")
			return err
		}
	} else {
		logger.Info("Executing initializeConfig: PRODUCTION environment")
	}

	logger.Info("Finished initializeConfig: SUCCESSFUL")
	return nil
}
