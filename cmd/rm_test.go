package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AWDDude/ent/internal/picker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunRm_ConfirmedByArg(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	mkDir(t, root, "myorg", "myrepo")

	stderr := &bytes.Buffer{}
	stdin := strings.NewReader("myorg/myrepo\n")
	rmCmd.SetErr(stderr)
	rmCmd.SetIn(stdin)
	t.Cleanup(func() { rmCmd.SetErr(nil); rmCmd.SetIn(nil) })

	p := &picker.MockPicker{}
	err := runRm(rmCmd, []string{"myrepo"}, p)
	require.NoError(t, err)

	assert.NoDirExists(t, filepath.Join(root, "myorg", "myrepo"))
	assert.Contains(t, stderr.String(), "Deleted:")
}

func TestRunRm_WrongConfirmation_Aborts(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	mkDir(t, root, "myorg", "myrepo")

	stdin := strings.NewReader("wrongname\n")
	rmCmd.SetIn(stdin)
	t.Cleanup(func() { rmCmd.SetIn(nil) })

	p := &picker.MockPicker{}
	err := runRm(rmCmd, []string{"myrepo"}, p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "did not match")

	assert.DirExists(t, filepath.Join(root, "myorg", "myrepo"))
}

func TestRunRm_NoArg_UsesPicker(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	mkDir(t, root, "myorg", "myrepo")

	stderr := &bytes.Buffer{}
	stdin := strings.NewReader("myorg/myrepo\n")
	rmCmd.SetErr(stderr)
	rmCmd.SetIn(stdin)
	t.Cleanup(func() { rmCmd.SetErr(nil); rmCmd.SetIn(nil) })

	p := &picker.MockPicker{Selection: "myorg/myrepo"}
	err := runRm(rmCmd, []string{}, p)
	require.NoError(t, err)

	assert.NoDirExists(t, filepath.Join(root, "myorg", "myrepo"))
}

func TestRunRm_PickerAborted_NoOp(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	mkDir(t, root, "myorg", "myrepo")

	p := &picker.MockPicker{Err: picker.ErrAborted}
	err := runRm(rmCmd, []string{}, p)
	require.NoError(t, err)

	assert.DirExists(t, filepath.Join(root, "myorg", "myrepo"))
}

func TestRunRm_MultipleMatches_UsesPicker(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	mkDir(t, root, "org1", "app")
	mkDir(t, root, "org2", "app")

	stderr := &bytes.Buffer{}
	stdin := strings.NewReader("org1/app\n")
	rmCmd.SetErr(stderr)
	rmCmd.SetIn(stdin)
	t.Cleanup(func() { rmCmd.SetErr(nil); rmCmd.SetIn(nil) })

	p := &picker.MockPicker{Selection: "org1/app"}
	err := runRm(rmCmd, []string{"app"}, p)
	require.NoError(t, err)

	assert.NoDirExists(t, filepath.Join(root, "org1", "app"))
	assert.DirExists(t, filepath.Join(root, "org2", "app"))
}

func TestRunRm_OrgDirCleanedUpWhenEmpty(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	mkDir(t, root, "myorg", "onlyrepo")

	stdin := strings.NewReader("myorg/onlyrepo\n")
	rmCmd.SetIn(stdin)
	t.Cleanup(func() { rmCmd.SetIn(nil) })

	p := &picker.MockPicker{}
	err := runRm(rmCmd, []string{"onlyrepo"}, p)
	require.NoError(t, err)

	assert.NoDirExists(t, filepath.Join(root, "myorg"))
}

func TestRunRm_OrgDirRetainedWhenNotEmpty(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	mkDir(t, root, "myorg", "repo1")
	mkDir(t, root, "myorg", "repo2")

	stdin := strings.NewReader("myorg/repo1\n")
	rmCmd.SetIn(stdin)
	t.Cleanup(func() { rmCmd.SetIn(nil) })

	p := &picker.MockPicker{}
	err := runRm(rmCmd, []string{"repo1"}, p)
	require.NoError(t, err)

	assert.NoDirExists(t, filepath.Join(root, "myorg", "repo1"))
	assert.DirExists(t, filepath.Join(root, "myorg", "repo2"))
	assert.DirExists(t, filepath.Join(root, "myorg"))
}

func TestRunRm_NoProjectFound_Errors(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	p := &picker.MockPicker{}
	err := runRm(rmCmd, []string{"notfound"}, p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "notfound")
}

// Verify warning message content.
func TestRunRm_WarningMessageContent(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	require.NoError(t, os.MkdirAll(filepath.Join(root, "myorg", "myrepo"), 0o755))

	stderr := &bytes.Buffer{}
	stdin := strings.NewReader("myorg/myrepo\n")
	rmCmd.SetErr(stderr)
	rmCmd.SetIn(stdin)
	t.Cleanup(func() { rmCmd.SetErr(nil); rmCmd.SetIn(nil) })

	p := &picker.MockPicker{}
	_ = runRm(rmCmd, []string{"myrepo"}, p)

	out := stderr.String()
	assert.Contains(t, out, "irreversible")
	assert.Contains(t, out, "not been pushed")
	assert.Contains(t, out, "myorg/myrepo")
}
