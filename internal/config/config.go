package config

import (
	"flag"
	"io/ioutil"
	"os"
	"sync"

	"github.com/romarq/visualtez-testing/internal/logger"
	"gopkg.in/yaml.v2"
)

// Config holds API configurations
type Config struct {
	Port  string      `yaml:"port,omitempty"`
	Tezos TezosConfig `yaml:"tezos,omitempty"`
	Log   LogConfig   `yaml:"log,omitempty"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Location string          `yaml:"location,omitempty"`
	Level    logger.LogLevel `yaml:"level,omitempty"`
}

// TezosConfig holds tezos configurations
type TezosConfig struct {
	BaseDirectory   string `yaml:"dir,omitempty"`
	DefaultProtocol string `yaml:"default_protocol,omitempty"`
}

// EnvironmentProperty - Known environment properties
type EnvironmentProperty string

const (
	logLocation EnvironmentProperty = "LOG_LOCATION"
	apiPort     EnvironmentProperty = "API_PORT"
)

var once sync.Once
var singleton Config

// GetConfig - Get Configurations (Singleton pattern)
func GetConfig() Config {
	// Load environment variables only once
	once.Do(func() {
		var configPath string
		flag.StringVar(&configPath, "config", "./config.yaml", "API config file location")
		flag.Parse()
		singleton = load(configPath)
	})

	return singleton
}

// Load configuration from yaml and environment variables
func load(file string) Config {
	logger.Info("Loading configurations from: %s", file)

	// Config instance
	c := Config{}

	// Load config from YAML file
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Warn("Error reading configuration file: %s. %v", file, err)
	}
	if err := yaml.Unmarshal(fileContents, &c); err != nil {
		logger.Warn("Failed to parse configuration file: %s.\n %s", file, err)
	}

	// Override configurations by ENV values (if provided)

	logLocationFromEnv := os.Getenv(string(logLocation))
	if logLocationFromEnv != "" {
		c.Log.Location = logLocationFromEnv
	}

	apiPortFromEnv := os.Getenv(string(apiPort))
	if apiPortFromEnv != "" {
		c.Port = apiPortFromEnv
	}

	return c
}
