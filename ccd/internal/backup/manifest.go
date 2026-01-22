package backup

import (
	"encoding/json"
	"time"
)

const (
	ManifestVersion  = "1.0"
	ManifestFilename = "manifest.json"
)

type BackupManifest struct {
	Version   string      `json:"version"`
	Timestamp time.Time   `json:"timestamp"`
	TargetDir string      `json:"target_dir"`
	Files     []FileEntry `json:"files"`
}

type FileEntry struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

func NewManifest(timestamp time.Time, targetDir string) *BackupManifest {
	return &BackupManifest{
		Version:   ManifestVersion,
		Timestamp: timestamp,
		TargetDir: targetDir,
		Files:     []FileEntry{},
	}
}

func (m *BackupManifest) AddFile(path string, size int64) {
	m.Files = append(m.Files, FileEntry{Path: path, Size: size})
}

func (m *BackupManifest) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

func ParseManifest(data []byte) (*BackupManifest, error) {
	var m BackupManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
