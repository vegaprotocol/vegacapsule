package datanode

import (
	"fmt"
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
	homeDir, err := filepath.Abs(path.Join(conf.OutputDir, conf.DataNodePrefix))
	if err != nil {
		return nil, err
	}

	return &ConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
	}, nil
}

func (tg *ConfigGenerator) Initiate(index int, dataNodeBinary string) (*types.DataNode, error) {
	nodeDir := tg.nodeDir(index)

	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return nil, err
	}

	b, err := utils.ExecuteBinary(dataNodeBinary, []string{"init", "-f", "--root-path", nodeDir}, nil)
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

func (tg ConfigGenerator) nodeDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", tg.conf.NodeDirPrefix, i)
	return filepath.Join(tg.homeDir, nodeDirName)
}

func (tg ConfigGenerator) configFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config.toml")
}
