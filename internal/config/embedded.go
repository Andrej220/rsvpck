package config

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/azargarov/rsvpck/internal/domain"
)

//go:embed defaults/geconfig.yaml defaults/geconfig.json
var defaultsFS embed.FS

func LoadEmbedded() (domain.NetTestConfig, error) {
	// Prefer YAML, fall back to JSON if YAML missing or invalid
	if b, err := defaultsFS.ReadFile("defaults/geconfig.yaml"); err == nil {
		return parseConfigBytes(b, ".yaml")
	}
	if b, err := defaultsFS.ReadFile("defaults/geconfig.json"); err == nil {
		return parseConfigBytes(b, ".json")
	}
	return domain.NetTestConfig{}, fmt.Errorf("no embedded defaults found")
}

func LoadFromFileOrEmbedded(path string) (domain.NetTestConfig, error) {
	if path != "" {
		if b, err := defaultsFS.ReadFile(path); err == nil {
			// If someone embeds a different path inside defaultsFS
			ext := filepath.Ext(path)
			return parseConfigBytes(b, ext)
		}
		// Or from disk
		cfg, err := LoadFromFile(path)
		if err == nil {
			return cfg, nil
		}
		return domain.NetTestConfig{}, fmt.Errorf("read %s: %w", path, err)
	}
	return LoadEmbedded()
}
