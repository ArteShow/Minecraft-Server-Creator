package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort   string `env:"HOST_METADATA_SERVICE_GRPC_PORT" env-default:"50052"`
	DBHost     string `env:"POSTGRES_HOST" env-default:"postgres-hub"`
	DBPort     string `env:"POSTGRES_PORT" env-default:"2345"`
	DBUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	DBPassword string `env:"POSTGRES_PASSWORD" env-default:"dev_only"`
	DBName     string `env:"POSTGRES_DB" env-default:"minecraft_server_creator_db"`
}

type BundleConfig struct {
	Bundles map[string]Bundle `json:"bundles"`
}

type Bundle struct {
	RAM   string `json:"RAM"`
	Cores string `json:"Cores"`
}

func Read() (*Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetBundles() (*BundleConfig, error) {
	file, err := os.Open("bundles.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg BundleConfig
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
