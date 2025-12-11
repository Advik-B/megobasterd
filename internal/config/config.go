package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents application configuration
type Config struct {
	DownloadPath    string `mapstructure:"download_path"`
	MaxWorkers      int    `mapstructure:"max_workers"`
	UseProxy        bool   `mapstructure:"use_proxy"`
	ProxyHost       string `mapstructure:"proxy_host"`
	ProxyPort       int    `mapstructure:"proxy_port"`
	Language        string `mapstructure:"language"`
	Theme           string `mapstructure:"theme"`
	SmartProxy      bool   `mapstructure:"smart_proxy"`
	VerifyMAC       bool   `mapstructure:"verify_mac"`
	UseSlots        bool   `mapstructure:"use_slots"`
}

// Load reads configuration from file
func Load() (*Config, error) {
	viper.SetConfigName("megobasterd")
	viper.SetConfigType("yaml")
	
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	configDir := filepath.Join(homeDir, ".megobasterd")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")
	
	// Set defaults
	setDefaults()
	
	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create default
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return nil, err
			}
			// Save default config
			cfg := GetDefault()
			if err := cfg.Save(); err != nil {
				return nil, err
			}
			return cfg, nil
		}
		return nil, err
	}
	
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}

// Save writes configuration to file
func (c *Config) Save() error {
	viper.Set("download_path", c.DownloadPath)
	viper.Set("max_workers", c.MaxWorkers)
	viper.Set("use_proxy", c.UseProxy)
	viper.Set("proxy_host", c.ProxyHost)
	viper.Set("proxy_port", c.ProxyPort)
	viper.Set("language", c.Language)
	viper.Set("theme", c.Theme)
	viper.Set("smart_proxy", c.SmartProxy)
	viper.Set("verify_mac", c.VerifyMAC)
	viper.Set("use_slots", c.UseSlots)
	
	return viper.WriteConfig()
}

// GetDefault returns default configuration
func GetDefault() *Config {
	homeDir, _ := os.UserHomeDir()
	downloadPath := filepath.Join(homeDir, "Downloads")
	
	return &Config{
		DownloadPath: downloadPath,
		MaxWorkers:   6,
		UseProxy:     false,
		ProxyHost:    "",
		ProxyPort:    9999,
		Language:     "en",
		Theme:        "light",
		SmartProxy:   false,
		VerifyMAC:    false,
		UseSlots:     true,
	}
}

func setDefaults() {
	homeDir, _ := os.UserHomeDir()
	downloadPath := filepath.Join(homeDir, "Downloads")
	
	viper.SetDefault("download_path", downloadPath)
	viper.SetDefault("max_workers", 6)
	viper.SetDefault("use_proxy", false)
	viper.SetDefault("proxy_host", "")
	viper.SetDefault("proxy_port", 9999)
	viper.SetDefault("language", "en")
	viper.SetDefault("theme", "light")
	viper.SetDefault("smart_proxy", false)
	viper.SetDefault("verify_mac", false)
	viper.SetDefault("use_slots", true)
}
