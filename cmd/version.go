package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These variables are injected at build time by GoReleaser via -ldflags.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "ent %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
	},
}
