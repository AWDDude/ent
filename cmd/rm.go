package cmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/AWDDude/ent/internal/picker"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Permanently delete a project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRm(cmd, args, defaultPicker())
	},
}

func runRm(cmd *cobra.Command, args []string, p picker.Picker) error {
	m, err := newManager()
	if err != nil {
		return err
	}

	var name string
	if len(args) > 0 {
		name = args[0]
	}

	path, err := resolveSingle(m, p, name)
	if err != nil {
		return err
	}
	if path == "" {
		// User aborted the picker.
		return nil
	}

	// Derive org/project display name by stripping root prefix.
	displayName := trimPrefix(path, m.Root+"/")

	stderr := cmd.ErrOrStderr()
	fmt.Fprintln(stderr, "WARNING: This action is irreversible.")
	fmt.Fprintln(stderr, "Any changes that have not been pushed to a git server will be lost.")
	fmt.Fprintf(stderr, "\nProject to delete: %s\n", displayName)
	fmt.Fprintf(stderr, "\nType the project name to confirm: ")

	scanner := bufio.NewScanner(cmd.InOrStdin())
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if input != displayName {
		return fmt.Errorf("confirmation did not match — project not deleted")
	}

	if err := m.Remove(path); err != nil {
		return err
	}

	fmt.Fprintf(stderr, "Deleted: %s\n", path)
	return nil
}
