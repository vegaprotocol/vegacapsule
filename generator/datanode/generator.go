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

type ConfigGenerator struct {
	conf    *config.Config
	homeDir string
}

func NewConfigGenerator(conf *config.Config) (*ConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(*conf.OutputDir, *conf.DataNodePrefix))
	if err != nil {
		return nil, err
	}

	return &ConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
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
