package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

func ParseConfig(conf []byte, outputDir string) (*Config, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, err
	}

	if outputDir != "" {
		config.OutputDir = &outputDir
	}

	if err := hclsimple.Decode("config.hcl", conf, nil, config); err != nil {
		return nil, fmt.Errorf("failed to load decode configuration: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return config, nil
}

func ParseConfigFile(filePath, outputDir string) (*Config, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, err
	}

	if outputDir != "" {
		config.OutputDir = &outputDir
	}

	if err := hclsimple.DecodeFile(filePath, nil, config); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return config, nil
}
