package config

import (
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGenerateDefault_ContainsAllFields(t *testing.T) {
	output := GenerateDefault()

	requiredFields := []string{
		"target:",
		"ignore_patterns:",
		"backup:",
		"enabled:",
		"dir:",
		"max_snapshots:",
		"default_mode:",
		"confirm_deletes:",
	}

	for _, field := range requiredFields {
		if !strings.Contains(output, field) {
			t.Errorf("GenerateDefault() missing field %q", field)
		}
	}
}

func TestGenerateDefault_HasComments(t *testing.T) {
	output := GenerateDefault()
	lines := strings.Split(output, "\n")

	var commentCount int
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			commentCount++
		}
	}

	if commentCount < 10 {
		t.Errorf("GenerateDefault() has too few comments: got %d, want at least 10", commentCount)
	}
}

func TestGenerateDefault_ValidYAML(t *testing.T) {
	output := GenerateDefault()

	var cfg Config
	if err := yaml.Unmarshal([]byte(output), &cfg); err != nil {
		t.Errorf("GenerateDefault() produced invalid YAML: %v", err)
	}
}

func TestGenerateDefault_MatchesDefaults(t *testing.T) {
	output := GenerateDefault()

	var cfg Config
	if err := yaml.Unmarshal([]byte(output), &cfg); err != nil {
		t.Fatalf("Failed to parse generated YAML: %v", err)
	}

	defaults := Default()

	if cfg.Target != defaults.Target {
		t.Errorf("Target = %q, want %q", cfg.Target, defaults.Target)
	}

	if cfg.DefaultMode != defaults.DefaultMode {
		t.Errorf("DefaultMode = %q, want %q", cfg.DefaultMode, defaults.DefaultMode)
	}

	if cfg.ConfirmDeletes != defaults.ConfirmDeletes {
		t.Errorf("ConfirmDeletes = %v, want %v", cfg.ConfirmDeletes, defaults.ConfirmDeletes)
	}

	if cfg.Backup.Enabled != defaults.Backup.Enabled {
		t.Errorf("Backup.Enabled = %v, want %v", cfg.Backup.Enabled, defaults.Backup.Enabled)
	}

	if cfg.Backup.Dir != defaults.Backup.Dir {
		t.Errorf("Backup.Dir = %q, want %q", cfg.Backup.Dir, defaults.Backup.Dir)
	}

	if cfg.Backup.MaxSnapshots != defaults.Backup.MaxSnapshots {
		t.Errorf("Backup.MaxSnapshots = %d, want %d", cfg.Backup.MaxSnapshots, defaults.Backup.MaxSnapshots)
	}

	if len(cfg.IgnorePatterns) != len(defaults.IgnorePatterns) {
		t.Errorf("IgnorePatterns length = %d, want %d\ngot: %v\nwant: %v",
			len(cfg.IgnorePatterns), len(defaults.IgnorePatterns),
			cfg.IgnorePatterns, defaults.IgnorePatterns)
	} else {
		for i, v := range cfg.IgnorePatterns {
			if v != defaults.IgnorePatterns[i] {
				t.Errorf("IgnorePatterns[%d] = %q, want %q", i, v, defaults.IgnorePatterns[i])
			}
		}
	}
}

func TestGetConfigOutputPath(t *testing.T) {
	tests := []struct {
		name     string
		execPath string
		want     string
	}{
		{
			name:     "binary in current directory",
			execPath: "/usr/local/bin/ccd",
			want:     "/usr/local/bin/config.yaml",
		},
		{
			name:     "binary in home directory",
			execPath: "/home/user/ccd",
			want:     "/home/user/config.yaml",
		},
		{
			name:     "binary with complex path",
			execPath: "/opt/tools/ccd/bin/ccd",
			want:     "/opt/tools/ccd/bin/config.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetConfigOutputPath(tt.execPath)
			if got != tt.want {
				t.Errorf("GetConfigOutputPath(%q) = %q, want %q", tt.execPath, got, tt.want)
			}
		})
	}
}

func TestGetConfigOutputPath_ReturnsAbsolutePath(t *testing.T) {
	execPath := "/some/path/to/ccd"
	got := GetConfigOutputPath(execPath)

	if !filepath.IsAbs(got) {
		t.Errorf("GetConfigOutputPath() returned non-absolute path: %q", got)
	}
}
