package config

import (
	"errors"
	"os"
	"path/filepath"

	"sub-muse/internal/keyring"

	"gopkg.in/yaml.v3"
)

var ErrNotConfigured = errors.New("not configured")

type Config struct {
	Configured bool   `yaml:"configured"`
	ServerURL  string `yaml:"server_url"`
	Username   string `yaml:"username"`
	ClientName string `yaml:"client_name"`
	Password   string `yaml:"-"`
}

func LoadConfig() (*Config, error) {
	config := &Config{
		ServerURL:  getEnv("SUBSONIC_URL", "http://localhost:4040"),
		Username:   getEnv("SUBSONIC_USERNAME", ""),
		Password:   getEnv("SUBSONIC_PASSWORD", ""),
		ClientName: getEnv("SUBSONIC_CLIENT_NAME", "sub-muse"),
	}

	if envUsername := os.Getenv("SUBSONIC_USERNAME"); envUsername != "" {
		config.Username = envUsername
	}
	if envPassword := os.Getenv("SUBSONIC_PASSWORD"); envPassword != "" {
		config.Password = envPassword
	}

	if config.Username == "" || config.Password == "" {
		fileConfig, err := loadConfigFile()
		if err != nil {
			if err == ErrNotConfigured {
				return nil, ErrNotConfigured
			}
			return nil, err
		}

		if fileConfig != nil {
			config.ServerURL = fileConfig.ServerURL
			config.Username = fileConfig.Username
			config.ClientName = fileConfig.ClientName
		}
	}

	if config.Password == "" {
		if pwd, err := keyring.GetPassword(config.Username); err == nil {
			config.Password = pwd
		}
	}

	if config.Username == "" || config.Password == "" {
		return nil, ErrNotConfigured
	}

	return config, nil
}

func loadConfigFile() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotConfigured
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if !cfg.Configured {
		return nil, ErrNotConfigured
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o600)
}

func IsConfigured() (bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return false, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return false, err
	}

	return cfg.Configured, nil
}

var configPath = ""

func GetConfigPath() (string, error) {
	if configPath != "" {
		return configPath, nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return configDir + "/sub-muse/config.yaml", nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
