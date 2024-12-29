// internal/config/config.go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Auth     AuthConfig     `json:"auth"`
	App      AppConfig      `json:"app"`
}

type ServerConfig struct {
	Port         string `json:"port"`
	Environment  string `json:"environment"`
	AllowOrigins string `json:"allow_origins"`
}

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type AuthConfig struct {
	Secret string `json:"secret"`
}

type AppConfig struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	BaseURL     string `json:"base_url"`
}

// LoadConfig loads configuration from both JSON and environment variables
func LoadConfig(environment string) (*Config, error) {
	// Default configuration
	config := &Config{
		Server: ServerConfig{
			Port:        "8080",
			Environment: "development",
		},
		App: AppConfig{
			Title:       "My Blog & Portfolio",
			Description: "Personal blog and portfolio website",
			BaseURL:     "http://localhost:8080",
		},
	}

	// Load from config file if exists
	configFile := filepath.Join("config", environment+".json")
	if _, err := os.Stat(configFile); err == nil {
		file, err := os.Open(configFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(config); err != nil {
			return nil, err
		}
	}

	// Override with environment variables if they exist
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Server.Environment = env
	}
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		// Parse database URL and set config
	}

	return config, nil
}
