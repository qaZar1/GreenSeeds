package config

import (
	"os"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"gopkg.in/yaml.v2"
)

func MakeConfig(path string) models.Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return models.Config{}
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return models.Config{}
	}

	return cfg
}
