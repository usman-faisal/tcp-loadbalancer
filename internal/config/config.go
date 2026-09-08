package config

import (
	"os"
	"usman-faisal/tcp-loadbalancer/internal/types"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port int `yaml:"port"`
	BackendList []string `yaml:"backend_list"`
	Algorithm   types.Algorithm  `yaml:"algorithm"`
}

func Load(path string) (*Config, error){
	data, err:=os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

