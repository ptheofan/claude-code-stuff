package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type BackupConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Dir          string `yaml:"dir"`
	MaxSnapshots int    `yaml:"max_snapshots"`
}

type Config struct {
	Target          string       `yaml:"target"`
	Mappings        []Mapping    `yaml:"mappings"`
	IgnorePatterns  []string     `yaml:"ignore_patterns"`
	Backup          BackupConfig `yaml:"backup"`
	DefaultMode     string       `yaml:"default_mode"`
	ConfirmDeletes  bool         `yaml:"confirm_deletes"`
}

func Default() *Config {
	return &Config{
		Target: "~/.claude",
		IgnorePatterns: []string{
			".DS_Store",
			"Thumbs.db",
			".git",
			".idea",
			"*.tmp",
			"*.log",
			"*.swp",
			"*~",
		},
		Backup: BackupConfig{
			Enabled:      true,
			Dir:          "~/.claude-backups",
			MaxSnapshots: 5,
		},
		DefaultMode:    "merge",
		ConfirmDeletes: true,
	}
}

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func Load(execPath string) (*Config, error) {
	cfg := Default()

	locations := []string{
		filepath.Join(filepath.Dir(execPath), "config.yaml"),
		ExpandPath("~/.config/claude-deploy/config.yaml"),
		ExpandPath("~/.claude-deploy/config.yaml"),
	}

	var configPath string
	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			configPath = loc
			break
		}
	}

	if configPath == "" {
		cfg.Target = ExpandPath(cfg.Target)
		cfg.Backup.Dir = ExpandPath(cfg.Backup.Dir)
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.Target = ExpandPath(cfg.Target)
	cfg.Backup.Dir = ExpandPath(cfg.Backup.Dir)

	return cfg, nil
}
