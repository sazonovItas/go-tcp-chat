package app

import (
	"fmt"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/config"
	"github.com/sazonovItas/gochat-tcp/internal/utils"
)

type AppConfig struct {
	Server  config.TCPServer
	Redis   config.Redis
	Storage config.Storage
}

func InitAppConfig() (*AppConfig, error) {
	const op = "gochat.internal.config.app.InitAppConfig"

	serverCfg, err := utils.LoadCfgFromEnv[config.TCPServer]()
	if err != nil {
		return nil, fmt.Errorf("%s: error load server config: %w", op, err)
	}

	storageCfg, err := utils.LoadCfgFromEnv[config.Storage]()
	if err != nil {
		return nil, fmt.Errorf("%s: error load storage config %w", op, err)
	}

	redisCfg, err := utils.LoadCfgFromEnv[config.Redis]()
	if err != nil {
		return nil, fmt.Errorf("%s: error load redis config %w", op, err)
	}

	return &AppConfig{
		Server:  *serverCfg,
		Storage: *storageCfg,
		Redis:   *redisCfg,
	}, nil
}
