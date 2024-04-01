package config

import "time"

type TCPServer struct {
	Addr        string        `yaml:"addr"         env:"SERVER_ADDR"`
	Timeout     time.Duration `yaml:"timeout"      env:"SERVER_TIMEOUT"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"SERVER_IDLE_TIMEOUT"`
}
