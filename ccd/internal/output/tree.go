package output

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileChange struct {
	Path      string
	Operation string // "create", "update", "delete"
	Size      int64
	ModTime   time.Time
	IsDir     bool
}

type TreeNode struct {
	Name     string
	Change   *FileChange
	Children map[string]*TreeNode
	IsDir    bool
}

func NewTreeNode(name string, isDir bool) *TreeNode {
	return &TreeNode{
		Name:     name,
		Children: make(map[string]*TreeNode),
		IsDir:    isDir,
	}
}

func BuildTree(changes []FileChange, rootPath string) *TreeNode {
	root := NewTreeNode(rootPath, true)

	for i := range changes {
		change := &changes[i]
		parts := strings.Split(change.Path, string(filepath.Separator))
		current := root

		for j, part := range parts {
			if part == "" {
				continue
			}
			isLast := j == len(parts)-1
			if _, exists := current.Children[part]; !exists {
				node := NewTreeNode(part, !isLast || change.IsDir)
				current.Children[part] = node
			}
			if isLast {
				current.Children[part].Change = change
			}
			current = current.Children[part]
		}
	}

	return root
}

func RenderTree(node *TreeNode, prefix string, isLast bool) string {
	var sb strings.Builder

	if node.Name != "" {
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		line := prefix + connector
		line += formatNode(node)
		sb.WriteString(line + "\n")
	}

	childPrefix := prefix
	if node.Name != "" {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	names := make([]string, 0, len(node.Children))
	for name := range node.Children {
		names = append(names, name)
	}
	sort.Strings(names)

	for i, name := range names {
		child := node.Children[name]
		isLastChild := i == len(names)-1
		sb.WriteString(RenderTree(child, childPrefix, isLastChild))
	}

	return sb.String()
}

func formatNode(node *TreeNode) string {
	var sb strings.Builder

	if node.Change != nil {
		switch node.Change.Operation {
		case "create":
			sb.WriteString(Colorize(Green, "[+] "))
		case "update":
			sb.WriteString(Colorize(Yellow, "[~] "))
		case "delete":
			sb.WriteString(Colorize(Red, "[-] "))
		}
	}

	if node.IsDir {
		sb.WriteString(Colorize(Cyan, node.Name+"/"))
	} else {
		sb.WriteString(node.Name)
	}

	if node.Change != nil && !node.IsDir {
		sb.WriteString(fmt.Sprintf(" (%s)", formatSize(node.Change.Size)))
		if node.Change.Operation == "delete" && !node.Change.ModTime.IsZero() {
			sb.WriteString(fmt.Sprintf(", modified %s", formatAge(node.Change.ModTime)))
		}
	}

	return sb.String()
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

func formatAge(t time.Time) string {
	days := int(time.Since(t).Hours() / 24)
	if days == 0 {
		return "today"
	} else if days == 1 {
		return "1d ago"
	}
	return fmt.Sprintf("%dd ago", days)
}

func PrintTreeHeader(targetPath string) {
	fmt.Println(Colorize(Cyan, targetPath+"/"))
}
