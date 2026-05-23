package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all projects",
	Long:    `List all projects under the projects root, printed as org/project.`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList(cmd)
	},
}

func runList(cmd *cobra.Command) error {
	m, err := newManager()
	if err != nil {
		return err
	}

	items, err := m.List()
	if err != nil {
		return err
	}

	for _, item := range items {
		fmt.Fprintln(cmd.OutOrStdout(), item)
	}
	return nil
}
