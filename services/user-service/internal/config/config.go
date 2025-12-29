package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DBHost     string `env:"POSTGRES_HOST" env-default:"postgres"`
	DBPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	DBUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	DBPassword string `env:"POSTGRES_PASSWORD" env-default:"dev_only"`
	DBName     string `env:"POSTGRES_DB" env-default:"minecraft_server_creator_db"`
	GRPCPort   string `env:"USER_SERVICE_GRPC_SERVER_PORT" env-default:"50052"`
}

func Read() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
