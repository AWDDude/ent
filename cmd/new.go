package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new local project",
	Long: `Create a new local project at $projects_root/local/<name> and run git init.

Examples:
  ent new myapp
  ent new my-cool-service`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNew(cmd, args[0])
	},
}

func runNew(cmd *cobra.Command, name string) error {
	m, err := newManager()
	if err != nil {
		return err
	}

	path, err := m.New(name)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Created project: %s\n", path)
	fmt.Fprintln(cmd.OutOrStdout(), path)
	return nil
}
