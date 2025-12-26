package config

import "testing"

func TestRead(t *testing.T) {
	cfg, err := Read()
	if err != nil {
		t.Fatal("Failed to read the config, ", err.Error())
	}

	if cfg.Port != ":8081" {
		t.Fatalf("expected port :8081, got: %s", cfg.Port)
	}

	if cfg.JWTSecret != "dev-only" {
		t.Fatalf("expected version dev-only, got: %s", cfg.JWTSecret)
	}
}