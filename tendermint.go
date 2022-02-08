package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

const (
	nodeDirPerm = 0755
)

type InitTendermintNode struct {
	HomeDir         string
	GenesisFilePath string
}

type TendermingTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	NodeNumber           int
	NodesCount           int
	NodeIDs              []string
}

type TendermintConfigGenerator struct {
	conf    *Config
	homeDir string

	genValidators []types.GenesisValidator
	nodeIDs       []string
}

func NewTendermintConfigGenerator(conf *Config) (*TendermintConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(conf.OutputDir, conf.TendermintNodePrefix))
	if err != nil {
		return nil, err
	}

	return &TendermintConfigGenerator{
		conf:          conf,
		homeDir:       homeDir,
		genValidators: []types.GenesisValidator{},
		nodeIDs:       []string{},
	}, nil
}

func (tg TendermintConfigGenerator) NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	return t, nil
}

func (tg TendermintConfigGenerator) HomeDir() string {
	return tg.homeDir
}

func (tg *TendermintConfigGenerator) Initiate(index int, mode string) (*InitTendermintNode, error) {
	nodeDir := tg.nodeDir(index)

	err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
	if err != nil {
		_ = os.RemoveAll(tg.conf.OutputDir)
		return nil, err
	}
	err = os.MkdirAll(filepath.Join(nodeDir, "data"), nodeDirPerm)
	if err != nil {
		_ = os.RemoveAll(tg.conf.OutputDir)
		return nil, err
	}

	b, err := executeBinary(tg.conf.VegaBinary, []string{"tm", "init", "--home", nodeDir}, nil)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, string(b))

	config := cfg.DefaultConfig()
	config.SetRoot(nodeDir)

	if mode != string(NodeModeValidator) {
		return nil, nil
	}

	pv := privval.LoadFilePV(config.BaseConfig.PrivValidatorKeyFile(), config.BaseConfig.PrivValidatorStateFile())

	pubKey, err := pv.GetPubKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get pubkey: %w", err)
	}

	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to get node key: %w", err)
	}

	tg.nodeIDs = append(tg.nodeIDs, string(nodeKey.ID()))
	tg.genValidators = append(tg.genValidators, types.GenesisValidator{
		Address: pubKey.Address(),
		PubKey:  pubKey,
		Power:   1,
		Name:    nodeDir,
	})

	return &InitTendermintNode{
		HomeDir:         nodeDir,
		GenesisFilePath: config.BaseConfig.GenesisFile(),
	}, nil
}

func (tg *TendermintConfigGenerator) OverwriteConfig(index int, configTemplate *template.Template) error {
	nodeDir := tg.nodeDir(index)
	configFilePath := tg.configFilePath(nodeDir)

	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %q: %w", configFilePath, err)
	}

	templateCtx := TendermingTemplateContext{
		Prefix:               tg.conf.Prefix,
		TendermintNodePrefix: tg.conf.TendermintNodePrefix,
		VegaNodePrefix:       tg.conf.VegaNodePrefix,
		NodeNumber:           index,
		NodesCount:           len(tg.nodeIDs),
		NodeIDs:              tg.nodeIDs,
	}

	buff := bytes.NewBuffer([]byte{})

	err := configTemplate.Execute(buff, templateCtx)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := viper.MergeConfig(buff); err != nil {
		return fmt.Errorf("failed to merge config override with config file %q: %w", configFilePath, err)
	}

	config := &cfg.Config{}
	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal merged config file %q: %w", configFilePath, err)
	}
	if err := config.ValidateBasic(); err != nil {
		return fmt.Errorf("failed to validated merged config file %q: %w", configFilePath, err)
	}

	config.SetRoot(nodeDir)
	cfg.WriteConfigFile(configFilePath, config)

	return nil
}

func (tg TendermintConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", tg.conf.NodeDirPrefix, i)
	return filepath.Join(tg.homeDir, nodeDirName)
}

func (tg TendermintConfigGenerator) configFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "config.toml")
}
