package tendermint

import (
	"bytes"
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
	tmp2p "github.com/tendermint/tendermint/p2p"
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
	NodeSet              types.NodeSet
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

func NewConfigGenerator(conf *config.Config, generatedNodeSets []types.NodeSet) (*ConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(*conf.OutputDir, *conf.TendermintNodePrefix))
	if err != nil {
		return nil, err
	}

	nodesIDs := make([]string, 0, len(generatedNodeSets))
	for _, tn := range generatedNodeSets {
		nodesIDs = append(nodesIDs, tn.Tendermint.NodeID)
	}

	return &ConfigGenerator{
		conf:          conf,
		homeDir:       homeDir,
		genValidators: []tmtypes.GenesisValidator{},
		nodeIDs:       nodesIDs,
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
		_ = os.RemoveAll(*tg.conf.OutputDir)
		return nil, err
	}

	if err := os.MkdirAll(filepath.Join(nodeDir, "data"), nodeDirPerm); err != nil {
		_ = os.RemoveAll(*tg.conf.OutputDir)
		return nil, err
	}

	args := []string{"tm", "init", mode, "--home", nodeDir}

	log.Printf("Initiating Tendermint node %q with: %s %v", mode, *tg.conf.VegaBinary, args)

	b, err := utils.ExecuteBinary(*tg.conf.VegaBinary, args, nil)
	if err != nil {
		return nil, err
	}
	log.Println(string(b))

	config := tmconfig.DefaultConfig()
	config.SetRoot(nodeDir)

	nodeKey, err := tmp2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to get node key: %w", err)
	}

	tg.nodeIDs = append(tg.nodeIDs, string(nodeKey.ID()))

	initNode := &types.TendermintNode{
		Name:            fmt.Sprintf("tendermint-%s-%d", mode, index),
		HomeDir:         nodeDir,
		NodeID:          string(nodeKey.ID()),
		GenesisFilePath: config.BaseConfig.GenesisFile(),
		BinaryPath:      *tg.conf.VegaBinary,
	}

	if mode != string(types.NodeModeValidator) {
		return initNode, nil
	}

	pv := privval.LoadFilePV(config.BaseConfig.PrivValidatorKeyFile(), config.BaseConfig.PrivValidatorStateFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load FilePV for tendermint node: %w", err)
	}

	pubKey, err := pv.GetPubKey()
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

func (tg *ConfigGenerator) OverwriteConfig(ns types.NodeSet, configTemplate *template.Template) error {
	nodeDir := tg.nodeDir(ns.Index)
	configFilePath := tg.configFilePath(nodeDir)

	templateCtx := ConfigTemplateContext{
		Prefix:               *tg.conf.Prefix,
		TendermintNodePrefix: *tg.conf.TendermintNodePrefix,
		VegaNodePrefix:       *tg.conf.VegaNodePrefix,
		NodeNumber:           ns.Index,
		NodesCount:           len(tg.nodeIDs),
		NodeIDs:              tg.nodeIDs,
		NodePeers:            tg.getNodePeers(ns.Index),
		NodeSet:              ns,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return tg.mergeAndSaveConfig(nodeDir, configFilePath, buff)
}

func (tg ConfigGenerator) mergeAndSaveConfig(nodeDir, configFilePath string, merge *bytes.Buffer) error {
	v := viper.New()
	v.SetConfigFile(configFilePath)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %q: %w", configFilePath, err)
	}

	if err := v.MergeConfig(merge); err != nil {
		return fmt.Errorf("failed to merge config override with config file %q: %w", configFilePath, err)
	}
	conf := tmconfig.DefaultConfig()
	if err := v.Unmarshal(conf); err != nil {
		return fmt.Errorf("failed to unmarshal merged config file %q: %w", configFilePath, err)
	}

	if err := conf.ValidateBasic(); err != nil {
		return fmt.Errorf("failed to validated merged config file %q: %w", configFilePath, err)
	}

	conf.SetRoot(nodeDir)
	tmconfig.WriteConfigFile(configFilePath, conf)

	return nil
}

func (tg ConfigGenerator) GenesisValidators() []tmtypes.GenesisValidator {
	return tg.genValidators
}

func (tg ConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", *tg.conf.NodeDirPrefix, i)
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
