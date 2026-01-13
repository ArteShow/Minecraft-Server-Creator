package config

import "testing"

func TestRead(t *testing.T) {
	cfg, err := Read()
	if err != nil {
		t.Fatal("Failed to read the config, ", err.Error())
	}

	if cfg.Port != ":8002" {
		t.Fatalf("expected port :8002, got: %s", cfg.Port)
	}
}