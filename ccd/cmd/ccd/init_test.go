package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pt/ccd/internal/config"
)

func runResetConfigWithPath(configPath string) error {
	content := config.GenerateDefault()
	return os.WriteFile(configPath, []byte(content), 0644)
}

func TestRunResetConfig_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := runResetConfigWithPath(configPath)
	if err != nil {
		t.Fatalf("runResetConfig() failed: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("runResetConfig() did not create config file")
	}
}

func TestRunResetConfig_OverwritesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(configPath, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	err := runResetConfigWithPath(configPath)
	if err != nil {
		t.Fatalf("runResetConfig() failed: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if string(content) == "existing" {
		t.Error("runResetConfig() did not overwrite file")
	}
}

func TestRunResetConfig_OutputContent(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := runResetConfigWithPath(configPath)
	if err != nil {
		t.Fatalf("runResetConfig() failed: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	expected := config.GenerateDefault()
	if string(content) != expected {
		t.Errorf("Config file content does not match GenerateDefault()\ngot:\n%s\nwant:\n%s", string(content), expected)
	}
}
