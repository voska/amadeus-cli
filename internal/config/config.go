package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey      string `json:"api_key,omitempty"`
	APISecret   string `json:"api_secret,omitempty"`
	Environment string `json:"environment,omitempty"`
}

func Dir() (string, error) {
	if d := os.Getenv("AMADEUS_CONFIG_DIR"); d != "" {
		return d, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "amadeus"), nil
}

func Load() (*Config, error) {
	dir, err := Dir()
	if err != nil {
		return &Config{}, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if os.IsNotExist(err) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0o600)
}

func (c *Config) ResolveAPIKey() string {
	if v := os.Getenv("AMADEUS_API_KEY"); v != "" {
		return v
	}
	return c.APIKey
}

func (c *Config) ResolveAPISecret() string {
	if v := os.Getenv("AMADEUS_API_SECRET"); v != "" {
		return v
	}
	return c.APISecret
}

func (c *Config) ResolveEnvironment(flagTest bool) string {
	if flagTest {
		return "test"
	}
	if c.Environment != "" {
		return c.Environment
	}
	return "test"
}

func BaseURL(env string) string {
	if env == "production" {
		return "https://api.amadeus.com"
	}
	return "https://test.api.amadeus.com"
}
