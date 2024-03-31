package utils

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Loading config from the environment variable
func LoadCfgFromEnv[T any]() (*T, error) {
	const op = "internal.config.configutils.LoadCfgFromEnv"

	var cfg T
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: cannot load config from env %w", op, err)
	}

	return &cfg, nil
}

// Loading config from the file
func LoadCfgFromFile[T any](configPath string) (*T, error) {
	const op = "internal.config.utils.LoadCfgFromFile"

	// check presence of the file
	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("%s: error to open config file: %s", op, err.Error())
	}

	var cfg T
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: error to read config from file: %s", op, err.Error())
	}

	return &cfg, nil
}

// GetEnv returns value of environment variable ENV or local if ENV
func GetEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	return env
}
