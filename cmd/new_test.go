package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunNew_Success_PrintsPathToStdoutAndMessageToStderr(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	newCmd.SetOut(stdout)
	newCmd.SetErr(stderr)
	t.Cleanup(func() {
		newCmd.SetOut(nil)
		newCmd.SetErr(nil)
	})

	err := runNew(newCmd, "myapp")
	require.NoError(t, err)

	expected := filepath.Join(root, "local", "myapp")
	assert.Equal(t, expected+"\n", stdout.String())
	assert.Contains(t, stderr.String(), "Created project:")
}

func TestRunNew_EmptyName_Errors(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	err := runNew(newCmd, "")
	require.Error(t, err)
}

func TestRunNew_SlashInName_Errors(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	err := runNew(newCmd, "org/project")
	require.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "slash")
}

func TestRunNew_DuplicateName_Errors(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	require.NoError(t, os.MkdirAll(filepath.Join(root, "local", "existing"), 0o755))

	err := runNew(newCmd, "existing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}
