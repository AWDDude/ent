// Package projects manages the on-disk project tree rooted at ProjectsRoot.
// The tree has the shape:
//
//	$ProjectsRoot/<org>/<project>
//
// All public functions accept a Manager that carries the root path and an
// executor for git commands, making them fully testable without touching real
// disk or spawning processes.
package projects

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AWDDude/ent/internal/urlparse"
)

// GitRunner is a function that runs a git command.  The default implementation
// uses exec.Command; tests can substitute a no-op or validator.
type GitRunner func(dir string, args ...string) error

// Manager holds configuration for project operations.
type Manager struct {
	Root      string
	RunGit    GitRunner
}

// New returns a Manager with the default (real) git runner.
func New(root string) *Manager {
	return &Manager{
		Root:   root,
		RunGit: defaultGitRunner,
	}
}

// List returns all org/project entries found under Root, sorted.
func (m *Manager) List() ([]string, error) {
	if _, err := os.Stat(m.Root); os.IsNotExist(err) {
		return nil, fmt.Errorf("projects root not found: %s", m.Root)
	}

	var results []string
	entries, err := os.ReadDir(m.Root)
	if err != nil {
		return nil, err
	}

	for _, orgEntry := range entries {
		if !orgEntry.IsDir() {
			continue
		}
		orgPath := filepath.Join(m.Root, orgEntry.Name())
		projEntries, err := os.ReadDir(orgPath)
		if err != nil {
			continue
		}
		for _, projEntry := range projEntries {
			if projEntry.IsDir() {
				results = append(results, orgEntry.Name()+"/"+projEntry.Name())
			}
		}
	}

	sort.Strings(results)
	return results, nil
}

// Resolve finds projects matching name and returns their full paths.
// If name contains a slash it matches against "org/project"; otherwise it
// matches only the project component.
func (m *Manager) Resolve(name string) ([]string, error) {
	all, err := m.List()
	if err != nil {
		return nil, err
	}

	var matches []string
	for _, entry := range all {
		if strings.Contains(name, "/") {
			// org/project style search
			if strings.Contains(entry, name) {
				matches = append(matches, filepath.Join(m.Root, entry))
			}
		} else {
			parts := strings.SplitN(entry, "/", 2)
			if len(parts) == 2 && strings.Contains(parts[1], name) {
				matches = append(matches, filepath.Join(m.Root, entry))
			}
		}
	}
	return matches, nil
}

// New creates a new local project at $Root/local/<name>, initialises a git
// repository, and returns its path.
func (m *Manager) New(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("project name is required")
	}
	if strings.Contains(name, "/") {
		return "", fmt.Errorf("project name must not contain slashes; use 'ent clone' for URLs")
	}

	target := filepath.Join(m.Root, "local", name)
	if _, err := os.Stat(target); err == nil {
		return target, fmt.Errorf("project already exists: %s", target)
	}

	if err := os.MkdirAll(target, 0o755); err != nil {
		return "", fmt.Errorf("could not create directory: %w", err)
	}

	if err := m.RunGit(target, "init"); err != nil {
		return "", fmt.Errorf("git init failed: %w", err)
	}

	return target, nil
}

// Clone clones a remote repository into $Root/<org>/<project> and returns the
// target path.
func (m *Manager) Clone(rawURL string) (string, error) {
	parsed, err := urlparse.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsed.CloneURL == "" {
		return "", fmt.Errorf("%q looks like a local name; use 'ent new' instead", rawURL)
	}

	target := filepath.Join(m.Root, parsed.Org, parsed.Project)
	if _, err := os.Stat(target); err == nil {
		return target, fmt.Errorf("project already exists: %s", target)
	}

	orgDir := filepath.Join(m.Root, parsed.Org)
	if err := os.MkdirAll(orgDir, 0o755); err != nil {
		return "", fmt.Errorf("could not create org directory: %w", err)
	}

	if err := m.RunGit(orgDir, "clone", parsed.CloneURL, target); err != nil {
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return target, nil
}

// Remove deletes the project directory at the given absolute path and removes
// the parent org directory if it becomes empty afterward.
func (m *Manager) Remove(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("could not remove project: %w", err)
	}
	// Clean up empty org dir.
	parent := filepath.Dir(path)
	entries, err := os.ReadDir(parent)
	if err == nil && len(entries) == 0 {
		_ = os.Remove(parent)
	}
	return nil
}

func defaultGitRunner(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
