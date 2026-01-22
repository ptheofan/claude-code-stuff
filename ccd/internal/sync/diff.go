package sync

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pt/ccd/internal/output"
)

func CalculateDiff(sourceDir, targetDir string, ignorePatterns []string, syncMode bool, mappings *MappingSet) ([]output.FileChange, error) {
	var changes []output.FileChange

	sourceFiles := make(map[string]os.FileInfo)
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(sourceDir, path)
		if relPath == "." {
			return nil
		}

		baseName := filepath.Base(path)

		if shouldIgnore(baseName, ignorePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if mappings != nil && !mappings.IsSourceMapped(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		sourceFiles[relPath] = info
		return nil
	})
	if err != nil {
		return nil, err
	}

	targetFiles := make(map[string]os.FileInfo)
	if _, err := os.Stat(targetDir); err == nil {
		err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, _ := filepath.Rel(targetDir, path)
			if relPath == "." {
				return nil
			}

			baseName := filepath.Base(path)
			if shouldIgnore(baseName, ignorePatterns) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			targetFiles[relPath] = info
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	for relPath, srcInfo := range sourceFiles {
		targetRelPath := relPath
		if mappings != nil {
			targetRelPath = mappings.GetTargetPath(relPath)
		}

		targetInfo, existsInTarget := targetFiles[targetRelPath]

		if !existsInTarget {
			changes = append(changes, output.FileChange{
				Path:      targetRelPath,
				Operation: "create",
				Size:      srcInfo.Size(),
				ModTime:   srcInfo.ModTime(),
				IsDir:     srcInfo.IsDir(),
			})
		} else if !srcInfo.IsDir() && !targetInfo.IsDir() {
			if srcInfo.ModTime().After(targetInfo.ModTime()) || srcInfo.Size() != targetInfo.Size() {
				changes = append(changes, output.FileChange{
					Path:      targetRelPath,
					Operation: "update",
					Size:      srcInfo.Size(),
					ModTime:   srcInfo.ModTime(),
					IsDir:     false,
				})
			}
		}
	}

	if syncMode {
		for relPath, targetInfo := range targetFiles {
			if mappings != nil && !mappings.IsManagedPath(relPath) {
				continue
			}

			sourceRelPath := relPath
			if mappings != nil {
				sourceRelPath = mappings.GetSourcePath(relPath)
			}

			if _, existsInSource := sourceFiles[sourceRelPath]; !existsInSource {
				baseName := filepath.Base(relPath)
				if !shouldIgnore(baseName, ignorePatterns) {
					changes = append(changes, output.FileChange{
						Path:      relPath,
						Operation: "delete",
						Size:      targetInfo.Size(),
						ModTime:   targetInfo.ModTime(),
						IsDir:     targetInfo.IsDir(),
					})
				}
			}
		}
	}

	return changes, nil
}

func shouldIgnore(name string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
		if strings.HasPrefix(pattern, "*") && strings.HasSuffix(name, pattern[1:]) {
			return true
		}
	}
	return false
}

func GetDeletions(changes []output.FileChange) []output.FileChange {
	var deletions []output.FileChange
	for _, c := range changes {
		if c.Operation == "delete" {
			deletions = append(deletions, c)
		}
	}
	return deletions
}

func CalculateTotalSize(changes []output.FileChange) int64 {
	var total int64
	for _, c := range changes {
		if !c.IsDir {
			total += c.Size
		}
	}
	return total
}

func FormatDeleteSummary(deletions []output.FileChange) string {
	var sb strings.Builder
	var totalSize int64

	for _, d := range deletions {
		if !d.IsDir {
			age := formatAge(d.ModTime)
			sb.WriteString("  - " + d.Path)
			sb.WriteString(" (")
			sb.WriteString(formatSize(d.Size))
			if age != "" {
				sb.WriteString(", modified ")
				sb.WriteString(age)
			}
			sb.WriteString(")\n")
			totalSize += d.Size
		}
	}

	fileCount := 0
	for _, d := range deletions {
		if !d.IsDir {
			fileCount++
		}
	}

	sb.WriteString("  Total: ")
	sb.WriteString(string(rune('0') + rune(fileCount)))
	if fileCount != 1 {
		sb.WriteString(" files")
	} else {
		sb.WriteString(" file")
	}
	sb.WriteString(", ")
	sb.WriteString(formatSize(totalSize))

	return sb.String()
}

func formatAge(t time.Time) string {
	days := int(time.Since(t).Hours() / 24)
	if days == 0 {
		return "today"
	} else if days == 1 {
		return "1d ago"
	}
	return string(rune('0')+rune(days/10)) + string(rune('0')+rune(days%10)) + "d ago"
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
	)

	if bytes >= MB {
		return string(rune('0'+bytes/MB)) + "." + string(rune('0'+(bytes%MB)/(MB/10))) + " MB"
	} else if bytes >= KB {
		return string(rune('0'+bytes/KB)) + "." + string(rune('0'+(bytes%KB)/(KB/10))) + " KB"
	}
	return string(rune('0'+bytes)) + " B"
}
