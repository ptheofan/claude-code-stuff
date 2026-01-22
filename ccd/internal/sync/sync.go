package sync

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/pt/ccd/internal/config"
	"github.com/pt/ccd/internal/output"
)

type SyncOptions struct {
	SourceDir      string
	TargetDir      string
	Mappings       []config.Mapping
	IgnorePatterns []string
	SyncMode       bool
	DryRun         bool
}

type SyncResult struct {
	Changes []output.FileChange
	Summary output.Summary
}

func Sync(opts SyncOptions) (*SyncResult, error) {
	var mappingSet *MappingSet
	if len(opts.Mappings) > 0 {
		var err error
		mappingSet, err = ResolveMappings(opts.SourceDir, opts.TargetDir, opts.Mappings)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve mappings: %w", err)
		}
	}

	changes, err := CalculateDiff(
		opts.SourceDir,
		opts.TargetDir,
		opts.IgnorePatterns,
		opts.SyncMode,
		mappingSet,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate diff: %w", err)
	}

	result := &SyncResult{
		Changes: changes,
	}

	for _, c := range changes {
		result.Summary.Add(c.Operation)
	}

	if opts.DryRun {
		return result, nil
	}

	sortedChanges := sortChangesForExecution(changes)

	for _, change := range sortedChanges {
		var srcPath string
		if mappingSet != nil {
			srcRelPath := mappingSet.GetSourcePath(change.Path)
			if srcRelPath != "" {
				srcPath = filepath.Join(opts.SourceDir, srcRelPath)
			} else {
				srcPath = filepath.Join(opts.SourceDir, change.Path)
			}
		} else {
			srcPath = filepath.Join(opts.SourceDir, change.Path)
		}
		dstPath := filepath.Join(opts.TargetDir, change.Path)

		switch change.Operation {
		case "create", "update":
			if change.IsDir {
				if err := CopyDir(srcPath, dstPath); err != nil {
					return result, fmt.Errorf("failed to create directory %s: %w", change.Path, err)
				}
			} else {
				if err := CopyFile(srcPath, dstPath); err != nil {
					return result, fmt.Errorf("failed to copy %s: %w", change.Path, err)
				}
			}

		case "delete":
			if err := Delete(dstPath); err != nil {
				return result, fmt.Errorf("failed to delete %s: %w", change.Path, err)
			}
		}
	}

	return result, nil
}

func sortChangesForExecution(changes []output.FileChange) []output.FileChange {
	sorted := make([]output.FileChange, len(changes))
	copy(sorted, changes)

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Operation == "delete" && sorted[j].Operation != "delete" {
			return false
		}
		if sorted[i].Operation != "delete" && sorted[j].Operation == "delete" {
			return true
		}

		if sorted[i].Operation == "delete" && sorted[j].Operation == "delete" {
			return len(sorted[i].Path) > len(sorted[j].Path)
		}

		if sorted[i].IsDir && !sorted[j].IsDir {
			return true
		}
		if !sorted[i].IsDir && sorted[j].IsDir {
			return false
		}

		return sorted[i].Path < sorted[j].Path
	})

	return sorted
}
