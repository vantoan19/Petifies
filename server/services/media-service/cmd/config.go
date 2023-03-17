package cmd

import (
	common "github.com/vantoan19/Petifies/server/libs/common-utils"
	config "github.com/vantoan19/Petifies/server/libs/config-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.New("MediaService.Cmd")

type Config struct {
	ServerMode string `yaml:"ServerMode,omitempty"`

	WebDomain      string `yaml:"WebDomain,omitempty"`
	CDNDomain      string `yaml:"CDNDomain"`
	GrpcPort       int    `yaml:"GrpcPort,omitempty"`
	StorageRootDir string `yaml:"StorageRootDir,omitempty"`
	MaxFileSize    int    `yaml:"MaxFileSize,omitempty"`
	BucketName     string `yaml:"BucketName"`
	CredentialFile string `yaml:"CredentialFile"`
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
