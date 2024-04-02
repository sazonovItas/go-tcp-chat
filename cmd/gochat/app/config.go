package app

import (
	"fmt"
	"io"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/config"
	"github.com/sazonovItas/gochat-tcp/internal/utils"
)

type Options struct {
	Env string

	LogWriter io.Writer
}

type Config struct {
	TCPServer    config.TCPServer
	CacheStorage config.Redis
	Storage      config.Storage

	Options *Options
}

func InitAppConfig(opts *Options) (*Config, error) {
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

	return &Config{
		TCPServer:    *serverCfg,
		Storage:      *storageCfg,
		CacheStorage: *redisCfg,
		Options:      opts,
	}, nil
}
