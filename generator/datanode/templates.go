package datanode

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
	"github.com/zannen/toml"
)

type ConfigTemplateContext struct {
	Prefix      string
	NodeHomeDir string
	NodeNumber  int
	NodeSet     types.NodeSet
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config for data node: %w", err)
	}

	return t, nil
}

func (dng ConfigGenerator) TemplateConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	templateCtx := ConfigTemplateContext{
		Prefix:      *dng.conf.Prefix,
		NodeNumber:  ns.Index,
		NodeHomeDir: dng.homeDir,
		NodeSet:     ns,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template for data node: %w", err)
	}

	return buff, nil
}

func (dng *ConfigGenerator) TemplateAndMergeConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	tempFileName := fmt.Sprintf("datanode-%s.config", ns.Name)

	f, err := os.CreateTemp("", tempFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary config file %q: %w", tempFileName, err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	buff, err := dng.TemplateConfig(ns, configTemplate)
	if err != nil {
		return nil, err
	}

	// Sometimes the DataNode field may be nil. Especiall when you want to template the data-node config
	// with merge for the validator node.
	if ns.DataNode == nil {
		return nil, fmt.Errorf("failed to merge and save data node configuration: data node is not initialized properly")
	}

	if err := dng.mergeAndSaveConfig(ns, buff, originalConfigFilePath(ns.DataNode.HomeDir), f.Name()); err != nil {
		return nil, err
	}

	buffOut := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffOut, f); err != nil {
		return nil, fmt.Errorf("failed to copy content of config file %q: %w", f.Name(), err)
	}

	return buffOut, nil
}

func (dng ConfigGenerator) OverwriteConfig(ns types.NodeSet, configTemplate *template.Template) error {
	buff, err := dng.TemplateConfig(ns, configTemplate)
	if err != nil {
		return err
	}

	configFilePath := ConfigFilePath(ns.DataNode.HomeDir)
	return dng.mergeAndSaveConfig(ns, buff, configFilePath, configFilePath)
}

func (dng ConfigGenerator) mergeAndSaveConfig(
	ns types.NodeSet,
	tmpldConf *bytes.Buffer,
	configPath string,
	saveConfigPath string,
) error {
	overrideConfig := map[string]interface{}{}

	if _, err := toml.DecodeReader(tmpldConf, &overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	config := map[string]interface{}{}
	if err := paths.ReadStructuredFile(configPath, &config); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", configPath, err)
	}

	if err := mergo.Map(&config, overrideConfig, mergo.WithOverride); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := paths.WriteStructuredFile(saveConfigPath, config); err != nil {
		return fmt.Errorf("failed to write configuration file for data node at %s: %w", saveConfigPath, err)
	}

	return nil
}
