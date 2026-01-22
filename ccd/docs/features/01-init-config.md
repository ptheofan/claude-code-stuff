# Init Config Command - Technical Design Document

**Author:** Claude
**Date:** 2026-01-22
**Status:** Implemented

---

## 1. Problem Statement

Users need a way to generate a default configuration file without manually copying or creating one. Currently, users must either:
1. Find and copy the example `config.yaml` shipped with the binary
2. Manually create a config file from documentation
3. Run with defaults and have no visibility into available options

This creates friction for new users and makes it harder to discover and customize ccd behavior.

## 2. Goals & Non-Goals

### Goals
- Provide a CLI flag/command to generate a default `config.yaml` next to the executable
- Include comprehensive comments explaining each option
- Use correct default values matching `config.Default()`
- Warn if config file already exists (avoid accidental overwrites)
- Support `--force` to overwrite existing config

### Non-Goals
- Interactive config wizard (ask questions to build config)
- Config validation command (separate feature)
- Generating config to arbitrary paths (use standard location next to binary)

## 3. Proposed Solution

Add an `init` subcommand that generates a well-documented `config.yaml` file next to the executable.

```
ccd init           # Generate config.yaml next to binary
ccd init --force   # Overwrite existing config.yaml
```

### 3.1 Architecture

The feature fits cleanly into existing architecture:

```
cmd/ccd/main.go          # Add init subcommand
    │
    └──▶ internal/config/
              │
              ├── config.go      # Existing: Default(), Load()
              └── generate.go    # New: GenerateDefault() -> commented YAML string
```

**Layer responsibilities:**
- `cmd/ccd`: CLI flag handling, file write orchestration, user feedback
- `internal/config`: Pure generation of config content (no I/O)

### 3.2 Module Boundaries

| Module | Exposes | Consumes |
|--------|---------|----------|
| `internal/config` | `GenerateDefault() string` | Nothing new |
| `cmd/ccd` | `init` subcommand | `config.GenerateDefault()`, `output.*` |

### 3.3 Data Model

No new data types required. The feature generates a string (YAML content) from existing `Config` struct.

### 3.4 API Design

**New function in `internal/config/generate.go`:**

```go
// GenerateDefault returns a YAML string containing the default configuration
// with comprehensive comments explaining each option.
func GenerateDefault() string
```

**CLI subcommand:**

```
ccd init [flags]

Flags:
  --force   Overwrite existing config.yaml without prompting
```

**Output path determination:**

```go
// GetConfigOutputPath returns the path where config.yaml should be written.
// This is always next to the executable.
func GetConfigOutputPath(execPath string) string
```

### 3.5 Error Handling

| Error Class | When Thrown | Data Included |
|-------------|-------------|---------------|
| `ConfigExistsError` | Config file already exists and --force not specified | `Path string` |
| `ConfigWriteError` | Failed to write config file to disk | `Path string`, `Cause error` |
| `ExecutablePathError` | Cannot determine executable path | `Cause error` |

## 4. Alternatives Considered

| Option | Pros | Cons | Verdict |
|--------|------|------|---------|
| **A: `init` subcommand** | Clear intent, familiar pattern (git init, npm init) | Adds a command | ✅ Selected |
| **B: `--init` flag on root** | No new command | Mixes concerns (deploy vs init), confusing UX | Rejected |
| **C: `config init` subcommand** | Groups config-related commands | Over-engineered for single operation | Rejected |
| **D: Embed config in binary, extract with `--dump-config`** | Single binary, no separate file needed | Harder to customize, unusual pattern | Rejected |

## 5. Testing Strategy

### Unit Tests

**`internal/config/generate_test.go`:**
- `TestGenerateDefault_ContainsAllFields`: Generated YAML contains all Config fields
- `TestGenerateDefault_HasComments`: Output includes comment lines (starts with `#`)
- `TestGenerateDefault_ValidYAML`: Output parses as valid YAML
- `TestGenerateDefault_MatchesDefaults`: Parsed output matches `Default()` values
- `TestGetConfigOutputPath`: Returns correct path next to executable

### Integration Tests

**`cmd/ccd/init_test.go`:**
- `TestInitCommand_CreatesFile`: Running `init` creates config.yaml
- `TestInitCommand_ExistingFile_NoForce`: Returns error when file exists
- `TestInitCommand_ExistingFile_WithForce`: Overwrites when --force specified
- `TestInitCommand_OutputContent`: Written file matches `GenerateDefault()`

### E2E Tests

Not required - integration tests cover the critical paths.

## 6. Migration / Rollout Plan

- [ ] No feature flag needed (additive feature)
- [ ] No database migrations
- [ ] Fully backward compatible (new command, no changes to existing behavior)
- [ ] Update README with `ccd init` usage

## 7. Open Questions

- Should `init` also create the target directory (`~/.claude`) if it doesn't exist?
  - **Recommendation:** No, keep `init` focused on config generation only.

## 8. References

- Existing config module: `internal/config/config.go`
- Example config: `config.yaml` in repo root
- Similar patterns: `git init`, `npm init`, `go mod init`

---

## Appendix: Expected Output Format

```yaml
# CCD Configuration
# Generated by: ccd init
# Documentation: https://github.com/pt/ccd
#
# Place this file next to the binary, or at:
# - ~/.config/claude-deploy/config.yaml
# - ~/.claude-deploy/config.yaml

# Target directory for deployment
# Supports ~ for home directory
target: ~/.claude

# Files and directories to exclude from deployment
# These items in the source directory will not be copied
exclude:
  - deploy
  - LICENSE
  - README.md
  - .git
  - .idea
  - ccd

# Glob patterns to always ignore (never synced)
# Matched against filename only, not full path
ignore_patterns:
  - .DS_Store
  - Thumbs.db
  - "*.tmp"
  - "*.log"
  - "*.swp"
  - "*~"

# Backup configuration
backup:
  # Enable automatic backups before each deploy
  enabled: true

  # Directory to store backup snapshots
  dir: ~/.claude-backups

  # Maximum number of snapshots to keep (oldest pruned first)
  max_snapshots: 5

# Default sync mode
# - "merge": Add and update files only (safe)
# - "sync": Also delete files not in source (destructive)
default_mode: merge

# Prompt for confirmation before deleting files in sync mode
confirm_deletes: true
```
