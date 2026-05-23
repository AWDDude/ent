// Package cmd wires up all cobra commands for the ent CLI.
package cmd

import (
	"github.com/AWDDude/ent/internal/config"
	"github.com/AWDDude/ent/internal/picker"
	"github.com/AWDDude/ent/internal/projects"
	"github.com/spf13/cobra"
)

var projectsRoot string

// rootCmd is the base command. All subcommands are registered on it.
var rootCmd = &cobra.Command{
	Use:   "ent",
	Short: "ent — project manager",
	Long: `ent organises your projects under a single root directory.

Projects are stored as $projects_root/<org>/<project>.

Add shell integration by running:
  eval "$(ent completion zsh)"    # zsh
  eval "$(ent completion bash)"   # bash
  ent completion fish | source    # fish
  ent completion powershell       # paste into $PROFILE`,
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Execute is the entry point called by main.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&projectsRoot,
		"projects-root", "",
		"override the projects root directory",
	)

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(versionCmd)
}

// loadConfig resolves config for use inside command Run functions.
func loadConfig() (*config.Config, error) {
	return config.Load(projectsRoot)
}

// newManager returns a Manager wired to the resolved projects root.
func newManager() (*projects.Manager, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return projects.New(cfg.ProjectsRoot), nil
}

// defaultPicker returns the real fuzzy picker; commands accept a Picker so
// tests can inject a mock.
func defaultPicker() picker.Picker {
	return picker.New()
}
