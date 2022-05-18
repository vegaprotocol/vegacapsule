package vega

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"code.vegaprotocol.io/shared/paths"
	vgconfig "code.vegaprotocol.io/vega/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
	"github.com/zannen/toml"
)

type ConfigTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	DataNodePrefix       string
	ETHEndpoint          string
	NodeMode             string
	FaucetPublicKey      string
	NodeNumber           int
	NodeSet              types.NodeSet
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	return t, nil
}

func (vg ConfigGenerator) TemplateConfig(ns types.NodeSet, fc *types.Faucet, configTemplate *template.Template) (*bytes.Buffer, error) {
	templateCtx := ConfigTemplateContext{
		Prefix:               *vg.conf.Prefix,
		TendermintNodePrefix: *vg.conf.TendermintNodePrefix,
		VegaNodePrefix:       *vg.conf.VegaNodePrefix,
		DataNodePrefix:       *vg.conf.DataNodePrefix,
		ETHEndpoint:          vg.conf.Network.Ethereum.Endpoint,
		NodeMode:             ns.Mode,
		NodeNumber:           ns.Index,
		NodeSet:              ns,
	}

	if fc != nil {
		templateCtx.FaucetPublicKey = fc.PublicKey
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buff, nil
}

// TemplateAndMergeConfig templates provided template and merge it with originally initated Tendermint instance's config
func (vg *ConfigGenerator) TemplateAndMergeConfig(ns types.NodeSet, fc *types.Faucet, configTemplate *template.Template) (*bytes.Buffer, error) {
	tempFileName := fmt.Sprintf("vega-%s.config", ns.Name)

	f, err := os.CreateTemp("", tempFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary config file %q: %w", tempFileName, err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	buff, err := vg.TemplateConfig(ns, fc, configTemplate)
	if err != nil {
		return nil, err
	}

	if err := vg.mergeAndSaveConfig(ns, buff, vg.originalConfigFilePath(ns.Vega.HomeDir), f.Name()); err != nil {
		return nil, err
	}

	buffOut := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffOut, f); err != nil {
		return nil, fmt.Errorf("failed to copy content of config file %q: %w", f.Name(), err)
	}

	return buffOut, nil
}

func (vg ConfigGenerator) OverwriteConfig(ns types.NodeSet, fc *types.Faucet, configTemplate *template.Template) error {
	buff, err := vg.TemplateConfig(ns, fc, configTemplate)
	if err != nil {
		return err
	}

	configFilePath := vg.configFilePath(ns.Vega.HomeDir)
	return vg.mergeAndSaveConfig(ns, buff, configFilePath, configFilePath)
}

func (vg ConfigGenerator) mergeAndSaveConfig(
	ns types.NodeSet,
	tmpldConf *bytes.Buffer,
	configPath string,
	saveConfigPath string,
) error {
	overrideConfig := vgconfig.Config{}

	if _, err := toml.DecodeReader(tmpldConf, &overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	vegaConfig := vgconfig.Config{}
	if err := paths.ReadStructuredFile(configPath, &vegaConfig); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", configPath, err)
	}

	if err := mergo.MergeWithOverwrite(&vegaConfig, overrideConfig); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := paths.WriteStructuredFile(saveConfigPath, vegaConfig); err != nil {
		return fmt.Errorf("failed to write configuration file at %s: %w", saveConfigPath, err)
	}

	return nil
}
