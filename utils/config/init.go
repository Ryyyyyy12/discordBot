package config

import (
	"io/ioutil"
	"kmuttBot/utils/logger"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// initialize Config struct (empty)
var C = &Config{}

func init() { //special func in Go that runs before main()

	// Read YAML configuration file
	yml, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		logger.Log(logrus.Fatal, "UNABLE TO READ YAML CONFIGURATION FILE")
	}

	// Unmarshal (decode) YAML into Config struct (C)
	err = yaml.Unmarshal(yml, C)
	if err != nil {
		logger.Log(logrus.Fatal, "UNABLE TO PARSE YAML CONFIGURATION FILE")
	}

	// Apply configurations
	logrus.SetLevel(logrus.WarnLevel)

}
