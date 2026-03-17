package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port string `env:"API_GATEWAY_PORT" env-default:":8000"`
	APIVersion string `env:"API_VERSION" env-default:"v1"`
	JWTKey string `env:"JWT_SECRET" env-default:"dev-only"`
}

func Read() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil{
		return nil, err
	}

	return &cfg, nil
}