package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pt/ccd/internal/config"
)

const (
	TimestampFormat = "2006-01-02_15-04-05"
	SnapshotPrefix  = "backup_"
	SnapshotSuffix  = ".zip"
)

type Snapshot struct {
	Path      string
	Name      string
	Timestamp time.Time
	Size      int64
}

func CreateSnapshot(targetDir, backupDir string, mappings []config.Mapping) (*Snapshot, error) {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now()
	name := SnapshotPrefix + timestamp.Format(TimestampFormat) + SnapshotSuffix
	zipPath := filepath.Join(backupDir, name)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	manifest := NewManifest(timestamp, targetDir)

	// Determine which paths to backup
	var pathsToBackup []string
	if len(mappings) > 0 {
		// Scoped backup: only mapped target paths
		for _, m := range mappings {
			targetPath := filepath.Join(targetDir, m.Target)
			if _, err := os.Stat(targetPath); err == nil {
				pathsToBackup = append(pathsToBackup, m.Target)
			}
		}
	} else {
		// Legacy: backup everything
		pathsToBackup = []string{""}
	}

	for _, basePath := range pathsToBackup {
		walkRoot := filepath.Join(targetDir, basePath)

		err = filepath.Walk(walkRoot, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(targetDir, path)
			if err != nil {
				return err
			}

			if relPath == "." {
				return nil
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Name = relPath
			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
				manifest.AddFile(relPath, info.Size())
			}

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(writer, file)
				if err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			os.Remove(zipPath)
			return nil, fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Write manifest.json to zip
	manifestData, err := manifest.ToJSON()
	if err != nil {
		os.Remove(zipPath)
		return nil, fmt.Errorf("failed to create manifest: %w", err)
	}

	manifestWriter, err := zipWriter.Create(ManifestFilename)
	if err != nil {
		os.Remove(zipPath)
		return nil, fmt.Errorf("failed to write manifest: %w", err)
	}
	if _, err := manifestWriter.Write(manifestData); err != nil {
		os.Remove(zipPath)
		return nil, fmt.Errorf("failed to write manifest: %w", err)
	}

	// Close zip to flush
	if err := zipWriter.Close(); err != nil {
		os.Remove(zipPath)
		return nil, fmt.Errorf("failed to finalize backup: %w", err)
	}

	info, _ := os.Stat(zipPath)

	return &Snapshot{
		Path:      zipPath,
		Name:      name,
		Timestamp: timestamp,
		Size:      info.Size(),
	}, nil
}

func ListSnapshots(backupDir string) ([]Snapshot, error) {
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var snapshots []Snapshot
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, SnapshotPrefix) || !strings.HasSuffix(name, SnapshotSuffix) {
			continue
		}

		timestampStr := strings.TrimPrefix(name, SnapshotPrefix)
		timestampStr = strings.TrimSuffix(timestampStr, SnapshotSuffix)

		timestamp, err := time.Parse(TimestampFormat, timestampStr)
		if err != nil {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		snapshots = append(snapshots, Snapshot{
			Path:      filepath.Join(backupDir, name),
			Name:      name,
			Timestamp: timestamp,
			Size:      info.Size(),
		})
	}

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp.After(snapshots[j].Timestamp)
	})

	return snapshots, nil
}

func RestoreSnapshot(snapshotPath, targetDir string) error {
	reader, err := zip.OpenReader(snapshotPath)
	if err != nil {
		return fmt.Errorf("failed to open snapshot: %w", err)
	}
	defer reader.Close()

	// Check for manifest
	var manifest *BackupManifest
	for _, file := range reader.File {
		if file.Name == ManifestFilename {
			rc, err := file.Open()
			if err != nil {
				return fmt.Errorf("failed to read manifest: %w", err)
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return fmt.Errorf("failed to read manifest: %w", err)
			}
			manifest, err = ParseManifest(data)
			if err != nil {
				return fmt.Errorf("failed to parse manifest: %w", err)
			}
			break
		}
	}

	if manifest != nil {
		// Scoped restore: only restore files listed in manifest
		// First, delete existing files that are in the manifest
		for _, entry := range manifest.Files {
			destPath := filepath.Join(targetDir, entry.Path)
			os.Remove(destPath) // Ignore error if doesn't exist
		}
	} else {
		// Legacy restore: nuke everything
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("failed to clear target directory: %w", err)
		}
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}
	}

	for _, file := range reader.File {
		// Skip manifest file
		if file.Name == ManifestFilename {
			continue
		}

		destPath := filepath.Join(targetDir, file.Name)

		if !strings.HasPrefix(destPath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path in archive: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		srcFile, err := file.Open()
		if err != nil {
			destFile.Close()
			return err
		}

		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		destFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func PruneSnapshots(backupDir string, maxSnapshots int) ([]string, error) {
	snapshots, err := ListSnapshots(backupDir)
	if err != nil {
		return nil, err
	}

	if len(snapshots) <= maxSnapshots {
		return nil, nil
	}

	var pruned []string
	for i := maxSnapshots; i < len(snapshots); i++ {
		if err := os.Remove(snapshots[i].Path); err != nil {
			return pruned, fmt.Errorf("failed to remove %s: %w", snapshots[i].Name, err)
		}
		pruned = append(pruned, snapshots[i].Name)
	}

	return pruned, nil
}

func FindSnapshot(backupDir, identifier string) (*Snapshot, error) {
	snapshots, err := ListSnapshots(backupDir)
	if err != nil {
		return nil, err
	}

	if len(snapshots) == 0 {
		return nil, fmt.Errorf("no snapshots found")
	}

	if identifier == "" {
		return &snapshots[0], nil
	}

	for _, s := range snapshots {
		if strings.Contains(s.Name, identifier) {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("snapshot not found: %s", identifier)
}

func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
	)

	switch {
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
