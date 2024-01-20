package main

import (
	"strings"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type UserConfig struct {
	Editor string `yaml:"editor"`
	EditorArgs []string `yaml:"editorArgs"`
	DefaultNotesDir string `yaml:"defaultNotesDir"`
}

func ReadConfig(path string) (UserConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return UserConfig{}, err
	}

	// default config values
	config := UserConfig{
		Editor: "vim",
		DefaultNotesDir: "~/Documents/notes",
	}

	// load yaml config only when it exists, otherwise, use defaults
	if !errors.Is(err, os.ErrNotExist) {
		err = yaml.Unmarshal(content, &config)
		if err != nil {
			return UserConfig{}, err
		}
	}

	config.DefaultNotesDir = expandUser(config.DefaultNotesDir)

	return config, nil
}

func expandUser(path string) string {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		// no need to check since this shouldn't fail
		home, _ := os.UserHomeDir()

		return filepath.Join(home, path[2:])
	}

	return path
}
