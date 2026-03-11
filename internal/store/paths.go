package store

import (
	"os"
	"path/filepath"
)

// DefaultPath returns the platform-appropriate path for the SQLite database file.
//
// Platform resolution:
//   - Linux:   $XDG_CONFIG_HOME/gamevault-go/state.db  (fallback: ~/.config/...)
//   - Windows: %APPDATA%\gamevault-go\state.db
//   - macOS:   ~/Library/Application Support/gamevault-go/state.db
func DefaultPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(configDir, "gamevault-go", "state.db")
}
