package config

import (
	"fmt"
	"os"

	"code.vegaprotocol.io/vegacapsule/utils"
)

type BinaryConfig struct {
	/*
		description: Name of the service that is going to be used as an identifier when service runs.
		example:
			type: hcl
			value: |
					binary_service "service-name" {
						...
					}
	*/
	Name string `hcl:"name,label"`

	/*
		description: Command that will run at the app startup.
		example:
			type: hcl
			value: |
					cmd = "serve"
	*/
	Command string `hcl:"cmd,optional"`

	/*
		description: List of arguments that will be added to cmd.
		example:
			type: hcl
			value: |
					args = [
						"--config", "--config=config/config_capsule.yaml",
					]
	*/
	Args []string `hcl:"args,optional"`

	/*
		description: |
			[Go template](templates.md) of a Binary Service config.

			The [binary.ConfigTemplateContext](templates.md#binaryconfigtemplatecontext) can be used in the template.
			Example can be found in [default network config](net_confs/config.hcl).
		examples:
			- type: hcl
			  value: |
						config_template = <<EOH
							...
						EOH

	*/

	/*
		description: |
					Allows user to define a Service binary to be used.
					A relative or absolute path can be used. If only the binary name is defined, it automatically looks for it in $PATH.
					This can help with testing different version compatibilities or a protocol upgrade.
		note: Using versions that are not compatible could break the network - therefore this should be used in advanced cases only.
	*/
	BinaryFile *string `hcl:"binary_path,optional"`

	/*
		description: |
			The [binary.ConfigTemplateContext](templates.md#binaryconfigtemplatecontext) can be used in the template.
			Example can be found in [default network config](net_confs/config.hcl).
		examples:
			- type: hcl
			  value: |
						template = <<EOH
							...
						EOH

	*/
	ConfigTemplate *string `hcl:"config_template,optional"`

	ConfigTemplateFile *string `hcl:"config_template_file,optional"`

	Sync bool `hcl:"sync,optional"`
}

func (b *BinaryConfig) GetConfigTemplate(configDir string) (*string, error) {
	if b.ConfigTemplate != nil {
		return b.ConfigTemplate, nil
	}

	if b.ConfigTemplateFile == nil {
		return nil, nil
	}

	templateFile, err := utils.AbsPathWithPrefix(configDir, *b.ConfigTemplateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute file path %q: %w", *b.ConfigTemplateFile, err)
	}

	template, err := os.ReadFile(templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", templateFile, err)
	}

	str := string(template)
	b.ConfigTemplate = &str

	return &str, nil
}
