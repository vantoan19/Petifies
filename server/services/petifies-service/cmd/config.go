package cmd

import (
	common "github.com/vantoan19/Petifies/server/libs/common-utils"
	config "github.com/vantoan19/Petifies/server/libs/config-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.New("PetifiesService.Cmd")

type Config struct {
	ServerMode string `yaml:"ServerMode,omitempty"`

	WebDomain    string `yaml:"WebDomain,omitempty"`
	GrpcPort     int    `yaml:"GrpcPort,omitempty"`
	MongoUrl     string `yaml:"MongoUrl,omitempty"`
	DatabaseName string `yaml:"DatabaseName,omitempty"`

	Brokers                    []string `yaml:"Brokers,omitempty"`
	PetfiesEventTopic          string   `yaml:"PetfiesEventTopic,omitempty"`
	PetifiesProposalEventTopic string   `yaml:"PetifiesProposalEventTopic,omitempty"`
	PetifiesSessionEventTopic  string   `yaml:"PetifiesSessionEventTopic,omitempty"`
	ReviewEventTopic           string   `yaml:"ReviewEventTopic,omitempty"`
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
			logger.ErrorData("Finished initializeConfig: FAILED", logging.Data{"error": err.Error()})
			return err
		}
	} else {
		logger.Info("Executing initializeConfig: PRODUCTION environment")
	}

	logger.Info("Finished initializeConfig: SUCCESSFUL")
	return nil
}
