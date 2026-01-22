package output

import (
	"fmt"
)

type Summary struct {
	Created int
	Updated int
	Deleted int
}

func (s *Summary) Add(operation string) {
	switch operation {
	case "create":
		s.Created++
	case "update":
		s.Updated++
	case "delete":
		s.Deleted++
	}
}

func (s *Summary) HasChanges() bool {
	return s.Created > 0 || s.Updated > 0 || s.Deleted > 0
}

func (s *Summary) Print() {
	fmt.Println()
	fmt.Println(Colorize(Blue, "Summary:"))

	if !s.HasChanges() {
		fmt.Println("  No changes")
		return
	}

	if s.Created > 0 {
		fmt.Printf("  %s: %d %s\n",
			Colorize(Green, "Created"),
			s.Created,
			pluralize("file", s.Created))
	}
	if s.Updated > 0 {
		fmt.Printf("  %s: %d %s\n",
			Colorize(Yellow, "Updated"),
			s.Updated,
			pluralize("file", s.Updated))
	}
	if s.Deleted > 0 {
		fmt.Printf("  %s: %d %s\n",
			Colorize(Red, "Deleted"),
			s.Deleted,
			pluralize("file", s.Deleted))
	}
}

func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func PrintMode(isDryRun, isSync bool) {
	mode := "merge"
	if isSync {
		mode = "sync"
	}

	if isDryRun {
		fmt.Printf("%s DRY RUN: Analyzing changes (%s mode)\n\n",
			Colorize(Yellow, "🔍"),
			mode)
	} else {
		fmt.Printf("%s Deploying (%s mode)\n\n",
			Colorize(Blue, "🚀"),
			mode)
	}
}

func PrintPaths(source, target string) {
	fmt.Printf("Source: %s\n", Colorize(Blue, source))
	fmt.Printf("Target: %s\n\n", Colorize(Blue, target))
}

func PrintSuccess(isDryRun bool) {
	if isDryRun {
		fmt.Printf("\n%s Dry run completed. No changes made.\n",
			Colorize(Yellow, "✅"))
	} else {
		fmt.Printf("\n%s Deployment completed successfully!\n",
			Colorize(Green, "✅"))
	}
}

func PrintError(msg string) {
	fmt.Printf("%s %s\n", Colorize(Red, "❌"), msg)
}

func PrintWarning(msg string) {
	fmt.Printf("%s %s\n", Colorize(Yellow, "⚠️"), msg)
}

func PrintInfo(msg string) {
	fmt.Printf("%s %s\n", Colorize(Blue, "ℹ️"), msg)
}
