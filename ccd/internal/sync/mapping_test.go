package sync

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pt/ccd/internal/config"
)

func TestResolveMappings_SingleFile(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "CLAUDE.md", "content")

	mappings := []config.Mapping{
		{Source: "CLAUDE.md", Target: "CLAUDE.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ms == nil {
		t.Fatal("expected non-nil MappingSet")
	}
	if len(ms.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(ms.Items))
	}
	if ms.Items[0].IsDir {
		t.Error("expected IsDir=false for file mapping")
	}
	if ms.Items[0].RelSource != "CLAUDE.md" {
		t.Errorf("expected RelSource=CLAUDE.md, got %s", ms.Items[0].RelSource)
	}
}

func TestResolveMappings_Directory(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")
	createFile(t, sourceDir, "commands/test.md", "content")

	mappings := []config.Mapping{
		{Source: "commands", Target: "commands"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ms == nil {
		t.Fatal("expected non-nil MappingSet")
	}
	if len(ms.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(ms.Items))
	}
	if !ms.Items[0].IsDir {
		t.Error("expected IsDir=true for directory mapping")
	}
}

func TestResolveMappings_NestedDirectory(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "agents/nested")
	createFile(t, sourceDir, "agents/nested/file.md", "content")

	mappings := []config.Mapping{
		{Source: "agents/nested", Target: "agents/nested"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ms == nil {
		t.Fatal("expected non-nil MappingSet")
	}
	expectedSourcePath := filepath.Join(sourceDir, "agents/nested")
	if ms.Items[0].SourcePath != expectedSourcePath {
		t.Errorf("expected SourcePath=%s, got %s", expectedSourcePath, ms.Items[0].SourcePath)
	}
}

func TestResolveMappings_Renaming(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "old.md", "content")

	mappings := []config.Mapping{
		{Source: "old.md", Target: "new.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ms.Items[0].RelSource != "old.md" {
		t.Errorf("expected RelSource=old.md, got %s", ms.Items[0].RelSource)
	}
	if ms.Items[0].RelTarget != "new.md" {
		t.Errorf("expected RelTarget=new.md, got %s", ms.Items[0].RelTarget)
	}
}

func TestResolveMappings_SourceNotFound(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	mappings := []config.Mapping{
		{Source: "nonexistent.md", Target: "foo.md"},
	}

	_, err := ResolveMappings(sourceDir, targetDir, mappings)

	if err == nil {
		t.Fatal("expected error for nonexistent source")
	}
	var notFoundErr *MappingSourceNotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("expected MappingSourceNotFoundError, got %T: %v", err, err)
	}
	if notFoundErr.Source != "nonexistent.md" {
		t.Errorf("expected Source=nonexistent.md, got %s", notFoundErr.Source)
	}
}

func TestResolveMappings_EmptySource(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	mappings := []config.Mapping{
		{Source: "", Target: "foo.md"},
	}

	_, err := ResolveMappings(sourceDir, targetDir, mappings)

	if err == nil {
		t.Fatal("expected error for empty source")
	}
	var invalidErr *InvalidMappingError
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected InvalidMappingError, got %T: %v", err, err)
	}
	if invalidErr.Reason != "empty source" {
		t.Errorf("expected Reason='empty source', got %s", invalidErr.Reason)
	}
}

func TestResolveMappings_EmptyTarget(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "foo.md", "content")

	mappings := []config.Mapping{
		{Source: "foo.md", Target: ""},
	}

	_, err := ResolveMappings(sourceDir, targetDir, mappings)

	if err == nil {
		t.Fatal("expected error for empty target")
	}
	var invalidErr *InvalidMappingError
	if !errors.As(err, &invalidErr) {
		t.Errorf("expected InvalidMappingError, got %T: %v", err, err)
	}
	if invalidErr.Reason != "empty target" {
		t.Errorf("expected Reason='empty target', got %s", invalidErr.Reason)
	}
}

func TestResolveMappings_OverlappingTargets(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "source1.md", "content1")
	createFile(t, sourceDir, "source2.md", "content2")

	mappings := []config.Mapping{
		{Source: "source1.md", Target: "CLAUDE.md"},
		{Source: "source2.md", Target: "CLAUDE.md"},
	}

	_, err := ResolveMappings(sourceDir, targetDir, mappings)

	if err == nil {
		t.Fatal("expected error for overlapping targets")
	}
	var overlapErr *MappingOverlapError
	if !errors.As(err, &overlapErr) {
		t.Errorf("expected MappingOverlapError, got %T: %v", err, err)
	}
	if overlapErr.TargetPath != "CLAUDE.md" {
		t.Errorf("expected TargetPath=CLAUDE.md, got %s", overlapErr.TargetPath)
	}
}

func TestResolveMappings_NilMappings(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	ms, err := ResolveMappings(sourceDir, targetDir, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ms != nil {
		t.Errorf("expected nil MappingSet for nil mappings, got %v", ms)
	}
}

func TestResolveMappings_EmptySlice(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	ms, err := ResolveMappings(sourceDir, targetDir, []config.Mapping{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ms != nil {
		t.Errorf("expected nil MappingSet for empty mappings, got %v", ms)
	}
}

func TestMappingSet_IsManagedPath_ExactMatch(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")

	mappings := []config.Mapping{
		{Source: "commands", Target: "commands"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ms.IsManagedPath("commands/foo.md") {
		t.Error("expected commands/foo.md to be managed")
	}
}

func TestMappingSet_IsManagedPath_FileMapping(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "CLAUDE.md", "content")

	mappings := []config.Mapping{
		{Source: "CLAUDE.md", Target: "CLAUDE.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ms.IsManagedPath("CLAUDE.md") {
		t.Error("expected CLAUDE.md to be managed")
	}
}

func TestMappingSet_IsManagedPath_NestedUnderDir(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "agents")

	mappings := []config.Mapping{
		{Source: "agents", Target: "agents"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ms.IsManagedPath("agents/core/dev.md") {
		t.Error("expected agents/core/dev.md to be managed under agents/")
	}
}

func TestMappingSet_IsManagedPath_NotManaged(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")

	mappings := []config.Mapping{
		{Source: "commands", Target: "commands"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ms.IsManagedPath("projects/foo.md") {
		t.Error("expected projects/foo.md to NOT be managed")
	}
}

func TestMappingSet_IsManagedPath_PartialNameMatch(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")

	mappings := []config.Mapping{
		{Source: "commands", Target: "commands"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ms.IsManagedPath("commands-old/foo.md") {
		t.Error("expected commands-old/foo.md to NOT be managed (partial name match)")
	}
}

func TestMappingSet_IsManagedPath_NilSet(t *testing.T) {
	var ms *MappingSet = nil

	if !ms.IsManagedPath("anything.md") {
		t.Error("expected nil MappingSet to treat everything as managed (legacy mode)")
	}
}

func TestMappingSet_GetTargetPath_SameSourceTarget(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "CLAUDE.md", "content")

	mappings := []config.Mapping{
		{Source: "CLAUDE.md", Target: "CLAUDE.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := ms.GetTargetPath("CLAUDE.md")
	if result != "CLAUDE.md" {
		t.Errorf("expected CLAUDE.md, got %s", result)
	}
}

func TestMappingSet_GetTargetPath_Renamed(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "old.md", "content")

	mappings := []config.Mapping{
		{Source: "old.md", Target: "new.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := ms.GetTargetPath("old.md")
	if result != "new.md" {
		t.Errorf("expected new.md, got %s", result)
	}
}

func TestMappingSet_GetTargetPath_DirectoryChild(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "commands")

	mappings := []config.Mapping{
		{Source: "commands", Target: "commands"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := ms.GetTargetPath("commands/foo.md")
	if result != "commands/foo.md" {
		t.Errorf("expected commands/foo.md, got %s", result)
	}
}

func TestMappingSet_GetTargetPath_DirectoryRenamed(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "src")

	mappings := []config.Mapping{
		{Source: "src", Target: "dest"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := ms.GetTargetPath("src/nested/file.md")
	if result != "dest/nested/file.md" {
		t.Errorf("expected dest/nested/file.md, got %s", result)
	}
}

func TestMappingSet_GetTargetPath_NotMapped(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "mapped.md", "content")

	mappings := []config.Mapping{
		{Source: "mapped.md", Target: "mapped.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := ms.GetTargetPath("unmapped.md")
	if result != "" {
		t.Errorf("expected empty string for unmapped path, got %s", result)
	}
}

func TestMappingSet_GetTargetPath_NilSet(t *testing.T) {
	var ms *MappingSet = nil

	result := ms.GetTargetPath("anything.md")
	if result != "anything.md" {
		t.Errorf("expected anything.md for nil set (legacy mode), got %s", result)
	}
}

func TestMappingSet_GetSourcePath_SameSourceTarget(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "CLAUDE.md", "content")

	mappings := []config.Mapping{
		{Source: "CLAUDE.md", Target: "CLAUDE.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := ms.GetSourcePath("CLAUDE.md")
	if result != "CLAUDE.md" {
		t.Errorf("expected CLAUDE.md, got %s", result)
	}
}

func TestMappingSet_GetSourcePath_Renamed(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "old.md", "content")

	mappings := []config.Mapping{
		{Source: "old.md", Target: "new.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := ms.GetSourcePath("new.md")
	if result != "old.md" {
		t.Errorf("expected old.md, got %s", result)
	}
}

func TestMappingSet_GetSourcePath_DirectoryChild(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createDir(t, sourceDir, "src")

	mappings := []config.Mapping{
		{Source: "src", Target: "dest"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := ms.GetSourcePath("dest/nested/file.md")
	if result != "src/nested/file.md" {
		t.Errorf("expected src/nested/file.md, got %s", result)
	}
}

func TestMappingSet_IsSourceMapped_True(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "mapped.md", "content")

	mappings := []config.Mapping{
		{Source: "mapped.md", Target: "mapped.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ms.IsSourceMapped("mapped.md") {
		t.Error("expected mapped.md to be source mapped")
	}
}

func TestMappingSet_IsSourceMapped_False(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	createFile(t, sourceDir, "mapped.md", "content")

	mappings := []config.Mapping{
		{Source: "mapped.md", Target: "mapped.md"},
	}

	ms, err := ResolveMappings(sourceDir, targetDir, mappings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ms.IsSourceMapped("unmapped.md") {
		t.Error("expected unmapped.md to NOT be source mapped")
	}
}

func TestMappingSet_IsSourceMapped_NilSet(t *testing.T) {
	var ms *MappingSet = nil

	if !ms.IsSourceMapped("anything.md") {
		t.Error("expected nil MappingSet to treat everything as source mapped (legacy mode)")
	}
}

func createFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create parent dirs: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
}

func createDir(t *testing.T, parent, name string) {
	t.Helper()
	path := filepath.Join(parent, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
}
