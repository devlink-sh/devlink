package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	P2PPort    int    `json:"p2p_port"`
	P2PNetwork string `json:"p2p_network"`

	EncryptionKey string `json:"encryption_key"`

	DataDir string `json:"data_dir"`

	DefaultExpiry string `json:"default_expiry"`
	MaxFileSize   int64  `json:"max_file_size"`

	LogLevel string `json:"log_level"`
	LogFile  string `json:"log_file"`
}

func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return &Config{
		P2PPort:       8080,
		P2PNetwork:    "devlink",
		DataDir:       filepath.Join(homeDir, ".devlink"),
		DefaultExpiry: "1h",
		MaxFileSize:   1024 * 1024,
		LogLevel:      "info",
		LogFile:       "",
	}
}

func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	if err := config.loadFromEnv(); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func (c *Config) loadFromEnv() error {
	if port := os.Getenv("DEVLINK_P2P_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.P2PPort = p
		}
	}

	if network := os.Getenv("DEVLINK_P2P_NETWORK"); network != "" {
		c.P2PNetwork = network
	}

	if key := os.Getenv("DEVLINK_ENCRYPTION_KEY"); key != "" {
		c.EncryptionKey = key
	}

	if dataDir := os.Getenv("DEVLINK_DATA_DIR"); dataDir != "" {
		c.DataDir = dataDir
	}

	if expiry := os.Getenv("DEVLINK_DEFAULT_EXPIRY"); expiry != "" {
		c.DefaultExpiry = expiry
	}

	if maxSize := os.Getenv("DEVLINK_MAX_FILE_SIZE"); maxSize != "" {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			c.MaxFileSize = size
		}
	}

	if logLevel := os.Getenv("DEVLINK_LOG_LEVEL"); logLevel != "" {
		c.LogLevel = strings.ToLower(logLevel)
	}

	if logFile := os.Getenv("DEVLINK_LOG_FILE"); logFile != "" {
		c.LogFile = logFile
	}

	return nil
}

func (c *Config) validate() error {
	if c.P2PPort < 1 || c.P2PPort > 65535 {
		return fmt.Errorf("invalid P2P port: %d (must be between 1-65535)", c.P2PPort)
	}

	if c.P2PNetwork == "" {
		return fmt.Errorf("P2P network name cannot be empty")
	}

	if c.DataDir == "" {
		return fmt.Errorf("data directory cannot be empty")
	}

	if _, err := time.ParseDuration(c.DefaultExpiry); err != nil {
		return fmt.Errorf("invalid default expiry format: %s", c.DefaultExpiry)
	}

	if c.MaxFileSize <= 0 {
		return fmt.Errorf("max file size must be positive")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	return nil
}

func (c *Config) EnsureDataDir() error {
	return os.MkdirAll(c.DataDir, 0755)
}

func (c *Config) GetDataDir() string {
	return c.DataDir
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

func (c *Config) GetMaxFileSize() int64 {
	return c.MaxFileSize
}

func (c *Config) GetDefaultExpiry() time.Duration {
	duration, _ := time.ParseDuration(c.DefaultExpiry)
	return duration
}
