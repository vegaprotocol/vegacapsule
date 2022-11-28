package datanode

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type node struct {
	name  string
	index int
}

type ConfigGenerator struct {
	conf    *config.Config
	homeDir string

	nodes []node
}

func NewConfigGenerator(conf *config.Config, generatedNodeSets []types.NodeSet) (*ConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(*conf.OutputDir, *conf.DataNodePrefix))
	if err != nil {
		return nil, err
	}

	nodes := []node{}
	for _, n := range generatedNodeSets {
		if n.DataNode == nil {
			continue
		}
		nodes = append(nodes, node{
			name:  n.DataNode.Name,
			index: n.Index,
		})
	}

	return &ConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
		nodes:   nodes,
	}, nil
}

func (dng *ConfigGenerator) Initiate(index int, chainID string, optVegaBinary *string) (*types.DataNode, error) {
	nodeDir := dng.nodeDir(index)
	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return nil, err
	}

	vegaBinary := *dng.conf.VegaBinary
	if optVegaBinary != nil {
		vegaBinary = *optVegaBinary
	}

	args := []string{
		config.DataNodeSubCmd, "init",
		"-f",
		"--home", nodeDir,
		chainID,
	}

	log.Printf("Initiating data node with: %s %v", vegaBinary, args)

	b, err := utils.ExecuteBinary(vegaBinary, args, nil)
	if err != nil {
		return nil, err
	}
	log.Println(string(b))

	confFilePath := ConfigFilePath(nodeDir)
	origConFilePath := originalConfigFilePath(nodeDir)

	if err := utils.CopyFile(confFilePath, origConFilePath); err != nil {
		return nil, fmt.Errorf("failed to copy initiated config from %q to %q: %w", confFilePath, origConFilePath, err)
	}

	initNode := &types.DataNode{
		GeneratedService: types.GeneratedService{
			Name:           fmt.Sprintf("data-node-%d", index),
			HomeDir:        nodeDir,
			ConfigFilePath: confFilePath,
		},
		BinaryPath: vegaBinary,
	}

	dng.nodes = append(dng.nodes, node{
		name:  initNode.Name,
		index: index,
	})

	return initNode, nil
}

func (dng ConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", *dng.conf.NodeDirPrefix, i)
	return filepath.Join(dng.homeDir, nodeDirName)
}

func ConfigFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "data-node", "config.toml")
}

func originalConfigFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "data-node", "original-config.toml")
}
