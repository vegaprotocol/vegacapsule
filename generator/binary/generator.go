package binary

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
)

type ConfigTemplateContext struct {
	// description: Path to home directory of the Binary.
	HomeDir    string
	Validators []types.NodeSet
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the binary's template config: %w", err)
	}

	return t, nil
}

type ConfigGenerator struct {
	conf    *config.Config
	homeDir string
}

func NewConfigGenerator(conf *config.Config) (*ConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(*conf.OutputDir, "binary_services"))
	if err != nil {
		return nil, err
	}

	return &ConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
	}, nil
}

func (cg *ConfigGenerator) InitiateAndConfigure(conf *config.BinaryConfig) (*types.Binary, error) {
	initBinary, err := cg.Initiate(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to init binary for %q: %w", conf.Name, err)
	}

	if conf.ConfigTemplate == nil {
		return initBinary, nil
	}

	configTemplate, err := NewConfigTemplate(*conf.ConfigTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to create new template config for %q: %w", conf.Name, err)
	}

	if err := cg.OverwriteConfig(initBinary, configTemplate); err != nil {
		return nil, fmt.Errorf("failed to override config for %q: %w", conf.Name, err)
	}

	return initBinary, nil
}

func (cg *ConfigGenerator) Initiate(conf *config.BinaryConfig) (*types.Binary, error) {
	binHome := path.Join(cg.homeDir, conf.Name)
	if err := os.MkdirAll(binHome, os.ModePerm); err != nil {
		return nil, err
	}
	return &types.Binary{
		GeneratedService: types.GeneratedService{
			Name:           fmt.Sprintf("%s-%s", cg.conf.Network.Name, conf.Name),
			ConfigFilePath: path.Join(binHome, "config.toml"),
			HomeDir:        binHome,
		},
		Args:       conf.Args,
		BinaryPath: *conf.BinaryFile,
	}, nil
}

func (cg *ConfigGenerator) OverwriteConfig(bc *types.Binary, configTemplate *template.Template) error {
	templateCtx := ConfigTemplateContext{
		HomeDir: cg.homeDir,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return fmt.Errorf("failed to execute template for binary: %w", err)
	}

	confMap := make(map[string]interface{})

	if _, err := toml.NewDecoder(buff).Decode(&confMap); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	if err := paths.WriteStructuredFile(bc.ConfigFilePath, confMap); err != nil {
		return fmt.Errorf("failed to write configuration file for binary at %s: %w", bc.ConfigFilePath, err)
	}

	return nil
}
