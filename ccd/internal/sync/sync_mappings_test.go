package sync

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pt/ccd/internal/config"
)

func TestSync_WithMappings_CreatesOnlyMapped(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "mapped.md", "mapped content")
	createFile(t, sourceDir, "unmapped.md", "unmapped content")

	opts := SyncOptions{
		SourceDir: sourceDir,
		TargetDir: targetDir,
		Mappings: []config.Mapping{
			{Source: "mapped.md", Target: "mapped.md"},
		},
		DryRun: false,
	}

	result, err := Sync(opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Summary.Created != 1 {
		t.Errorf("expected 1 create, got %d", result.Summary.Created)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "mapped.md")); os.IsNotExist(err) {
		t.Error("expected mapped.md to exist in target")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "unmapped.md")); !os.IsNotExist(err) {
		t.Error("expected unmapped.md to NOT exist in target")
	}
}

func TestSync_WithMappings_DeletesOnlyWithinMapped(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")
	createFile(t, sourceDir, "commands/new.md", "new content")

	createDir(t, targetDir, "commands")
	createFile(t, targetDir, "commands/old.md", "old content")
	createFile(t, targetDir, "projects/data.md", "protected data")

	opts := SyncOptions{
		SourceDir: sourceDir,
		TargetDir: targetDir,
		Mappings: []config.Mapping{
			{Source: "commands", Target: "commands"},
		},
		SyncMode: true,
		DryRun:   false,
	}

	result, err := Sync(opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Summary.Deleted != 1 {
		t.Errorf("expected 1 delete, got %d", result.Summary.Deleted)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "commands/old.md")); !os.IsNotExist(err) {
		t.Error("expected commands/old.md to be deleted")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "projects/data.md")); os.IsNotExist(err) {
		t.Error("expected projects/data.md to be PROTECTED (not deleted)")
	}
}

func TestSync_WithMappings_Renaming(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "old.md", "content")

	opts := SyncOptions{
		SourceDir: sourceDir,
		TargetDir: targetDir,
		Mappings: []config.Mapping{
			{Source: "old.md", Target: "new.md"},
		},
		DryRun: false,
	}

	result, err := Sync(opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Summary.Created != 1 {
		t.Errorf("expected 1 create, got %d", result.Summary.Created)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "new.md")); os.IsNotExist(err) {
		t.Error("expected new.md to exist in target")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "old.md")); !os.IsNotExist(err) {
		t.Error("expected old.md to NOT exist in target")
	}
}

func TestSync_NoMappings_LegacyBehavior(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "file1.md", "content1")
	createFile(t, sourceDir, "file2.md", "content2")

	opts := SyncOptions{
		SourceDir: sourceDir,
		TargetDir: targetDir,
		Mappings:  nil,
		DryRun:    false,
	}

	result, err := Sync(opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Summary.Created != 2 {
		t.Errorf("expected 2 creates in legacy mode, got %d", result.Summary.Created)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "file1.md")); os.IsNotExist(err) {
		t.Error("expected file1.md to exist in target")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "file2.md")); os.IsNotExist(err) {
		t.Error("expected file2.md to exist in target")
	}
}

func TestSync_MappingSourceNotFound_ReturnsError(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	opts := SyncOptions{
		SourceDir: sourceDir,
		TargetDir: targetDir,
		Mappings: []config.Mapping{
			{Source: "nonexistent.md", Target: "target.md"},
		},
		DryRun: false,
	}

	_, err := Sync(opts)

	if err == nil {
		t.Fatal("expected error for nonexistent source")
	}

	var notFoundErr *MappingSourceNotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected MappingSourceNotFoundError, got %T: %v", err, err)
	}
}

func TestSync_WithMappings_DirectoryRenaming(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "src")
	createFile(t, sourceDir, "src/nested/file.md", "content")

	opts := SyncOptions{
		SourceDir: sourceDir,
		TargetDir: targetDir,
		Mappings: []config.Mapping{
			{Source: "src", Target: "dest"},
		},
		DryRun: false,
	}

	_, err := Sync(opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "dest/nested/file.md")); os.IsNotExist(err) {
		t.Error("expected dest/nested/file.md to exist in target")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "src")); !os.IsNotExist(err) {
		t.Error("expected src/ to NOT exist in target")
	}
}

func TestSync_WithMappings_IgnorePatternsStillApply(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")
	createFile(t, sourceDir, "commands/test.md", "content")
	createFile(t, sourceDir, "commands/.DS_Store", "junk")

	opts := SyncOptions{
		SourceDir: sourceDir,
		TargetDir: targetDir,
		Mappings: []config.Mapping{
			{Source: "commands", Target: "commands"},
		},
		IgnorePatterns: []string{".DS_Store"},
		DryRun:         false,
	}

	_, err := Sync(opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "commands/.DS_Store")); !os.IsNotExist(err) {
		t.Error("expected .DS_Store to NOT exist in target (should be ignored)")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "commands/test.md")); os.IsNotExist(err) {
		t.Error("expected commands/test.md to exist in target")
	}
}

