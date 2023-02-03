package configutils

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"

	logging "github.com/vantoan19/Petifies/server/libs/logging-config"
)

var logger = logging.NewLogger("Libs.ConfigUtils")

func LoadFromYaml(configPath string, config interface{}) error {
	logger.Info("Loading config from yaml file")

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.ErrorData("Fail to load config from yaml file", logging.Data{"error": err.Error()})
		return err
	}

	err = yaml.Unmarshal(b, config)
	if err != nil {
		logger.ErrorData("Fail to load config from yaml file", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Loaded config from yaml file successfully")
	return nil
}
