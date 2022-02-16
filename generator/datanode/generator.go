package datanode

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
	"github.com/zannen/toml"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	datanodeconfig "code.vegaprotocol.io/data-node/config"
)

type ConfigTemplateContext struct {
	Prefix      string
	NodeHomeDir string
	NodeNumber  int
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config for data node: %w", err)
	}

	return t, nil
}

type ConfigGenerator struct {
	conf    *config.Config
	homeDir string
}

func NewConfigGenerator(conf *config.Config) (*ConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(conf.OutputDir, conf.DataNodePrefix))
	if err != nil {
		return nil, err
	}

	return &ConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
	}, nil
}

func (dng *ConfigGenerator) Initiate(index int, dataNodeBinary string) (*types.DataNode, error) {
	nodeDir := dng.nodeDir(index)

	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return nil, err
	}

	b, err := utils.ExecuteBinary(dataNodeBinary, []string{"init", "-f", "--home", nodeDir}, nil)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, string(b))

	initNode := &types.DataNode{
		HomeDir:    nodeDir,
		BinaryPath: dataNodeBinary,
	}

	return initNode, nil
}

func (dng ConfigGenerator) OverwriteConfig(index int, configTemplate *template.Template) error {
	templateCtx := ConfigTemplateContext{
		Prefix:      dng.conf.Prefix,
		NodeNumber:  index,
		NodeHomeDir: dng.homeDir,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return fmt.Errorf("failed to execute template for data node: %w", err)
	}

	overrideConfig := datanodeconfig.Config{}

	if _, err := toml.DecodeReader(buff, &overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	configFilePath := dng.configFilePath(dng.nodeDir(index))

	defaultConfig := datanodeconfig.NewDefaultConfig()
	if err := paths.ReadStructuredFile(configFilePath, &defaultConfig); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", configFilePath, err)
	}

	if err := mergo.Merge(&overrideConfig, defaultConfig); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := paths.WriteStructuredFile(configFilePath, overrideConfig); err != nil {
		return fmt.Errorf("failed to write configuration file for data node at %s: %w", configFilePath, err)
	}

	return nil
}

func (dng ConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", dng.conf.NodeDirPrefix, i)
	return filepath.Join(dng.homeDir, nodeDirName)
}

func (dng ConfigGenerator) configFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "data-node", "config.toml")
}
