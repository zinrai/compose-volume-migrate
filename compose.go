package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// ComposeFile represents docker-compose.yml structure
type ComposeFile struct {
	Volumes map[string]VolumeConfig `yaml:"volumes"`
}

// VolumeConfig represents volume configuration
type VolumeConfig struct {
	External interface{} `yaml:"external,omitempty"`
	Driver   string      `yaml:"driver,omitempty"`
}

// findComposeFile finds docker-compose.yml or docker-compose.yaml
func findComposeFile() (string, error) {
	candidates := []string{"docker-compose.yml", "docker-compose.yaml"}

	for _, name := range candidates {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
	}

	return "", fmt.Errorf("docker-compose.yml not found")
}

// parseComposeFile parses docker-compose.yml
func parseComposeFile(path string) (*ComposeFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read compose file: %w", err)
	}

	var compose ComposeFile
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, fmt.Errorf("failed to parse compose file: %w", err)
	}

	return &compose, nil
}

// getExternalVolumes returns list of external volume names
func getExternalVolumes(compose *ComposeFile) []string {
	var volumes []string

	for name, config := range compose.Volumes {
		if isExternal(config) {
			volumes = append(volumes, name)
		}
	}

	return volumes
}

// isExternal checks if volume is external
func isExternal(config VolumeConfig) bool {
	if config.External == nil {
		return false
	}

	switch v := config.External.(type) {
	case bool:
		return v
	case map[string]interface{}:
		// external: { name: xxx }
		return true
	default:
		return false
	}
}
