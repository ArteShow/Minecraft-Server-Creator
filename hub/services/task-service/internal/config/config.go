package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port                  string `env:"TASK_SERVICE_PORT" env-default:"8013"`
	DefaultHostServerPort string `env:"DEFAULT_HOST_SERVER_PORT" env-default:"8003"`
}

type BundleConfig struct {
	Bundles map[string]Bundle `json:"bundles"`
}

type Bundle struct {
	RAM   int `json:"RAM"`
	Cores int `json:"Cores"`
	Backups int `json:"Backups"`
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