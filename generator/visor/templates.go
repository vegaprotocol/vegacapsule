package visor

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	"code.vegaprotocol.io/shared/paths"
	vsconfig "code.vegaprotocol.io/vega/visor/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
	"github.com/zannen/toml"
)

type ConfigTemplateContext struct {
	NodeSet types.NodeSet
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("run-config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	return t, nil
}

func (vg Generator) TemplateConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	templateCtx := ConfigTemplateContext{
		NodeSet: ns,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buff, nil
}

// TemplateAndMergeConfig templates provided template and merge it with originally initated Visor genesis run config
func (vg *Generator) TemplateAndMergeConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	buff, err := vg.TemplateConfig(ns, configTemplate)
	if err != nil {
		return nil, err
	}

	genRunConfigFilePath := genesisRunConfigFilePath(ns.Visor.HomeDir)
	upRunConfigFilePath := upgradeRunConfigFilePath(ns.Visor.HomeDir)

	if err := vg.mergeAndSaveConfig(ns, buff, genRunConfigFilePath, upRunConfigFilePath); err != nil {
		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(genRunConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read save config file %q: %w", genRunConfigFilePath, err)
	}

	buffOut := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffOut, bytes.NewReader(fileBytes)); err != nil {
		return nil, fmt.Errorf("failed to copy content of config file %q: %w", genRunConfigFilePath, err)
	}

	return buffOut, nil
}

func (vg Generator) OverwriteConfig(ns types.NodeSet, configTemplate *template.Template) error {
	buff, err := vg.TemplateConfig(ns, configTemplate)
	if err != nil {
		return err
	}

	genRunConfigFilePath := genesisRunConfigFilePath(ns.Visor.HomeDir)
	upRunConfigFilePath := upgradeRunConfigFilePath(ns.Visor.HomeDir)

	return vg.mergeAndSaveConfig(ns, buff, genRunConfigFilePath, upRunConfigFilePath)
}

func (vg Generator) mergeAndSaveConfig(
	ns types.NodeSet,
	tmpldConf *bytes.Buffer,
	genesisConfigPath, upgradeConfigPath string,
) error {
	overrideConfig := vsconfig.RunConfig{}

	if _, err := toml.DecodeReader(tmpldConf, &overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	visorConfig := vsconfig.RunConfig{}
	if err := paths.ReadStructuredFile(genesisConfigPath, &visorConfig); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", genesisConfigPath, err)
	}

	if err := mergo.MergeWithOverwrite(&visorConfig, overrideConfig); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := paths.WriteStructuredFile(genesisConfigPath, visorConfig); err != nil {
		return fmt.Errorf("failed to write configuration file at %s: %w", genesisConfigPath, err)
	}

	if err := paths.WriteStructuredFile(upgradeConfigPath, visorConfig); err != nil {
		return fmt.Errorf("failed to write configuration file at %s: %w", upgradeConfigPath, err)
	}

	return nil
}
