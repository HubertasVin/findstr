package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/HubertasVin/findstr/models"
)

const defaultConfigJSON = `{
  "layout": {
    "align": "right",
    "autoWidth": true,
    "header": { "parts": ["---", " ", "{filepath}", ":"] },
    "match":  { "parts": ["{ln}", " | ", "{text}"] },
    "context":{ "parts": ["{ln}", " | ", "{text}"] }
  },
  "theme": {
    "styles": {
      "header":  { "fg": "#ffffff", "bold": true },
      "match":   { "fg": "#ffffff", "bg": "#008a00", "bold": true },
      "context": { "fg": "#cccccc", "bold": false }
    }
  }
}`

func LoadConfig() (models.CompiledLayout, models.Theme, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return models.CompiledLayout{}, models.Theme{}, err
	}
	path := filepath.Join(home, ".config", "findstr.conf")

	var raw []byte
	if b, err := os.ReadFile(path); err == nil {
		raw = b
	} else if os.IsNotExist(err) {
		raw = []byte(defaultConfigJSON)
	} else {
		return models.CompiledLayout{}, models.Theme{}, err
	}

	cl, theme, err := parseConfigRaw(string(raw))
	if err != nil {
		return models.CompiledLayout{}, models.Theme{}, err
	}

	return cl, theme, nil
}

func parseConfigRaw(raw string) (models.CompiledLayout, models.Theme, error) {
	var cfg models.ConfigJSON
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		var lj models.LayoutJSON
		if err2 := json.Unmarshal([]byte(raw), &lj); err2 != nil {
			return models.CompiledLayout{}, models.Theme{}, err
		}
		return CompileLayout(fillLayoutDefaults(lj)), defaultThemeResolved(), nil
	}
	return CompileLayout(fillLayoutDefaults(cfg.Layout)), resolveThemeWithDefaults(cfg.Theme), nil
}

func CreateDefaultConfig() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to resolve home directory: %w", err)
	}
	path := filepath.Join(home, ".config", "findstr.conf")

	if _, err := os.Stat(path); err == nil {
		return "", errors.New("config already exists: " + path)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("unable to check existing config: %w", err)
	}

	if err := os.WriteFile(path, []byte(defaultConfigJSON), 0o644); err != nil {
		return "", fmt.Errorf("failed to write config: %w", err)
	}
	return path, nil
}
