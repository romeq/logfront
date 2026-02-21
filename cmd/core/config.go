package main

import (
	"fmt"
	"os"

	"github.com/romeq/logfront/internal/sources"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Sources  map[string]sources.SourceConfigType `yaml:"sources"`
	Services map[string]sources.SourceConfigType `yaml:"services"`
}

func parseConfig(location string) (AppConfig, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return AppConfig{}, fmt.Errorf("error reading config file: %w", err)
	}

	var m AppConfig
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return AppConfig{}, err
	}

	return m, nil
}
