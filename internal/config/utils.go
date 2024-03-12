package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Loading config from the environment variable
func LoadCfgFromEnv(envVar string) (*Config, error) {
	const op = "internal.config.utils"

	// Get path to config file from env variable
	configPath := os.Getenv(envVar)
	if configPath == "" {
		return nil, fmt.Errorf("%s: error to get value of the env variable", op)
	}

	cfg, err := LoadCfgFromFile(configPath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Loading config from the file
func LoadCfgFromFile(configPath string) (*Config, error) {
	const op = "internal.config.utils"

	// check presence of the file
	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("%s: error to open config file: %s", op, err.Error())
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: error to read config from file: %s", op, err.Error())
	}

	return &cfg, nil
}
