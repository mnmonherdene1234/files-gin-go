package main

import (
	"os"
	"path/filepath"
	"testing"
)

// unsetEnv temporarily removes environment variables for the duration of the
// test and restores their original values (or removes them) on cleanup.
func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		old, had := os.LookupEnv(key)
		os.Unsetenv(key)
		t.Cleanup(func() {
			if had {
				os.Setenv(key, old)
			} else {
				os.Unsetenv(key)
			}
		})
	}
}

func TestLoadConfig_APIKeyEnabledWithoutKey(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("API_KEY_ENABLED=true\nAPI_KEY=\nSERVER_PORT=9999\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	unsetEnv(t, "API_KEY_ENABLED", "API_KEY", "SERVER_PORT")

	originalCwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(originalCwd)

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error when API_KEY_ENABLED=true but API_KEY is empty")
	}
}

func TestLoadConfig_APIKeyDisabled(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("API_KEY_ENABLED=false\nAPI_KEY=\nSERVER_PORT=9999\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	unsetEnv(t, "API_KEY_ENABLED", "API_KEY", "SERVER_PORT")

	originalCwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(originalCwd)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKeyEnabled {
		t.Fatal("expected APIKeyEnabled to be false")
	}
}
