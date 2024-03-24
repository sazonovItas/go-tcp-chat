package config

import "time"

type Config struct {
	Env       string     `yaml:"env"        env:"ENV" env-default:"local"`
	TCPServer *TCPServer `yaml:"tcp_server"`
	Storage   *Storage   `yaml:"storage"`
}

type TCPServer struct {
	Addr        string        `yaml:"addr"         env:"ADDR" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout"                 env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout"            env-default:"60s"`
}

type Storage struct {
	Name     string `yaml:"name"     env:"DBNAME"`
	Host     string `yaml:"host"     env:"DBUSER"`
	Port     string `yaml:"port"     env:"DBPORT"`
	User     string `yaml:"user"     env:"DBUSER"`
	Password string `yaml:"password" env:"DBPASSWORD"`
}
