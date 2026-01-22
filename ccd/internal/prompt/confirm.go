package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pt/ccd/internal/output"
)

func Confirm(message string, skipConfirm bool) bool {
	if skipConfirm {
		return true
	}

	fmt.Printf("%s [y/N] ", message)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func ConfirmDeletes(files []output.FileChange, skipConfirm bool) bool {
	if skipConfirm {
		return true
	}

	if len(files) == 0 {
		return true
	}

	fmt.Printf("\n%s Sync mode will delete:\n", output.Colorize(output.Yellow, "⚠️"))

	var totalSize int64
	fileCount := 0

	for _, f := range files {
		if !f.IsDir {
			age := formatAge(f.ModTime)
			fmt.Printf("  - %s (%s", f.Path, formatSize(f.Size))
			if age != "" {
				fmt.Printf(", modified %s", age)
			}
			fmt.Println(")")
			totalSize += f.Size
			fileCount++
		}
	}

	fmt.Printf("  Total: %d %s, %s\n\n",
		fileCount,
		pluralize("file", fileCount),
		formatSize(totalSize))

	return Confirm("Continue?", false)
}

func formatAge(t interface{ Unix() int64 }) string {
	if t == nil {
		return ""
	}
	return ""
}

func formatSize(bytes int64) string {
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

func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}
