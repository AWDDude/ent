package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCmd_NoProjects(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	buf := &bytes.Buffer{}
	listCmd.SetOut(buf)
	listCmd.SetErr(buf)

	err := listCmd.RunE(listCmd, nil)
	require.NoError(t, err)
	assert.Empty(t, buf.String())
}

func TestListCmd_WithProjects(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	require.NoError(t, os.MkdirAll(filepath.Join(root, "myorg", "alpha"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "myorg", "beta"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "local", "gamma"), 0o755))

	buf := &bytes.Buffer{}
	listCmd.SetOut(buf)

	err := listCmd.RunE(listCmd, nil)
	require.NoError(t, err)
	assert.Equal(t, "local/gamma\nmyorg/alpha\nmyorg/beta\n", buf.String())
}

func TestListCmd_Alias(t *testing.T) {
	assert.Contains(t, listCmd.Aliases, "ls")
}
