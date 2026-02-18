package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// BrokerConfig represents a named broker configuration
type BrokerConfig struct {
	Name    string   `json:"name"`
	Brokers []string `json:"brokers"`
}

// Config represents the application configuration
type Config struct {
	DefaultAlias string         `json:"defaultAlias,omitempty"`
	Brokers      []BrokerConfig `json:"brokers"`
}

// GetConfigPath returns the absolute path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".kafkaDebugger")
	return filepath.Join(configDir, "config.json"), nil
}

// EnsureConfigDir ensures the config directory exists
func EnsureConfigDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".kafkaDebugger")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{Brokers: []BrokerConfig{}}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetBrokerByAlias returns the broker configuration for the given alias
func (c *Config) GetBrokerByAlias(alias string) (*BrokerConfig, error) {
	for _, broker := range c.Brokers {
		if broker.Name == alias {
			return &broker, nil
		}
	}
	return nil, fmt.Errorf("broker alias '%s' not found", alias)
}

// AddBroker adds a new broker configuration
func (c *Config) AddBroker(name string, brokers []string) error {
	// Check if alias already exists
	for i, broker := range c.Brokers {
		if broker.Name == name {
			// Update existing
			c.Brokers[i].Brokers = brokers
			return nil
		}
	}

	// Add new
	c.Brokers = append(c.Brokers, BrokerConfig{
		Name:    name,
		Brokers: brokers,
	})

	return nil
}

// RemoveBroker removes a broker configuration by name
func (c *Config) RemoveBroker(name string) error {
	for i, broker := range c.Brokers {
		if broker.Name == name {
			c.Brokers = append(c.Brokers[:i], c.Brokers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("broker alias '%s' not found", name)
}

// ListBrokers returns all broker configurations
func (c *Config) ListBrokers() []BrokerConfig {
	return c.Brokers
}
