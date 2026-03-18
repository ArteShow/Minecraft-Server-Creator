package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DBHost     string `env:"POSTGRES_HOST" env-default:"postgres"`
	DBPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	DBUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	DBPassword string `env:"POSTGRES_PASSWORD" env-default:"dev_only"`
	DBName     string `env:"POSTGRES_DB" env-default:"minecraft_server_creator_db"`
	GRPCPort   string `env:"NETWORK_SERVICE_GRPC_PORT" env-default:"50050"`
}

func Read() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
