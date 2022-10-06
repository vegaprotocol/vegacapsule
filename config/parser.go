package config

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

func ParseConfigFile(filePath, outputDir string) (*Config, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, err
	}

	if outputDir != "" {
		config.OutputDir = &outputDir
	}

	if err := hclsimple.DecodeFile(filePath, defaultEvalContext, config); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	dir, _ := filepath.Split(filePath)
	if err := config.Validate(dir); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return config, nil
}
