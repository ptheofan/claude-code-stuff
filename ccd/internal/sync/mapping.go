package sync

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pt/ccd/internal/config"
)

// ResolvedMapping represents a validated source-to-target mapping.
type ResolvedMapping struct {
	SourcePath string // Absolute path to source
	TargetPath string // Absolute path to target
	RelSource  string // Relative from source root
	RelTarget  string // Relative from target root
	IsDir      bool
}

// MappingSet holds resolved mappings for efficient lookup.
type MappingSet struct {
	Items       []ResolvedMapping
	targetPaths map[string]bool
}

// ResolveMappings validates and expands config mappings.
// Returns nil if mappings is nil or empty (signals legacy mode).
func ResolveMappings(sourceDir, targetDir string, mappings []config.Mapping) (*MappingSet, error) {
	if len(mappings) == 0 {
		return nil, nil
	}

	ms := &MappingSet{
		Items:       make([]ResolvedMapping, 0, len(mappings)),
		targetPaths: make(map[string]bool),
	}

	for _, m := range mappings {
		if m.Source == "" {
			return nil, &InvalidMappingError{
				Mapping: formatMapping(m),
				Reason:  "empty source",
			}
		}
		if m.Target == "" {
			return nil, &InvalidMappingError{
				Mapping: formatMapping(m),
				Reason:  "empty target",
			}
		}

		targetKey := normalizeTargetKey(m.Target)
		if ms.targetPaths[targetKey] {
			for _, existing := range ms.Items {
				if normalizeTargetKey(existing.RelTarget) == targetKey {
					return nil, &MappingOverlapError{
						TargetPath: m.Target,
						Mapping1:   formatResolvedMapping(existing),
						Mapping2:   formatMapping(m),
					}
				}
			}
		}

		srcPath := filepath.Join(sourceDir, m.Source)
		info, err := os.Stat(srcPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, &MappingSourceNotFoundError{
					Source:  m.Source,
					Mapping: formatMapping(m),
				}
			}
			return nil, err
		}

		resolved := ResolvedMapping{
			SourcePath: srcPath,
			TargetPath: filepath.Join(targetDir, m.Target),
			RelSource:  m.Source,
			RelTarget:  m.Target,
			IsDir:      info.IsDir(),
		}

		ms.Items = append(ms.Items, resolved)
		ms.targetPaths[targetKey] = true
	}

	return ms, nil
}

// IsManagedPath returns true if targetRelPath falls under any mapping.
func (ms *MappingSet) IsManagedPath(targetRelPath string) bool {
	if ms == nil {
		return true
	}

	normalized := filepath.Clean(targetRelPath)
	for _, m := range ms.Items {
		mappingTarget := filepath.Clean(m.RelTarget)

		if m.IsDir {
			if normalized == mappingTarget ||
				strings.HasPrefix(normalized, mappingTarget+string(filepath.Separator)) {
				return true
			}
		} else {
			if normalized == mappingTarget {
				return true
			}
		}
	}
	return false
}

// GetTargetPath translates a source-relative path to its target-relative path.
// Returns empty string if path is not covered by any mapping.
func (ms *MappingSet) GetTargetPath(sourceRelPath string) string {
	if ms == nil {
		return sourceRelPath
	}

	normalized := filepath.Clean(sourceRelPath)
	for _, m := range ms.Items {
		mappingSource := filepath.Clean(m.RelSource)

		if m.IsDir {
			if normalized == mappingSource ||
				strings.HasPrefix(normalized, mappingSource+string(filepath.Separator)) {
				suffix := strings.TrimPrefix(normalized, mappingSource)
				return filepath.Clean(m.RelTarget + suffix)
			}
		} else {
			if normalized == mappingSource {
				return m.RelTarget
			}
		}
	}
	return ""
}

// GetSourcePath is the reverse of GetTargetPath - finds source for a target path.
func (ms *MappingSet) GetSourcePath(targetRelPath string) string {
	if ms == nil {
		return targetRelPath
	}

	normalized := filepath.Clean(targetRelPath)
	for _, m := range ms.Items {
		mappingTarget := filepath.Clean(m.RelTarget)

		if m.IsDir {
			if normalized == mappingTarget ||
				strings.HasPrefix(normalized, mappingTarget+string(filepath.Separator)) {
				suffix := strings.TrimPrefix(normalized, mappingTarget)
				return filepath.Clean(m.RelSource + suffix)
			}
		} else {
			if normalized == mappingTarget {
				return m.RelSource
			}
		}
	}
	return ""
}

// IsSourceMapped checks if a source-relative path is covered by any mapping.
func (ms *MappingSet) IsSourceMapped(sourceRelPath string) bool {
	if ms == nil {
		return true
	}
	return ms.GetTargetPath(sourceRelPath) != ""
}

func formatMapping(m config.Mapping) string {
	return m.Source + " -> " + m.Target
}

func formatResolvedMapping(m ResolvedMapping) string {
	return m.RelSource + " -> " + m.RelTarget
}

func normalizeTargetKey(target string) string {
	return filepath.Clean(strings.TrimSuffix(target, "/"))
}
