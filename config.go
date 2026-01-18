package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	APIKey    string   `json:"api_key" yaml:"api_key"`
	AccessKey string   `json:"access_key" yaml:"access_key"`
	Port      string   `json:"port" yaml:"port"`
	Shows     []string `json:"shows" yaml:"shows"`
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config

	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		err = yaml.NewDecoder(file).Decode(&config)
	case ".json":
		err = json.NewDecoder(file).Decode(&config)
	default:
		return nil, errors.New("unsupported config file format")
	}

	if err != nil {
		return nil, errors.New("failed to decode config")
	}

	// defaults
	if config.Port == "" {
		config.Port = "8000"
	}

	// validate
	if config.APIKey == "" {
		return nil, errors.New("api_key is required")
	}

	if len(config.Shows) == 0 {
		return nil, errors.New("at least one show must be specified")
	}

	return &config, nil
}
