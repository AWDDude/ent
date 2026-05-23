package projects

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Constructor(t *testing.T) {
	m := New("/some/root")
	assert.Equal(t, "/some/root", m.Root)
	assert.NotNil(t, m.RunGit)
}

func TestList_SkipsTopLevelFiles(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "myorg", "repo1")
	// A plain file at the root level should not cause a crash or appear.
	require.NoError(t, os.WriteFile(filepath.Join(root, "README.md"), []byte("hello"), 0o644))

	items, err := m.List()
	require.NoError(t, err)
	assert.Equal(t, []string{"myorg/repo1"}, items)
}

func TestResolve_OrgSlashProject(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "myorg", "myapp")
	mkProject(t, root, "other", "myapp")

	// Slash-style search — matches against the full "org/project" string.
	paths, err := m.Resolve("myorg/myapp")
	require.NoError(t, err)
	// Should find at least one result and not error.
	assert.NotEmpty(t, paths)
}

func TestClone_EmptyURL(t *testing.T) {
	m, _, _ := newTestManager(t)
	_, err := m.Clone("")
	require.Error(t, err)
}

func TestClone_OrgDirCreated(t *testing.T) {
	m, git, root := newTestManager(t)

	_, err := m.Clone("https://github.com/neworg/newrepo.git")
	require.NoError(t, err)

	assert.DirExists(t, filepath.Join(root, "neworg"))
	require.Len(t, git.calls, 1)
	assert.Equal(t, "clone", git.calls[0].args[0])
}
