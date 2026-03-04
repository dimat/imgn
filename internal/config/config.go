// Package config handles configuration loading from env vars, config file, and flags.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// ProviderConfig holds provider-specific configuration.
type ProviderConfig struct {
	APIKey string `mapstructure:"api_key"`
}

// Config holds the application configuration.
type Config struct {
	Provider  string                    `mapstructure:"provider"`
	Providers map[string]ProviderConfig `mapstructure:"providers"`
	Model     string                    `mapstructure:"model"`
	Aspect    string                    `mapstructure:"aspect"`
	Size      string                    `mapstructure:"size"`
	OutputDir string                    `mapstructure:"output_dir"`
}

// APIKey returns the API key for the active provider.
func (c Config) APIKey() string {
	if p, ok := c.Providers[c.Provider]; ok && p.APIKey != "" {
		return p.APIKey
	}
	return ""
}

// DefaultConfig returns configuration with default values.
func DefaultConfig() Config {
	return Config{
		Provider:  "google",
		Providers: map[string]ProviderConfig{},
		Model:     "flash2",
		Aspect:    "16:9",
		Size:      "2k",
		OutputDir: ".",
	}
}

// Load reads configuration from env vars and config file.
// Priority: env vars > config file > defaults.
func Load() (Config, error) {
	cfg := DefaultConfig()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	configDir, err := os.UserHomeDir()
	if err == nil {
		v.AddConfigPath(filepath.Join(configDir, ".config", "imgn"))
	}

	v.SetDefault("model", cfg.Model)
	v.SetDefault("aspect", cfg.Aspect)
	v.SetDefault("size", cfg.Size)
	v.SetDefault("output_dir", cfg.OutputDir)

	if err := v.ReadInConfig(); err != nil {
		// Config file not found is fine
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only warn about other errors, don't fail
			fmt.Fprintf(os.Stderr, "warning: error reading config: %v\n", err)
		}
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}

	// Ensure providers map is initialised
	if cfg.Providers == nil {
		cfg.Providers = map[string]ProviderConfig{}
	}

	// Env vars override config file. IMGN_API_KEY takes priority over GEMINI_API_KEY.
	// These override the google provider key specifically.
	if key := os.Getenv("IMGN_API_KEY"); key != "" {
		cfg.Providers["google"] = ProviderConfig{APIKey: key}
	} else if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		cfg.Providers["google"] = ProviderConfig{APIKey: key}
	}

	return cfg, nil
}

// ConfigFilePath returns the expected config file path.
func ConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.config/imgn/config.yaml"
	}
	return filepath.Join(home, ".config", "imgn", "config.yaml")
}
