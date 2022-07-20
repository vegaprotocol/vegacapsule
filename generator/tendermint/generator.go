package tendermint

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

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

type node struct {
	name      string
	groupName string
	id        string
	index     int
}

type ConfigGenerator struct {
	conf    *config.Config
	homeDir string

	genValidators []tmtypes.GenesisValidator
	nodes         []node
}

func newGenValidator(nodeDir string, config *tmconfig.Config) (*tmtypes.GenesisValidator, error) {

	pv, err := privval.LoadFilePV(config.PrivValidator.KeyFile(), config.PrivValidator.StateFile())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pubKey, err := pv.GetPubKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pubkey: %w", err)
	}

	return &tmtypes.GenesisValidator{
		Address: pubKey.Address(),
		PubKey:  pubKey,
		Power:   1,
		Name:    nodeDir,
	}, nil
}

func NewConfigGenerator(conf *config.Config, generatedNodeSets []types.NodeSet) (*ConfigGenerator, error) {
	homeDir, err := conf.StandarizePath(path.Join(*conf.OutputDir, *conf.TendermintNodePrefix))
	if err != nil {
		return nil, err
	}

	nodes := make([]node, 0, len(generatedNodeSets))
	genValidators := make([]tmtypes.GenesisValidator, 0, len(generatedNodeSets))
	for _, tn := range generatedNodeSets {
		nodes = append(nodes, node{
			name:      tn.Tendermint.Name,
			groupName: tn.GroupName,
			id:        tn.Tendermint.NodeID,
			index:     tn.Index,
		})

		if tn.Mode != types.NodeModeValidator {
			continue
		}
		config := tmconfig.DefaultConfig()
		config.SetRoot(tn.Tendermint.HomeDir)
		genValidator, err := newGenValidator(tn.Tendermint.HomeDir, config)
		if err != nil {
			return nil, err
		}

		genValidators = append(genValidators, *genValidator)
	}

	return &ConfigGenerator{
		conf:          conf,
		homeDir:       homeDir,
		nodes:         nodes,
		genValidators: genValidators,
	}, nil
}

func (tg ConfigGenerator) HomeDir() string {
	return tg.homeDir
}

func (tg *ConfigGenerator) Initiate(index int, mode, groupName string) (*types.TendermintNode, error) {
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

	confFilePath := ConfigFilePath(nodeDir)
	origConFilePath := originalConfigFilePath(nodeDir)

	if err := utils.CopyFile(confFilePath, origConFilePath); err != nil {
		return nil, fmt.Errorf("failed to copy initiated config from %q to %q: %w", confFilePath, origConFilePath, err)
	}

	config := tmconfig.DefaultConfig()
	config.SetRoot(nodeDir)

	nodeKey, err := tmtypes.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to get node key: %w", err)
	}

	nodeID := string(nodeKey.ID)
	nodeName := fmt.Sprintf("tendermint-%s-%d", mode, index)

	tg.nodes = append(tg.nodes, node{
		name:      nodeName,
		groupName: groupName,
		id:        nodeID,
		index:     index,
	})

	initNode := &types.TendermintNode{
		Name:            nodeName,
		HomeDir:         nodeDir,
		NodeID:          nodeID,
		GenesisFilePath: config.BaseConfig.GenesisFile(),
		BinaryPath:      *tg.conf.VegaBinary,
	}

	if mode != string(types.NodeModeValidator) {
		return initNode, nil
	}

	genValidator, err := newGenValidator(nodeDir, config)
	if err != nil {
		return nil, err
	}

	tg.genValidators = append(tg.genValidators, *genValidator)

	return initNode, nil
}

func (tg ConfigGenerator) GenesisValidators() []tmtypes.GenesisValidator {
	return tg.genValidators
}

func (tg ConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", *tg.conf.NodeDirPrefix, i)
	return filepath.Join(tg.homeDir, nodeDirName)
}

func ConfigFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "config.toml")
}

func originalConfigFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "config-original.toml")
}

func NodeKeyFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "node_key.json")
}

func PrivValidatorFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "priv_validator_key.json")
}
