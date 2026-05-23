package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <url>",
	Short: "Clone a remote git repository",
	Long: `Clone a remote git repository into $projects_root/<org>/<project>.

Both SSH and HTTPS URLs are supported:

  ent clone git@github.com:myorg/myrepo.git
  ent clone https://github.com/myorg/myrepo.git`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runClone(cmd, args[0])
	},
}

func runClone(cmd *cobra.Command, rawURL string) error {
	m, err := newManager()
	if err != nil {
		return err
	}

	path, err := m.Clone(rawURL)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Cloned into: %s\n", path)
	return nil
}
