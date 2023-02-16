package configutils

import (
	"os"

	"github.com/go-yaml/yaml"

	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.New("Libs.ConfigUtils")

func LoadFromYaml(configPath string, config interface{}) error {
	logger.Info("Start LoadFromYaml")

	b, err := os.ReadFile(configPath)
	if err != nil {
		logger.ErrorData("Finished LoadFromYaml: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	err = yaml.Unmarshal(b, config)
	if err != nil {
		logger.ErrorData("Finished LoadFromYaml: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finished LoadFromYaml: SUCCESSFUL")
	return nil
}
