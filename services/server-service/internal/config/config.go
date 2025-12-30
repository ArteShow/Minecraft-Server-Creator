package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	StartRAM string `env:"START_RAM" env-default:":2G"`
	RunRAM   string `env:"RUN_RAM" env-default:"4G"`
}

func Read() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
