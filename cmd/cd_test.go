package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AWDDude/ent/internal/picker"
	"github.com/AWDDude/ent/internal/projects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTestManager(t *testing.T) (*projects.Manager, string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	m := projects.New(root)
	return m, root
}

func mkDir(t *testing.T, root string, parts ...string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Join(append([]string{root}, parts...)...), 0o755))
}

func TestResolveSingle_NoArgs_PickerCalled(t *testing.T) {
	m, root := makeTestManager(t)
	mkDir(t, root, "myorg", "myrepo")

	p := &picker.MockPicker{Selection: "myorg/myrepo"}
	path, err := resolveSingle(m, p, "")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "myorg", "myrepo"), path)
}

func TestResolveSingle_ExactMatch(t *testing.T) {
	m, root := makeTestManager(t)
	mkDir(t, root, "myorg", "myrepo")

	p := &picker.MockPicker{}
	path, err := resolveSingle(m, p, "myrepo")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "myorg", "myrepo"), path)
}

func TestResolveSingle_MultipleMatches_PickerCalled(t *testing.T) {
	m, root := makeTestManager(t)
	mkDir(t, root, "org1", "app")
	mkDir(t, root, "org2", "app")

	p := &picker.MockPicker{Selection: "org1/app"}
	path, err := resolveSingle(m, p, "app")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "org1", "app"), path)
}

func TestResolveSingle_NoMatch(t *testing.T) {
	m, root := makeTestManager(t)
	mkDir(t, root, "myorg", "myrepo")

	p := &picker.MockPicker{}
	_, err := resolveSingle(m, p, "notfound")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "notfound")
}

func TestResolveSingle_Aborted(t *testing.T) {
	m, root := makeTestManager(t)
	mkDir(t, root, "myorg", "myrepo")

	p := &picker.MockPicker{Err: picker.ErrAborted}
	path, err := resolveSingle(m, p, "")
	// Aborted is not an error — user just cancelled.
	require.NoError(t, err)
	assert.Empty(t, path)
}

func TestRunCd_PrintsPath(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	mkDir(t, root, "myorg", "myrepo")

	buf := &bytes.Buffer{}
	cdCmd.SetOut(buf)

	p := &picker.MockPicker{Selection: "myorg/myrepo"}
	err := runCd(cdCmd, []string{"myrepo"}, p)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "myorg", "myrepo"), strings.TrimSpace(buf.String()))
}
