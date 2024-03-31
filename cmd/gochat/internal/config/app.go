package config

import (
	"fmt"

	"github.com/sazonovItas/gochat-tcp/internal/utils"
)

type AppConfig struct {
	Server  TCPServer
	Redis   Redis
	Storage Storage
}

func InitAppConfig() (*AppConfig, error) {
	const op = "gochat.internal.config.app.InitAppConfig"

	serverCfg, err := utils.LoadCfgFromEnv[TCPServer]()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storageCfg, err := utils.LoadCfgFromEnv[Storage]()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	redisCfg, err := utils.LoadCfgFromEnv[Redis]()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &AppConfig{
		Server:  *serverCfg,
		Storage: *storageCfg,
		Redis:   *redisCfg,
	}, nil
}
