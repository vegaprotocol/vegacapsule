package visor

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type Generator struct {
	conf    *config.Config
	homeDir string
}

func NewGenerator(conf *config.Config) (*Generator, error) {
	homeDir, err := filepath.Abs(path.Join(*conf.OutputDir, *conf.VisorPrefix))
	if err != nil {
		return nil, err
	}

	return &Generator{
		conf:    conf,
		homeDir: homeDir,
	}, nil
}

func (g *Generator) Initiate(
	index int,
	visorBinary string,
	vegaNode *types.VegaNode,
	tmNode *types.TendermintNode,
	dataNode *types.DataNode,
) (*types.Visor, error) {
	visorDir := g.visorDir(index)

	args := []string{"init", "--home", visorDir}

	log.Printf("Initiating visor with: %s %v", visorBinary, args)

	b, err := utils.ExecuteBinary(visorBinary, args, nil)
	if err != nil {
		return nil, err
	}
	log.Println(string(b))

	if err := utils.CopyFile(vegaNode.BinaryPath, path.Join(genesisFolder(visorDir), "vega")); err != nil {
		return nil, err
	}

	if err := utils.CopyFile(vegaNode.BinaryPath, path.Join(upgradeFolder(visorDir), "vega")); err != nil {
		return nil, err
	}

	if dataNode != nil {
		if err := utils.CopyFile(dataNode.BinaryPath, path.Join(genesisFolder(visorDir), "data-node")); err != nil {
			return nil, err
		}
		if err := utils.CopyFile(dataNode.BinaryPath, path.Join(upgradeFolder(visorDir), "data-node")); err != nil {
			return nil, err
		}
	}

	initNode := &types.Visor{
		Name:       fmt.Sprintf("visor-%d", index),
		HomeDir:    visorDir,
		BinaryPath: visorBinary,
	}

	return initNode, nil
}

func (g Generator) visorDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", *g.conf.NodeDirPrefix, i)
	return filepath.Join(g.homeDir, nodeDirName)
}

func genesisFolder(nodeDir string) string {
	return filepath.Join(nodeDir, "genesis")
}

func upgradeFolder(nodeDir string) string {
	return filepath.Join(nodeDir, "upgrade")
}

func genesisRunConfigFilePath(nodeDir string) string {
	return filepath.Join(genesisFolder(nodeDir), "run-config.toml")
}

func upgradeRunConfigFilePath(nodeDir string) string {
	return filepath.Join(upgradeFolder(nodeDir), "run-config.toml")
}
