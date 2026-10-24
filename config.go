package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Config holds user settings, loaded from <config dir>/noteit/config.json.
type Config struct {
	Path            string `json:"path"`
	Editor          string `json:"editor"`
	DefaultNotebook string `json:"defaultNotebook"`
}

// configPath returns the location of the config file, honouring
// $NOTEIT_CONFIG if set.
func configPath() (string, error) {
	if p := os.Getenv("NOTEIT_CONFIG"); p != "" {
		return p, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "noteit", "config.json"), nil
}

// loadConfig reads the config file, if present, and fills in defaults for
// anything it does not set. A missing config file is not an error.
func loadConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("finding home directory: %w", err)
	}

	var cfg Config
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
	case err != nil:
		return Config{}, fmt.Errorf("reading config: %w", err)
	default:
		if err := json.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parsing config %s: %w", path, err)
		}
	}

	if cfg.Path == "" {
		cfg.Path = filepath.Join(home, "noteit")
	} else if cfg.Path == "~" || strings.HasPrefix(cfg.Path, "~/") {
		cfg.Path = filepath.Join(home, cfg.Path[1:])
	}
	if cfg.DefaultNotebook == "" {
		cfg.DefaultNotebook = "default"
	}
	return cfg, nil
}

// editorCommand picks the editor to use: the config value, then $VISUAL,
// then $EDITOR, falling back to vi. The result is split on whitespace so
// values such as "code --wait" work.
func (c Config) editorCommand() []string {
	for _, e := range []string{c.Editor, os.Getenv("VISUAL"), os.Getenv("EDITOR")} {
		if fields := strings.Fields(e); len(fields) > 0 {
			return fields
		}
	}
	return []string{"vi"}
}
