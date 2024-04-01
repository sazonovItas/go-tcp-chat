package config

type Redis struct {
	Addr     string `yaml:"address"  env:"REDIS_ADDR"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db"       env:"REDIS_DB"`
}
