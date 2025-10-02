// Package config
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ta7eralla/ghget/flags"
)

const (
	appName string = "ghget"
	file    string = "config.json"
)

type Config struct {
	User   string `json:"name"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}

func New() *Config {
	cfg := &Config{}
	return cfg
}

func (cfg *Config) Read() error {
	configPath, err := cfg.getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return err
	}

	return nil
}

// TODO: Add atomic file write

func (cfg *Config) Write(fcfg flags.FlagConfig) error {
	fmt.Println("Writing config file....")

	filePath, err := cfg.getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Chdir()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(fcfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) String() string {
	format := "Username: %s\nRepo: %s\nBranch: %s\n"
	return fmt.Sprintf(format, cfg.User, cfg.Repo, cfg.Branch)
}

func (cfg *Config) getConfigFilePath() (string, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigPath := filepath.Join(configPath, appName, file)
	return appConfigPath, nil
}
