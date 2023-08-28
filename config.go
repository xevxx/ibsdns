package main

import (
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

func config() (credents Config, err error) {
	// Check for a local config.yaml in the current folder
	localConfigPath := "config.yaml"
	_, localErr := os.Stat(localConfigPath)
	if localErr == nil {
		localConfigFile, localErr := os.ReadFile(localConfigPath)
		if localErr != nil {
			return credents, localErr
		}

		err = yaml.Unmarshal(localConfigFile, &credents)
		if err != nil {
			return credents, err
		}

		return credents, nil
	}

	// Fallback to user-specific configuration directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return credents, err
	}

	configPath := filepath.Join(configDir, "ibsdns", "config.yaml")
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return credents, err
	}

	err = yaml.Unmarshal(configFile, &credents)
	if err != nil {
		return credents, err
	}

	return credents, nil
}

type Config struct {
	ApiKey   string `json:"apiKey"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	Url      string `json:"url"`
}
