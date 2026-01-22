package sync

import (
	"path/filepath"
	"testing"

	"github.com/pt/ccd/internal/config"
)

func TestCalculateDiff_WithMappings_OnlyMappedFiles(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "mapped.md", "content")
	createFile(t, sourceDir, "unmapped.md", "content")

	mappings := []config.Mapping{
		{Source: "mapped.md", Target: "mapped.md"},
	}
	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("failed to resolve mappings: %v", err)
	}

	changes, err := CalculateDiff(sourceDir, targetDir, nil, false, ms)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Path != "mapped.md" {
		t.Errorf("expected mapped.md, got %s", changes[0].Path)
	}
}

func TestCalculateDiff_WithMappings_DirectoryMapping(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")
	createFile(t, sourceDir, "commands/a.md", "content")
	createFile(t, sourceDir, "commands/b.md", "content")
	createFile(t, sourceDir, "other/c.md", "content")

	mappings := []config.Mapping{
		{Source: "commands", Target: "commands"},
	}
	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("failed to resolve mappings: %v", err)
	}

	changes, err := CalculateDiff(sourceDir, targetDir, nil, false, ms)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	paths := make(map[string]bool)
	for _, c := range changes {
		paths[c.Path] = true
	}

	if !paths["commands"] {
		t.Error("expected commands directory in changes")
	}
	if !paths["commands/a.md"] {
		t.Error("expected commands/a.md in changes")
	}
	if !paths["commands/b.md"] {
		t.Error("expected commands/b.md in changes")
	}
	if paths["other"] || paths["other/c.md"] {
		t.Error("did not expect other/ files in changes")
	}
}

func TestCalculateDiff_WithMappings_Renaming(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "old.md", "content")

	mappings := []config.Mapping{
		{Source: "old.md", Target: "new.md"},
	}
	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("failed to resolve mappings: %v", err)
	}

	changes, err := CalculateDiff(sourceDir, targetDir, nil, false, ms)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Path != "new.md" {
		t.Errorf("expected target path new.md, got %s", changes[0].Path)
	}
}

func TestCalculateDiff_WithMappings_DeletesScoped(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")
	createDir(t, targetDir, "commands")
	createFile(t, targetDir, "commands/old.md", "content")
	createFile(t, targetDir, "projects/data.md", "content")

	mappings := []config.Mapping{
		{Source: "commands", Target: "commands"},
	}
	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("failed to resolve mappings: %v", err)
	}

	changes, err := CalculateDiff(sourceDir, targetDir, nil, true, ms)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var deletions []string
	for _, c := range changes {
		if c.Operation == "delete" {
			deletions = append(deletions, c.Path)
		}
	}

	if len(deletions) != 1 {
		t.Fatalf("expected 1 deletion, got %d: %v", len(deletions), deletions)
	}
	if deletions[0] != "commands/old.md" {
		t.Errorf("expected commands/old.md deletion, got %s", deletions[0])
	}
}

func TestCalculateDiff_WithMappings_IgnorePatternsStillApply(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")
	createFile(t, sourceDir, "commands/test.md", "content")
	createFile(t, sourceDir, "commands/.DS_Store", "content")

	mappings := []config.Mapping{
		{Source: "commands", Target: "commands"},
	}
	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("failed to resolve mappings: %v", err)
	}

	ignorePatterns := []string{".DS_Store"}
	changes, err := CalculateDiff(sourceDir, targetDir, ignorePatterns, false, ms)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, c := range changes {
		if filepath.Base(c.Path) == ".DS_Store" {
			t.Error("expected .DS_Store to be ignored")
		}
	}
}

func TestCalculateDiff_NilMappings_LegacyBehavior(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "file1.md", "content")
	createFile(t, sourceDir, "file2.md", "content")

	changes, err := CalculateDiff(sourceDir, targetDir, nil, false, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 2 {
		t.Errorf("expected 2 changes for legacy mode, got %d", len(changes))
	}
}

func TestCalculateDiff_NilMappings_DeletesEverywhere(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, targetDir, "extra1.md", "content")
	createFile(t, targetDir, "extra2.md", "content")

	changes, err := CalculateDiff(sourceDir, targetDir, nil, true, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deleteCount := 0
	for _, c := range changes {
		if c.Operation == "delete" {
			deleteCount++
		}
	}

	if deleteCount != 2 {
		t.Errorf("expected 2 deletions in legacy mode, got %d", deleteCount)
	}
}

