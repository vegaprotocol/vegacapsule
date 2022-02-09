package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

func ParseConfig(conf []byte) (*Config, error) {
	config := &Config{}
	if err := hclsimple.Decode("config.hcl", conf, nil, config); err != nil {
		return nil, fmt.Errorf("failed to load decode configuration: %w", err)
	}

	return config, nil
}

func ParseConfigFile(filePath string) (*Config, error) {
	config := &Config{}
	if err := hclsimple.DecodeFile(filePath, nil, config); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return config, nil
}
