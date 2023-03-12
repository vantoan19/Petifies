package cmd

import (
	common "github.com/vantoan19/Petifies/server/libs/common-utils"
	config "github.com/vantoan19/Petifies/server/libs/config-utils"
)

type Config struct {
	ServerMode string `yaml:"ServerMode,omitempty"`

	WebDomain         string `yaml:"WebDomain,omitempty"`
	GrpcPort          int    `yaml:"GrpcPort,omitempty"`
	CassandraUrl      string `yaml:"CassandraUrl,omitempty"`
	CassandraUser     string `yaml:"CassandraUser,omitempty"`
	CassandraPassword string `yaml:"CassandraPassword,omitempty"`
	Keyspace          string `yaml:"Keyspace,omitempty"`

	UserEventTopic         string   `yaml:"UserEventTopic,omitempty"`
	UserEventConsumerGroup string   `yaml:"UserEventConsumerGroup,omitempty"`
	PostEventTopic         string   `yaml:"PostEventTopic,omitempty"`
	PostEventConsumerGroup string   `yaml:"PostEventConsumerGroup,omitempty"`
	Brokers                []string `yaml:"Brokers,omitempty"`

	RelationshipServiceHost string `yaml:"RelationshipServiceHost"`
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
