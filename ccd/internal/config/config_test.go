package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMapping_YAMLUnmarshal(t *testing.T) {
	yamlContent := `
source: CLAUDE.md
target: CLAUDE.md
`
	var m Mapping
	err := yaml.Unmarshal([]byte(yamlContent), &m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Source != "CLAUDE.md" {
		t.Errorf("expected Source=CLAUDE.md, got %s", m.Source)
	}
	if m.Target != "CLAUDE.md" {
		t.Errorf("expected Target=CLAUDE.md, got %s", m.Target)
	}
}

func TestConfig_WithMappings_YAMLUnmarshal(t *testing.T) {
	yamlContent := `
target: ~/.claude
mappings:
  - source: CLAUDE.md
    target: CLAUDE.md
  - source: commands/
    target: commands/
`
	var cfg Config
	err := yaml.Unmarshal([]byte(yamlContent), &cfg)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Mappings) != 2 {
		t.Fatalf("expected 2 mappings, got %d", len(cfg.Mappings))
	}
	if cfg.Mappings[0].Source != "CLAUDE.md" {
		t.Errorf("expected first mapping source=CLAUDE.md, got %s", cfg.Mappings[0].Source)
	}
	if cfg.Mappings[1].Source != "commands/" {
		t.Errorf("expected second mapping source=commands/, got %s", cfg.Mappings[1].Source)
	}
}

func TestConfig_WithoutMappings_YAMLUnmarshal(t *testing.T) {
	yamlContent := `
target: ~/.claude
exclude:
  - .git
`
	var cfg Config
	err := yaml.Unmarshal([]byte(yamlContent), &cfg)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Mappings != nil {
		t.Errorf("expected nil Mappings, got %v", cfg.Mappings)
	}
}

func TestDefault_MappingsNil(t *testing.T) {
	cfg := Default()

	if cfg.Mappings != nil {
		t.Errorf("expected nil Mappings in default config, got %v", cfg.Mappings)
	}
}

func TestLoad_WithMappings(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	execPath := filepath.Join(tempDir, "ccd")

	configContent := `
target: ~/.claude
mappings:
  - source: CLAUDE.md
    target: CLAUDE.md
  - source: commands/
    target: commands/
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(execPath)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Mappings) != 2 {
		t.Fatalf("expected 2 mappings, got %d", len(cfg.Mappings))
	}
	if cfg.Mappings[0].Source != "CLAUDE.md" {
		t.Errorf("expected first mapping source=CLAUDE.md, got %s", cfg.Mappings[0].Source)
	}
}

func TestLoad_WithoutMappings_LegacyBehavior(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	execPath := filepath.Join(tempDir, "ccd")

	configContent := `
target: ~/.claude
exclude:
  - .git
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(execPath)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Mappings != nil {
		t.Errorf("expected nil Mappings for legacy config, got %v", cfg.Mappings)
	}
}
