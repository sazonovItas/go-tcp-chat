package config

import "time"

type TCPServer struct {
	Addr    string        `yaml:"addr"    env:"SERVER_ADDR"`
	Timeout time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT"`
}
