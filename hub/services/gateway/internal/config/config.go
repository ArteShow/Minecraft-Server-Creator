package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port       string `env:"GATEWAY_PORT" env-default:":8010"`
	APIVersion string `env:"API_VERSION" env-default:"v1"`
	JWTSecret string `env:"JWT_SECRET" env-default:"dev_only"`
}

func Read() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
