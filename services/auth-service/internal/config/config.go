package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port      string `env:"AUTH_SERVICE_PORT" env-default:":8001"`
	JWTSecret string `env:"JWT_SECRET" env-default:"dev-only"`
}

func Read() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil{
		return nil, err
	}

	return &cfg, nil
}