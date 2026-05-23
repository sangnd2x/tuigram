package xdg

import (
	"os"
	"path/filepath"
)

func ConfigDir() string {
	if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
		return filepath.Join(d, "tui-telegram")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "tui-telegram")
}

func DataDir() string {
	if d := os.Getenv("XDG_DATA_HOME"); d != "" {
		return filepath.Join(d, "tui-telegram")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "tui-telegram")
}
