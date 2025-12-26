package config

import "testing"

func TestRead(t *testing.T) {
	cfg, err := Read()
	if err != nil {
		t.Fatal("Failed to read the config, ", err.Error())
	}

	if cfg.Port != ":8000" {
		t.Fatalf("expected port :8000, got: %s", cfg.Port)
	}
	if cfg.APIVersion != "v1" {
		t.Fatalf("expected version v1, got: %s", cfg.APIVersion)
	}
}