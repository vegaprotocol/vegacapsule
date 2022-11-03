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

func (g Generator) TemplateConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	templateCtx := ConfigTemplateContext{
		NodeSet: ns,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buff, nil
}

// TODO solve this for other then genesis config
// TemplateAndMergeConfig templates provided template and merge it with originally initated Visor genesis run config
func (vg *Generator) TemplateAndMergeConfig(ns types.NodeSet, configTemplate *template.Template) (*bytes.Buffer, error) {
	buff, err := vg.TemplateConfig(ns, configTemplate)
	if err != nil {
		return nil, err
	}

	genRunConfigFilePath := genesisRunConfigFilePath(ns.Visor.HomeDir)

	if err := mergeAndSaveConfig(ns, buff, genRunConfigFilePath, vsconfig.RunConfig{}, vsconfig.RunConfig{}); err != nil {
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

// OverwriteConfigs overwrites visor config and genesis run config
func (g Generator) OverwriteConfigs(ns types.NodeSet, visorConfTemplate, runConfTemplate *template.Template) error {
	if err := g.OverwriteConfig(ns, visorConfTemplate); err != nil {
		return fmt.Errorf("failed to overwrite visor config: %w", err)
	}

	if err := g.OverwriteRunConfig(ns, runConfTemplate, genesisRunConfigFilePath(ns.Visor.HomeDir)); err != nil {
		return fmt.Errorf("failed to overwrite visor genesis run config: %w", err)
	}

	return nil
}

func (g Generator) OverwriteConfig(ns types.NodeSet, configTemplate *template.Template) error {
	buff, err := g.TemplateConfig(ns, configTemplate)
	if err != nil {
		return err
	}

	configPath := configFilePath(ns.Visor.HomeDir)
	return mergeAndSaveConfig(ns, buff, configPath, vsconfig.VisorConfigFile{}, vsconfig.VisorConfigFile{})
}

// OverwriteRunConfig overwrites run config with template in a given path.
// Uses default genesis path if not given.
func (g Generator) OverwriteRunConfig(ns types.NodeSet, configTemplate *template.Template, configPath string) error {
	buff, err := g.TemplateConfig(ns, configTemplate)
	if err != nil {
		return err
	}

	if configPath == "" {
		configPath = genesisRunConfigFilePath(ns.Visor.HomeDir)
	}

	return mergeAndSaveConfig(ns, buff, configPath, vsconfig.RunConfig{}, vsconfig.RunConfig{})
}

func mergeAndSaveConfig[T vsconfig.RunConfig | vsconfig.VisorConfigFile](
	ns types.NodeSet,
	tmpldConf *bytes.Buffer,
	configPath string,
	overrideConfig, originalConfig T,
) error {
	if _, err := toml.DecodeReader(tmpldConf, &overrideConfig); err != nil {
		fmt.Println(tmpldConf)
		return fmt.Errorf("failed decode override config: %w", err)
	}

	if err := paths.ReadStructuredFile(configPath, &originalConfig); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", configPath, err)
	}

	if err := mergo.MergeWithOverwrite(&originalConfig, overrideConfig); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := paths.WriteStructuredFile(configPath, originalConfig); err != nil {
		return fmt.Errorf("failed to write configuration file at %s: %w", configPath, err)
	}

	return nil
}
