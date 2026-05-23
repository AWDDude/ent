package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Default(t *testing.T) {
	t.Setenv("PROJECTS_ROOT", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir()) // point away from real config

	cfg, err := Load("")
	require.NoError(t, err)

	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, "projects"), cfg.ProjectsRoot)
}

func TestLoad_FromSettingsFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("PROJECTS_ROOT", "")

	entDir := filepath.Join(dir, "ent")
	require.NoError(t, os.MkdirAll(entDir, 0o755))

	settings := Settings{ProjectsRoot: "/tmp/myprojects"}
	data, err := json.Marshal(settings)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(entDir, "settings.json"), data, 0o644))

	cfg, err := Load("")
	require.NoError(t, err)
	assert.Equal(t, "/tmp/myprojects", cfg.ProjectsRoot)
}

func TestLoad_EnvOverridesSettings(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("PROJECTS_ROOT", "/env/projects")

	entDir := filepath.Join(dir, "ent")
	require.NoError(t, os.MkdirAll(entDir, 0o755))
	settings := Settings{ProjectsRoot: "/settings/projects"}
	data, _ := json.Marshal(settings)
	require.NoError(t, os.WriteFile(filepath.Join(entDir, "settings.json"), data, 0o644))

	cfg, err := Load("")
	require.NoError(t, err)
	assert.Equal(t, "/env/projects", cfg.ProjectsRoot)
}

func TestLoad_FlagOverridesEnv(t *testing.T) {
	t.Setenv("PROJECTS_ROOT", "/env/projects")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg, err := Load("/flag/projects")
	require.NoError(t, err)
	assert.Equal(t, "/flag/projects", cfg.ProjectsRoot)
}

func TestLoad_TildeExpansion(t *testing.T) {
	t.Setenv("PROJECTS_ROOT", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg, err := Load("~/myprojects")
	require.NoError(t, err)

	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, "myprojects"), cfg.ProjectsRoot)
}

func TestLoad_MissingSettingsFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir()) // empty dir, no settings.json
	t.Setenv("PROJECTS_ROOT", "")

	cfg, err := Load("")
	require.NoError(t, err)

	home, _ := os.UserHomeDir()
	assert.Equal(t, filepath.Join(home, "projects"), cfg.ProjectsRoot)
}

func TestSettingsPath_XDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/custom/xdg")
	path, err := SettingsPath()
	require.NoError(t, err)
	assert.Equal(t, "/custom/xdg/ent/settings.json", path)
}
