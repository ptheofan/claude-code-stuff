# Sync Mappings - Technical Design Document

**Author:** Claude
**Date:** 2026-01-22
**Status:** Implemented

---

## 1. Problem Statement

The target directory (`~/.claude`) contains runtime data that CCD did not create:
- `projects/` - User's project data
- `statsig/` - Analytics cache
- `todos/` - User's todo files
- Other IDE-generated content

In `--sync` mode, CCD deletes files not present in source. This is dangerous because it would destroy user data that CCD should never touch.

**Current workaround:** Users avoid `--sync` mode entirely, losing the benefit of clean deployments.

**Root cause:** CCD treats the entire target directory as its domain, but it should only manage specific paths.

## 2. Goals & Non-Goals

### Goals
- Define explicit source→target mappings in config
- Only sync files that match a mapping
- In sync mode, only delete within mapped target paths
- Support both file and directory mappings
- Replace the current `exclude` config (simpler mental model)

### Non-Goals
- Wildcard/glob mappings (e.g., `*.md → docs/`)
- Bidirectional sync
- Conflict resolution between overlapping mappings
- Migration tool for existing configs (manual update)

## 3. Proposed Solution

Add a `mappings` config section that explicitly defines what gets synced and where:

```yaml
mappings:
  - source: CLAUDE.md
    target: CLAUDE.md
  - source: commands/
    target: commands/
  - source: settings.json
    target: settings.json

ignore_patterns:
  - .DS_Store
  - "*.tmp"
```

**Behavior when `mappings` is defined:**
1. Only mapped items sync - unmapped source files are skipped
2. Sync mode scoped to mappings - deletions only occur within mapped target paths
3. `exclude` is ignored - mappings are the authoritative source filter

**Behavior when `mappings` is empty/missing:**
- Legacy behavior preserved - sync everything (minus excludes)

### 3.1 Architecture

```
cmd/ccd/main.go
    │
    └──▶ internal/sync/
              │
              ├── sync.go       # Modified: Pass mappings to diff
              ├── diff.go       # Modified: Filter by mappings, scope deletions
              └── mapping.go    # New: Mapping resolution logic

         internal/config/
              │
              ├── config.go     # Modified: Add Mappings field
              └── types.go      # New: Mapping type definition
```

### 3.2 Module Boundaries

| Module | Exposes | Consumes |
|--------|---------|----------|
| `internal/config` | `Mapping` type, `Config.Mappings` | Nothing new |
| `internal/sync` | `ResolveMappings()`, `MappingSet` | `config.Mapping` |
| `cmd/ccd` | No new APIs | Uses existing `sync.Sync()` |

### 3.3 Data Model

**`internal/config/types.go`:**

```go
type Mapping struct {
    Source string `yaml:"source"` // Relative to working directory
    Target string `yaml:"target"` // Relative to target directory
}
```

**Modified `Config` struct:**

```go
type Config struct {
    Target          string       `yaml:"target"`
    Mappings        []Mapping    `yaml:"mappings"`        // New
    Exclude         []string     `yaml:"exclude"`         // Ignored when mappings defined
    IgnorePatterns  []string     `yaml:"ignore_patterns"`
    Backup          BackupConfig `yaml:"backup"`
    DefaultMode     string       `yaml:"default_mode"`
    ConfirmDeletes  bool         `yaml:"confirm_deletes"`
}
```

**`internal/sync/mapping.go`:**

```go
type ResolvedMapping struct {
    SourcePath string // Absolute path to source
    TargetPath string // Absolute path to target
    RelSource  string // Relative from source root
    RelTarget  string // Relative from target root
    IsDir      bool
}

type MappingSet struct {
    Items       []ResolvedMapping
    TargetPaths map[string]bool // Managed target paths (for lookup)
}
```

### 3.4 API Design

**`internal/sync/mapping.go`:**

```go
// ResolveMappings expands config mappings into concrete file paths.
// Directory mappings expand to include all children.
func ResolveMappings(sourceDir, targetDir string, mappings []config.Mapping) (*MappingSet, error)

// IsManagedPath returns true if targetRelPath falls under any mapping.
func (m *MappingSet) IsManagedPath(targetRelPath string) bool
```

**Modified `SyncOptions`:**

```go
type SyncOptions struct {
    SourceDir      string
    TargetDir      string
    Mappings       []config.Mapping  // New: nil = legacy mode
    Excludes       []string          // Ignored when Mappings set
    IgnorePatterns []string
    SyncMode       bool
    DryRun         bool
}
```

### 3.5 Error Handling

| Error Class | When Thrown | Data Included |
|-------------|-------------|---------------|
| `MappingSourceNotFoundError` | Source path in mapping doesn't exist | `Source`, `Mapping` |
| `MappingOverlapError` | Two mappings target the same path | `Path`, `Mapping1`, `Mapping2` |
| `InvalidMappingError` | Empty source or target | `Mapping`, `Reason` |

## 4. Alternatives Considered

| Option | Pros | Cons | Verdict |
|--------|------|------|---------|
| **A: Explicit mappings** | Precise control; safe by default | Config verbosity | ✅ Selected |
| **B: Target whitelist only** | Simpler config | Doesn't filter source | Rejected |
| **C: Target blacklist (protect paths)** | Non-breaking | Easy to miss paths | Rejected |
| **D: Separate "managed" directory** | Clean separation | Restructures user's source | Rejected |

## 5. Testing Strategy

### Unit Tests

**`internal/sync/mapping_test.go`:**
- `TestResolveMappings_SingleFile` - Maps single file correctly
- `TestResolveMappings_Directory` - Expands directory to all children
- `TestResolveMappings_NestedDirectory` - Handles nested dirs
- `TestResolveMappings_SourceNotFound` - Returns MappingSourceNotFoundError
- `TestResolveMappings_EmptyMapping` - Returns InvalidMappingError
- `TestMappingSet_IsManagedPath` - Identifies managed paths
- `TestMappingSet_IsManagedPath_Nested` - Nested paths under mapping

**`internal/sync/diff_test.go`:**
- `TestCalculateDiff_WithMappings_OnlyMappedFiles` - Unmapped skipped
- `TestCalculateDiff_WithMappings_DeletesScoped` - Deletes only within mapped
- `TestCalculateDiff_WithMappings_IgnorePatternsApply` - .DS_Store ignored
- `TestCalculateDiff_NilMappings_LegacyBehavior` - Backward compatible

### Integration Tests

**`cmd/ccd/sync_mappings_test.go`:**
- `TestSync_WithMappings_CreatesOnlyMapped` - Respects mappings
- `TestSync_WithMappings_DeletesOnlyWithinMapped` - Protected paths untouched
- `TestSync_NoMappings_LegacyBehavior` - Existing configs work

### E2E Tests

Not required - integration tests cover critical paths.

## 6. Migration / Rollout Plan

- [x] Backward compatible - Empty `mappings` = current behavior
- [ ] No database migrations
- [ ] Update `GenerateDefault()` with commented mappings example
- [ ] Add migration guidance to README

## 7. Open Questions

1. **Warn if sync mode without mappings?**
   - Recommendation: Yes, print warning about risk

2. **Support renaming (`source: foo.md` → `target: bar.md`)?**
   - Recommendation: Yes, same complexity, useful for reorganization

## 8. References

- Sync implementation: `internal/sync/sync.go`, `internal/sync/diff.go`
- Config module: `internal/config/config.go`

---

## Appendix: Example Config

```yaml
target: ~/.claude

mappings:
  - source: CLAUDE.md
    target: CLAUDE.md
  - source: commands/
    target: commands/
  - source: agents/
    target: agents/

ignore_patterns:
  - .DS_Store
  - "*.tmp"

backup:
  enabled: true
  dir: ~/.claude-backups
  max_snapshots: 5

default_mode: merge
confirm_deletes: true
```

**Result in sync mode:**
- `~/.claude/CLAUDE.md` - managed (can be deleted/updated)
- `~/.claude/commands/*` - managed (deletions possible)
- `~/.claude/projects/*` - **untouched** (not mapped)
- `~/.claude/statsig/*` - **untouched** (not mapped)
