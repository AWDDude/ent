package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunClone_LocalNameRejected(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	err := runClone(cloneCmd, "justanamenourl")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ent new")
}

func TestRunClone_InvalidURL(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECTS_ROOT", root)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	// A URL that parses but lacks org/project
	err := runClone(cloneCmd, "https://github.com/onlyorg")
	require.Error(t, err)
}
