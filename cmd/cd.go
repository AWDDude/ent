package cmd

import (
	"errors"
	"fmt"

	"github.com/AWDDude/ent/internal/picker"
	"github.com/AWDDude/ent/internal/projects"
	"github.com/spf13/cobra"
)

var cdCmd = &cobra.Command{
	Use:   "cd [name]",
	Short: "Print the path of a project (shell wrapper does the cd)",
	Long: `Resolve a project by name and print its absolute path to stdout.

The shell integration wrapper (installed via 'ent init') intercepts this
output and calls 'cd' in the current shell. Without the wrapper, you can use
it directly:

  cd "$(ent cd myproject)"`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCd(cmd, args, defaultPicker())
	},
}

func runCd(cmd *cobra.Command, args []string, p picker.Picker) error {
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

	fmt.Fprintln(cmd.OutOrStdout(), path)
	return nil
}

// resolveSingle returns exactly one project path, using the picker when
// multiple matches exist or when no name is provided.
func resolveSingle(m *projects.Manager, p picker.Picker, name string) (string, error) {
	if name == "" {
		// No name given — fuzzy pick from everything.
		all, err := m.List()
		if err != nil {
			return "", err
		}
		if len(all) == 0 {
			return "", fmt.Errorf("no projects found")
		}
		selected, err := p.Pick("Select project:", all)
		if err != nil {
			if errors.Is(err, picker.ErrAborted) {
				return "", nil
			}
			return "", err
		}
		return m.Root + "/" + selected, nil
	}

	matches, err := m.Resolve(name)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no project found matching: %s", name)
	}
	if len(matches) == 1 {
		return matches[0], nil
	}

	// Multiple matches — let user pick.
	// Convert absolute paths back to relative org/project for display.
	root := m.Root + "/"
	display := make([]string, len(matches))
	for i, abs := range matches {
		display[i] = trimPrefix(abs, root)
	}

	selected, err := p.Pick("Multiple matches:", display)
	if err != nil {
		if errors.Is(err, picker.ErrAborted) {
			return "", nil
		}
		return "", err
	}
	return m.Root + "/" + selected, nil
}

func trimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
