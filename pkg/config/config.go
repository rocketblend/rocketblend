package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type InstallConfig struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

type Config struct {
	InstallationDir string          `json:"installationDir"`
	Installed       []InstallConfig `json:"installed"`
	Remotes         []string        `json:"remotes"`
}

const (
	appName    string = "rocketblend"
	configName string = "config"
	configType string = "json"
)

var (
	configDir string
)

func LoadConfig() (Config, error) {
	vp := viper.New()
	var config Config

	vp.SetConfigType(configType)
	vp.SetConfigName(configName)
	vp.AddConfigPath(configDir)

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Config{}, err
		}

		return Config{}, fmt.Errorf("config file was found but another error ocurred: %v", err)
	}

	if err := vp.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("%v", err)
	}

	return config, nil
}

func SaveConfig(config Config) error {
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	filePath := filepath.Join(configDir, configName, configType)

	err = os.WriteFile(filePath, file, 0644)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func init() {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot find user home directory: %v", err))
	}

	configDir = filepath.Join(home, fmt.Sprintf(".%s", appName))
}
