package main

import (
	"os"
	"testing"
)

func TestLoadConfig_APIKeyEnabledWithoutKey(t *testing.T) {
	dir := t.TempDir()
	envFile := dir + "\\.env"
	os.WriteFile(envFile, []byte("API_KEY_ENABLED=true\nAPI_KEY=\nSERVER_PORT=9999\n"), 0o644)
	originalCwd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalCwd)

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error when API_KEY_ENABLED=true but API_KEY is empty")
	}
}

func TestLoadConfig_APIKeyDisabled(t *testing.T) {
	dir := t.TempDir()
	envFile := dir + "\\.env"
	os.WriteFile(envFile, []byte("API_KEY_ENABLED=false\nAPI_KEY=\nSERVER_PORT=9999\n"), 0o644)
	originalCwd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalCwd)

	os.Unsetenv("API_KEY_ENABLED")
	os.Unsetenv("API_KEY")
	os.Unsetenv("SERVER_PORT")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKeyEnabled {
		t.Fatal("expected APIKeyEnabled to be false")
	}
}
