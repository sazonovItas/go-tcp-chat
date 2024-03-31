package config

type Storage struct {
	Name     string `yaml:"name"     env:"DB_NAME"`
	Host     string `yaml:"host"     env:"DB_HOST"`
	Port     string `yaml:"port"     env:"DB_PORT"`
	User     string `yaml:"user"     env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
}
