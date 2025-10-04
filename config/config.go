// Package config
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	appName string = "ghget"
	file    string = "config.json"
)

type Config struct {
	Name   string `json:"name"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}

func New(name, repo, branch string) *Config {
	cfg := &Config{
		Name:   name,
		Repo:   repo,
		Branch: branch,
	}
	return cfg
}

func (cfg *Config) Read() error {
	configPath, err := cfg.getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Open(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("no config file in $HOME/.config/ghget/config.json")
		}
		return err
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("the config.json file is empty")
		}
		return err
	}

	return nil
}

// TODO: Add atomic file write

func (cfg *Config) Write() error {
	fmt.Println("Writing config file....")

	filePath, err := cfg.getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) IsNewFlagsEqualConfigValues() bool {
	oldCfg, err := readOldConfig()
	if err != nil {
		return false
	}

	return cfg.Name == oldCfg.Name &&
		cfg.Branch == oldCfg.Branch &&
		cfg.Repo == oldCfg.Repo
}

func readOldConfig() (Config, error) {
	cfg := Config{}

	configPath, err := cfg.getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, err
		}

		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return Config{}, err
		}

		return Config{}, err
	}

	return cfg, nil
}

func (cfg *Config) getConfigFilePath() (string, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigPath := filepath.Join(configPath, appName, file)
	return appConfigPath, nil
}

func (cfg *Config) String() string {
	format := "Username: %s\nRepo: %s\nBranch: %s\n"
	return fmt.Sprintf(format, cfg.Name, cfg.Repo, cfg.Branch)
}
