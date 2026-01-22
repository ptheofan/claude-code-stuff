package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/pt/ccd/internal/backup"
	"github.com/pt/ccd/internal/config"
	"github.com/pt/ccd/internal/output"
	"github.com/pt/ccd/internal/prompt"
	"github.com/pt/ccd/internal/sync"
)

var (
	version = "1.0.0"

	flagSync    bool
	flagDryRun  bool
	flagTarget  string
	flagNoColor bool
	flagYes     bool
	flagList    bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ccd",
		Short: "Claude Code Deploy - File synchronization tool",
		Long:  "Deploy claude-code-stuff configuration to target directory with tree-view output and rollback support.",
		RunE:  runDeploy,
	}

	rootCmd.Flags().BoolVar(&flagSync, "sync", false, "Remove files from destination that no longer exist in source")
	rootCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Preview changes without making them")
	rootCmd.Flags().StringVar(&flagTarget, "target", "", "Override target directory")
	rootCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Disable colored output")
	rootCmd.Flags().BoolVar(&flagYes, "yes", false, "Skip confirmation prompts")

	rollbackCmd := &cobra.Command{
		Use:   "rollback [timestamp]",
		Short: "Restore from a backup snapshot",
		Long:  "Restore the target directory from a previous backup snapshot.",
		RunE:  runRollback,
	}
	rollbackCmd.Flags().BoolVar(&flagList, "list", false, "List available snapshots")
	rootCmd.AddCommand(rollbackCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ccd version %s\n", version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	resetConfigCmd := &cobra.Command{
		Use:   "reset-config",
		Short: "Reset config.yaml to defaults",
		Long:  "Reset config.yaml to default values with comprehensive comments. Overwrites existing config.",
		RunE:  runResetConfig,
	}
	rootCmd.AddCommand(resetConfigCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runDeploy(cmd *cobra.Command, args []string) error {
	if flagNoColor {
		output.DisableColors()
	}

	execPath, err := os.Executable()
	if err != nil {
		execPath = os.Args[0]
	}

	configPath := config.GetConfigOutputPath(execPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		output.PrintInfo("No config.yaml found, generating default...")
		content := config.GenerateDefault()
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			output.PrintWarning(fmt.Sprintf("Failed to generate config: %v", err))
		} else {
			fmt.Printf("  Created: %s\n\n", configPath)
		}
	}

	cfg, err := config.Load(execPath)
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to load config: %v", err))
		return err
	}

	sourceDir, err := os.Getwd()
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to get working directory: %v", err))
		return err
	}

	targetDir := cfg.Target
	if flagTarget != "" {
		targetDir = config.ExpandPath(flagTarget)
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		output.PrintError(fmt.Sprintf("Target directory does not exist: %s", targetDir))
		return err
	}

	output.PrintMode(flagDryRun, flagSync)
	output.PrintPaths(sourceDir, targetDir)

	if flagSync && len(cfg.Mappings) == 0 {
		output.PrintWarning("Sync mode without mappings - ALL unmapped target files may be deleted")
	}

	syncResult, err := sync.Sync(sync.SyncOptions{
		SourceDir:      sourceDir,
		TargetDir:      targetDir,
		Mappings:       cfg.Mappings,
		IgnorePatterns: cfg.IgnorePatterns,
		SyncMode:       flagSync,
		DryRun:         true,
	})
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to calculate changes: %v", err))
		return err
	}

	if !syncResult.Summary.HasChanges() {
		output.PrintInfo("No changes detected")
		return nil
	}

	tree := output.BuildTree(syncResult.Changes, targetDir)
	output.PrintTreeHeader(targetDir)
	fmt.Print(output.RenderTree(tree, "", true))

	syncResult.Summary.Print()

	if flagDryRun {
		output.PrintSuccess(true)
		return nil
	}

	if flagSync && cfg.ConfirmDeletes {
		deletions := sync.GetDeletions(syncResult.Changes)
		if len(deletions) > 0 {
			if !prompt.ConfirmDeletes(deletions, flagYes) {
				output.PrintWarning("Aborted by user")
				return nil
			}
		}
	}

	if cfg.Backup.Enabled {
		fmt.Println()
		output.PrintInfo("Creating backup snapshot...")
		snapshot, err := backup.CreateSnapshot(targetDir, cfg.Backup.Dir)
		if err != nil {
			output.PrintWarning(fmt.Sprintf("Failed to create backup: %v", err))
		} else {
			fmt.Printf("  Backup created: %s (%s)\n", snapshot.Name, backup.FormatSize(snapshot.Size))

			pruned, err := backup.PruneSnapshots(cfg.Backup.Dir, cfg.Backup.MaxSnapshots)
			if err != nil {
				output.PrintWarning(fmt.Sprintf("Failed to prune old backups: %v", err))
			} else if len(pruned) > 0 {
				fmt.Printf("  Pruned %d old %s\n", len(pruned), pluralize("snapshot", len(pruned)))
			}
		}
	}

	fmt.Println()
	output.PrintInfo("Applying changes...")

	_, err = sync.Sync(sync.SyncOptions{
		SourceDir:      sourceDir,
		TargetDir:      targetDir,
		Mappings:       cfg.Mappings,
		IgnorePatterns: cfg.IgnorePatterns,
		SyncMode:       flagSync,
		DryRun:         false,
	})
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to sync: %v", err))
		return err
	}

	output.PrintSuccess(false)
	return nil
}

func runRollback(cmd *cobra.Command, args []string) error {
	if flagNoColor {
		output.DisableColors()
	}

	execPath, err := os.Executable()
	if err != nil {
		execPath = os.Args[0]
	}

	cfg, err := config.Load(execPath)
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to load config: %v", err))
		return err
	}

	targetDir := cfg.Target
	if flagTarget != "" {
		targetDir = config.ExpandPath(flagTarget)
	}

	if flagList {
		snapshots, err := backup.ListSnapshots(cfg.Backup.Dir)
		if err != nil {
			output.PrintError(fmt.Sprintf("Failed to list snapshots: %v", err))
			return err
		}

		if len(snapshots) == 0 {
			output.PrintInfo("No snapshots found")
			return nil
		}

		fmt.Println(output.Colorize(output.Blue, "Available snapshots:"))
		for _, s := range snapshots {
			age := formatAge(s.Timestamp)
			fmt.Printf("  %s (%s, %s)\n", s.Name, backup.FormatSize(s.Size), age)
		}
		return nil
	}

	var identifier string
	if len(args) > 0 {
		identifier = args[0]
	}

	snapshot, err := backup.FindSnapshot(cfg.Backup.Dir, identifier)
	if err != nil {
		output.PrintError(err.Error())
		return err
	}

	fmt.Printf("Restoring from: %s\n", output.Colorize(output.Cyan, snapshot.Name))
	fmt.Printf("Target: %s\n\n", output.Colorize(output.Blue, targetDir))

	if !prompt.Confirm(output.Colorize(output.Yellow, "⚠️")+" This will replace all files in the target. Continue?", flagYes) {
		output.PrintWarning("Aborted by user")
		return nil
	}

	output.PrintInfo("Restoring snapshot...")

	if err := backup.RestoreSnapshot(snapshot.Path, targetDir); err != nil {
		output.PrintError(fmt.Sprintf("Failed to restore: %v", err))
		return err
	}

	fmt.Printf("\n%s Restored successfully from %s\n",
		output.Colorize(output.Green, "✅"),
		snapshot.Name)

	return nil
}

func runResetConfig(cmd *cobra.Command, args []string) error {
	if flagNoColor {
		output.DisableColors()
	}

	execPath, err := os.Executable()
	if err != nil {
		return &config.ExecutablePathError{Cause: err}
	}

	configPath := config.GetConfigOutputPath(execPath)

	content := config.GenerateDefault()
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return &config.ConfigWriteError{Path: configPath, Cause: err}
	}

	output.PrintSuccess(false)
	fmt.Printf("  Reset: %s\n", configPath)
	return nil
}

func formatAge(t time.Time) string {
	duration := time.Since(t)
	hours := int(duration.Hours())

	if hours < 1 {
		return "just now"
	} else if hours < 24 {
		return fmt.Sprintf("%dh ago", hours)
	}

	days := hours / 24
	if days == 1 {
		return "1d ago"
	}
	return fmt.Sprintf("%dd ago", days)
}

func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func init() {
	cobra.OnInitialize(func() {
		if _, ok := os.LookupEnv("NO_COLOR"); ok {
			output.DisableColors()
		}
	})
}
