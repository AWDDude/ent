package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandHome_TildeAlone(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	got, err := expandHome("~", home)
	require.NoError(t, err)
	assert.Equal(t, home, got)
}

func TestExpandHome_NoTilde(t *testing.T) {
	got, err := expandHome("/absolute/path", "/home/user")
	require.NoError(t, err)
	assert.Equal(t, "/absolute/path", got)
}

func TestExpandHome_Empty(t *testing.T) {
	got, err := expandHome("", "/home/user")
	require.NoError(t, err)
	assert.Equal(t, "", got)
}

func TestSettingsPath_Default(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	path, err := SettingsPath()
	require.NoError(t, err)
	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, ".config", "ent", "settings.json"), path)
}

func TestLoad_MalformedJSON(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("PROJECTS_ROOT", "")

	entDir := filepath.Join(dir, "ent")
	require.NoError(t, os.MkdirAll(entDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(entDir, "settings.json"), []byte("{not valid json"), 0o644))

	_, err := Load("")
	require.Error(t, err)
}
