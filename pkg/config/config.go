package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	ParentDirectory    string `mapstructure:"parent_directory"`
	ConfigFileLocation string
}

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		ParentDirectory: filepath.Join(homeDir, "projects"),
	}
}

// LoadConfig loads the configuration from file or creates a default one if it doesn't exist
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Set up config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".config", "tuv")
	configFile := filepath.Join(configDir, "config.yaml")
	config.ConfigFileLocation = configFile

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Return the default config without saving it
		// This will trigger the first run screen
		return config, nil
	}

	// Load existing config
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// Save writes the config to the config file
func (c *Config) Save() error {
	viper.SetConfigFile(c.ConfigFileLocation)
	viper.SetConfigType("yaml")

	viper.Set("parent_directory", c.ParentDirectory)

	return viper.WriteConfig()
}

// IsFirstRun checks if this is the first run of the application
func (c *Config) IsFirstRun() bool {
	// If the config file doesn't exist, it's the first run
	_, err := os.Stat(c.ConfigFileLocation)
	return os.IsNotExist(err)
}
