package projects

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// noopGit is a GitRunner that records the calls made to it without actually
// running git.
type noopGit struct {
	calls []gitCall
}

type gitCall struct {
	dir  string
	args []string
}

func (n *noopGit) run(dir string, args ...string) error {
	n.calls = append(n.calls, gitCall{dir: dir, args: args})
	return nil
}

// newTestManager creates a Manager backed by a temp directory and a no-op git runner.
func newTestManager(t *testing.T) (*Manager, *noopGit, string) {
	t.Helper()
	root := t.TempDir()
	git := &noopGit{}
	m := &Manager{Root: root, RunGit: git.run}
	return m, git, root
}

// mkProject creates org/project directories in root.
func mkProject(t *testing.T, root, org, project string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Join(root, org, project), 0o755))
}

// ---- List ----------------------------------------------------------------

func TestList_Empty(t *testing.T) {
	m, _, _ := newTestManager(t)
	items, err := m.List()
	require.NoError(t, err)
	assert.Empty(t, items)
}

func TestList_MultipleProjects(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "myorg", "alpha")
	mkProject(t, root, "myorg", "beta")
	mkProject(t, root, "local", "gamma")

	items, err := m.List()
	require.NoError(t, err)
	assert.Equal(t, []string{"local/gamma", "myorg/alpha", "myorg/beta"}, items)
}

func TestList_IgnoresFiles(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "myorg", "repo")
	// A plain file at the org level should not appear.
	require.NoError(t, os.WriteFile(filepath.Join(root, "myorg", "somefile.txt"), []byte("x"), 0o644))

	items, err := m.List()
	require.NoError(t, err)
	assert.Equal(t, []string{"myorg/repo"}, items)
}

func TestList_RootMissing(t *testing.T) {
	m := &Manager{Root: "/nonexistent/path/xyz", RunGit: (&noopGit{}).run}
	_, err := m.List()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "projects root not found")
}

// ---- Resolve -------------------------------------------------------------

func TestResolve_ByProjectName(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "myorg", "myapp")
	mkProject(t, root, "other", "otherapp")

	paths, err := m.Resolve("myapp")
	require.NoError(t, err)
	require.Len(t, paths, 1)
	assert.Equal(t, filepath.Join(root, "myorg", "myapp"), paths[0])
}

func TestResolve_MultipleMatches(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "org1", "app")
	mkProject(t, root, "org2", "app")

	paths, err := m.Resolve("app")
	require.NoError(t, err)
	assert.Len(t, paths, 2)
}

func TestResolve_NoMatch(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "myorg", "myapp")

	paths, err := m.Resolve("notfound")
	require.NoError(t, err)
	assert.Empty(t, paths)
}

func TestResolve_PartialMatch(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "myorg", "my-cool-app")

	paths, err := m.Resolve("cool")
	require.NoError(t, err)
	require.Len(t, paths, 1)
	assert.Equal(t, filepath.Join(root, "myorg", "my-cool-app"), paths[0])
}

// ---- New -----------------------------------------------------------------

func TestNew_CreatesDirectory(t *testing.T) {
	m, git, root := newTestManager(t)

	path, err := m.New("myapp")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "local", "myapp"), path)
	assert.DirExists(t, path)
	require.Len(t, git.calls, 1)
	assert.Equal(t, []string{"init"}, git.calls[0].args)
}

func TestNew_AlreadyExists(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "local", "existing")

	_, err := m.New("existing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestNew_EmptyName(t *testing.T) {
	m, _, _ := newTestManager(t)
	_, err := m.New("")
	require.Error(t, err)
}

func TestNew_SlashInName(t *testing.T) {
	m, _, _ := newTestManager(t)
	_, err := m.New("org/repo")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "slashes")
}

// ---- Clone ---------------------------------------------------------------

func TestClone_SSH(t *testing.T) {
	m, git, root := newTestManager(t)

	path, err := m.Clone("git@github.com:myorg/myrepo.git")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "myorg", "myrepo"), path)
	require.Len(t, git.calls, 1)
	assert.Equal(t, []string{"clone", "git@github.com:myorg/myrepo.git", path}, git.calls[0].args)
}

func TestClone_HTTPS(t *testing.T) {
	m, git, root := newTestManager(t)

	path, err := m.Clone("https://github.com/acme/tool.git")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "acme", "tool"), path)
	require.Len(t, git.calls, 1)
	assert.Equal(t, []string{"clone", "https://github.com/acme/tool.git", path}, git.calls[0].args)
}

func TestClone_LocalNameRejected(t *testing.T) {
	m, _, _ := newTestManager(t)
	_, err := m.Clone("justanamenourl")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ent new")
}

func TestClone_AlreadyExists(t *testing.T) {
	m, _, root := newTestManager(t)
	mkProject(t, root, "myorg", "myrepo")

	_, err := m.Clone("git@github.com:myorg/myrepo.git")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}
