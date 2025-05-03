package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Qwen struct {
		APIKey  string `yaml:"api_key"`
		BaseURL string `yaml:"base_url"`
	} `yaml:"qwen"`

	Liblibai struct {
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		APIURL    string `yaml:"api_url"`
		QueryURL  string `yaml:"query_url"`
	} `yaml:"liblibai"`

	Video struct {
		APIURL  string `yaml:"api_url"`
		TaskURL string `yaml:"task_url"`
	} `yaml:"video"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
