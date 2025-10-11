package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server Server   `yaml:"server"`
	Routes []*Route `yaml:"routes"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Route struct {
	ID      string   `yaml:"id"`
	Path    string   `yaml:"path"`
	Targets []string `yaml:"targets"`
}

func LoadConfig(config string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
