package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DBHost     string `env:"POSTGRES_HOST" env-default:"postgres-hub"`
	DBPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	DBUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	DBPassword string `env:"POSTGRES_PASSWORD" env-default:"dev_only"`
	DBName     string `env:"POSTGRES_DB" env-default:"minecraft_server_creator_db"`
	Port   string `env:"BUNDLE_SERVICE_PORT" env-default:"8015"`
	GRPCPort   string `env:"BUNDLE_SERVICE_GRPC_PORT" env-default:"50055"`
}

func Read() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
