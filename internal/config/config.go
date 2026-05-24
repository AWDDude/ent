// Package config loads ent configuration from settings.json, environment
// variables, and flag overrides. Precedence (highest first):
//
//  1. Explicit override passed via SetOverride (from --projects-root flag)
//  2. PROJECTS_ROOT environment variable
//  3. ~/.config/ent/settings.json (or $XDG_CONFIG_HOME/ent/settings.json)
//  4. Default: ~/projects
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// Settings is the structure of settings.json.
type Settings struct {
	ProjectsRoot string `json:"projects_root"`
}

// Config holds the resolved configuration for ent.
type Config struct {
	ProjectsRoot string
}

// Load resolves the configuration. flagOverride is the value of --projects-root
// (empty string if the flag was not set).
func Load(flagOverride string) (*Config, error) {
	cfg := &Config{}

	// 1. Start with the default.
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfg.ProjectsRoot = filepath.Join(home, "projects")

	// 2. Apply settings.json (missing file is silently skipped; parse errors are fatal).
	s, settingsErr := loadSettings()
	if settingsErr != nil {
		return nil, settingsErr
	}
	if s.ProjectsRoot != "" {
		cfg.ProjectsRoot = s.ProjectsRoot
	}

	// 3. Apply env var.
	if v := os.Getenv("PROJECTS_ROOT"); v != "" {
		cfg.ProjectsRoot = v
	}

	// 4. Apply flag override.
	if flagOverride != "" {
		cfg.ProjectsRoot = flagOverride
	}

	// Expand ~ in whatever value we ended up with.
	cfg.ProjectsRoot, err = expandHome(cfg.ProjectsRoot, home)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// SettingsPath returns the platform-appropriate path to settings.json.
// If ENT_CONFIG_PATH is set it takes precedence over all other location logic.
func SettingsPath() (string, error) {
	if p := os.Getenv("ENT_CONFIG_PATH"); p != "" {
		return p, nil
	}

	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "ent", "settings.json"), nil
	}

	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "ent", "settings.json"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ent", "settings.json"), nil
}

func loadSettings() (*Settings, error) {
	path, err := SettingsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// No config file yet — write a default one so the user has a starting point.
		defaults := &Settings{ProjectsRoot: "~/projects"}
		if writeErr := writeDefaultSettings(path, defaults); writeErr != nil {
			// Non-fatal: carry on with the in-memory defaults.
			return defaults, nil
		}
		return defaults, nil
	}
	if err != nil {
		return nil, err
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// writeDefaultSettings serialises defaults to JSON and writes it to path,
// creating any missing parent directories first.
func writeDefaultSettings(path string, defaults *Settings) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(defaults, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func expandHome(path, home string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}
	if path == "~" {
		return home, nil
	}
	if len(path) > 1 && path[:2] == "~/" {
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}
