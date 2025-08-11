package config

import (
	"os"

	yaml "github.com/goccy/go-yaml"
)

func ReadConfig(fileName string) (Config, error) {
	configFile, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}

	var conf v1
	if err := yaml.Unmarshal(configFile, &conf); err != nil {
		return Config{}, err
	}

	return conf.transform(), nil
}
