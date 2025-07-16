package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3" // Импорт YAML пакета
)

type Configuration struct {
	Port        int    `yaml:"port" env:"PORT"`
	PathProxy   string `yaml:"pathProxy" env:"PATH_PROXY" envdefault:"https://jsonplaceholder.typicode.com/posts/"`
	LoggerLevel string `yaml:"loggerLevel" env:"LOGGER_LEVEL"`
}

func Load(path string) (*Configuration, error) {
	config := &Configuration{}

	if path != "" {
		configFile, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}

		if err := yaml.Unmarshal(configFile, config); err != nil {
			return nil, fmt.Errorf("error parsing YAML: %w", err)
		}
	}

	loadFromEnv(config)

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func loadFromEnv(cfg *Configuration) {

	if envPort := os.Getenv("PORT"); envPort != "" {
		if port, err := strconv.Atoi(envPort); err == nil {
			cfg.Port = port
		}
	}

	if envProxy := os.Getenv("PATH_PROXY"); envProxy != "" {
		cfg.PathProxy = envProxy
	} else if cfg.PathProxy == "" {

		cfg.PathProxy = "https://jsonplaceholder.typicode.com/posts/"
	}

	if envLogLevel := os.Getenv("LOGGER_LEVEL"); envLogLevel != "" {
		cfg.LoggerLevel = envLogLevel
	}
}

func (cfg *Configuration) Validate() error {
	if cfg.Port <= 0 {
		return errors.New("port must be positive")
	}
	if cfg.PathProxy == "" {
		return errors.New("proxy path is required")
	}
	return nil
}
