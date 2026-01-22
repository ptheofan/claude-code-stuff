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

func CreateSnapshot(sourceDir, backupDir string) (*Snapshot, error) {
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

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
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

	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to clear target directory: %w", err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	for _, file := range reader.File {
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
