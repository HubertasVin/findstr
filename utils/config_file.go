package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/HubertasVin/findstr/models"
	"github.com/pelletier/go-toml/v2"
)

const defaultConfigTOML = `
[layout]
align = "right"
autoWidth = true

[layout.header]
parts = ["---", " ", "{filepath}", ":"]

[layout.match]
parts = ["{ln}", " | ", "{text}"]

[layout.context]
parts = ["{ln}", " | ", "{text}"]

[theme.styles.header]
fg = "#ffffff"
bold = true

[theme.styles.match]
fg = "#5f5faf"
bold = true

[theme.styles.context]
fg = "#cccccc"
bold = false
`

func LoadConfig() (models.CompiledLayout, models.Theme, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return models.CompiledLayout{}, models.Theme{}, err
	}
	path := filepath.Join(home, ".config", "findstr.toml")

	var raw []byte
	if b, err := os.ReadFile(path); err == nil {
		raw = b
	} else if os.IsNotExist(err) {
		raw = []byte(defaultConfigTOML)
	} else {
		return models.CompiledLayout{}, models.Theme{}, err
	}

	var cfg models.ConfigJSON
	if err := toml.Unmarshal(raw, &cfg); err != nil {
		return models.CompiledLayout{}, models.Theme{}, err
	}

	return CompileLayout(fillLayoutDefaults(cfg.Layout)), resolveThemeWithDefaults(cfg.Theme), nil
}

func CreateDefaultConfig() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to resolve home directory: %w", err)
	}
	path := filepath.Join(home, ".config", "findstr.toml")

	if _, err := os.Stat(path); err == nil {
		return "", errors.New("config already exists: " + path)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("unable to check existing config: %w", err)
	}

	if err := os.WriteFile(path, []byte(defaultConfigTOML), 0o644); err != nil {
		return "", fmt.Errorf("failed to write config: %w", err)
	}
	return path, nil
}
