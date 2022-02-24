package tendermint

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/spf13/viper"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/privval"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	nodeDirPerm = 0755
)

type Peer struct {
	Index int
	ID    string
}

type ConfigTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	NodeNumber           int
	NodesCount           int
	NodeIDs              []string
	NodePeers            []Peer
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	return t, nil
}

type ConfigGenerator struct {
	conf    *config.Config
	homeDir string

	genValidators []tmtypes.GenesisValidator
	nodeIDs       []string
}

func NewConfigGenerator(conf *config.Config) (*ConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(conf.OutputDir, conf.TendermintNodePrefix))
	if err != nil {
		return nil, err
	}

	return &ConfigGenerator{
		conf:          conf,
		homeDir:       homeDir,
		genValidators: []tmtypes.GenesisValidator{},
		nodeIDs:       []string{},
	}, nil
}

func (tg ConfigGenerator) HomeDir() string {
	return tg.homeDir
}

func (tg *ConfigGenerator) Initiate(index int, mode string) (*types.TendermintNode, error) {
	nodeDir := tg.nodeDir(index)

	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
		_ = os.RemoveAll(tg.conf.OutputDir)
		return nil, err
	}

	if err := os.MkdirAll(filepath.Join(nodeDir, "data"), nodeDirPerm); err != nil {
		_ = os.RemoveAll(tg.conf.OutputDir)
		return nil, err
	}

	b, err := utils.ExecuteBinary(tg.conf.VegaBinary, []string{"tm", "init", mode, "--home", nodeDir}, nil)
	if err != nil {
		return nil, err
	}
	log.Println(string(b))

	config := tmconfig.DefaultConfig()
	config.SetRoot(nodeDir)

	nodeKey, err := tmtypes.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to get node key: %w", err)
	}

	tg.nodeIDs = append(tg.nodeIDs, string(nodeKey.ID))

	initNode := &types.TendermintNode{
		HomeDir:         nodeDir,
		GenesisFilePath: config.BaseConfig.GenesisFile(),
	}

	if mode != string(types.NodeModeValidator) {
		return initNode, nil
	}

	pv, err := privval.LoadFilePV(config.PrivValidator.KeyFile(), config.PrivValidator.StateFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load FilePV for tendermint node: %w", err)
	}

	// TODO: Pass context from higher function to avoid locking
	pubKey, err := pv.GetPubKey(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to get pubkey: %w", err)
	}

	tg.genValidators = append(tg.genValidators, tmtypes.GenesisValidator{
		Address: pubKey.Address(),
		PubKey:  pubKey,
		Power:   1,
		Name:    nodeDir,
	})

	return initNode, nil
}

func (tg *ConfigGenerator) OverwriteConfig(index int, configTemplate *template.Template) error {
	nodeDir := tg.nodeDir(index)
	configFilePath := tg.configFilePath(nodeDir)

	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %q: %w", configFilePath, err)
	}

	templateCtx := ConfigTemplateContext{
		Prefix:               tg.conf.Prefix,
		TendermintNodePrefix: tg.conf.TendermintNodePrefix,
		VegaNodePrefix:       tg.conf.VegaNodePrefix,
		NodeNumber:           index,
		NodesCount:           len(tg.nodeIDs),
		NodeIDs:              tg.nodeIDs,
		NodePeers:            tg.getNodePeers(index),
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := viper.MergeConfig(buff); err != nil {
		return fmt.Errorf("failed to merge config override with config file %q: %w", configFilePath, err)
	}

	conf := &tmconfig.Config{}
	if err := viper.Unmarshal(conf); err != nil {
		return fmt.Errorf("failed to unmarshal merged config file %q: %w", configFilePath, err)
	}

	if err := conf.ValidateBasic(); err != nil {
		return fmt.Errorf("failed to validated merged config file %q: %w", configFilePath, err)
	}

	conf.SetRoot(nodeDir)
	if err := tmconfig.WriteConfigFile(nodeDir, conf); err != nil {
		return fmt.Errorf("failed to write overwritten tendermint configuration: %w", err)
	}

	return nil
}

func (tg ConfigGenerator) GenesisValidators() []tmtypes.GenesisValidator {
	return tg.genValidators
}

func (tg ConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", tg.conf.NodeDirPrefix, i)
	return filepath.Join(tg.homeDir, nodeDirName)
}

func (tg ConfigGenerator) configFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "config.toml")
}

func (tg ConfigGenerator) getNodePeers(index int) []Peer {
	peers := []Peer{}

	for nodeIdx, nodeId := range tg.nodeIDs {
		if nodeIdx == index {
			continue
		}

		peers = append(peers, Peer{
			Index: nodeIdx,
			ID:    nodeId,
		})
	}

	return peers
}
