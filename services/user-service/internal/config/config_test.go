package config

import "testing"

func TestRead(t *testing.T) {
	cfg, err := Read()
	if err != nil {
		t.Fatal("Failed to read the config, ", err.Error())
	}

	if cfg.DBHost != "postgres" {
		t.Fatalf("expected port :8000, got: %s", cfg.DBHost)
	}

	if cfg.DBName != "minecraft_server_creator_db" {
		t.Fatalf("expected version v1, got: %s", cfg.DBName)
	}

	if cfg.DBPassword != "dev_only" {
		t.Fatalf("expected version v1, got: %s", cfg.DBPassword)
	}

	if cfg.DBPort != "5432" {
		t.Fatalf("expected version v1, got: %s", cfg.DBPort)
	}

	if cfg.DBUser != "postgres" {
		t.Fatalf("expected version v1, got: %s", cfg.DBUser)
	}
}
