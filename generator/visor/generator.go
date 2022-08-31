package visor

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

const (
	GenesisFolderName        = "genesis"
	DefaultUpgradeFolderName = "vX.X.X"
	runConfigFileName        = "run-config.toml"
	configFileName           = "config.toml"
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

	initNode := &types.Visor{
		Name:       fmt.Sprintf("visor-%d", index),
		HomeDir:    visorDir,
		BinaryPath: visorBinary,
	}

	return initNode, nil
}

func (g Generator) PrepareUpgrade(
	index int,
	releaseTag string,
	ns types.NodeSet,
	configTemplate *template.Template,
) error {
	visorDir := g.visorDir(index)

	upgradeFolderName := filepath.Join(visorDir, releaseTag)

	log.Printf("Preparing upgrade folder %q for visor %q", upgradeFolderName, visorDir)

	if err := os.Mkdir(upgradeFolderName, os.ModePerm); err != nil {
		return err
	}

	upgradeRunConfigPath := filepath.Join(releaseTag, runConfigFileName)
	if err := utils.CopyFile(defaultUpgradeRunConfigFilePath(visorDir), upgradeRunConfigPath); err != nil {
		return err
	}

	log.Printf("Overwriting upgrade run config %q", upgradeRunConfigPath)

	// TODO TemplateAndMergeConfig use template instead so the result can be printed to stdout
	if err := g.OverwriteRunConfig(ns, configTemplate, upgradeRunConfigPath); err != nil {
		return err
	}

	return nil
}

func (g Generator) visorDir(i int) string {
	nodeDirName := fmt.Sprintf("%s%d", *g.conf.NodeDirPrefix, i)
	return filepath.Join(g.homeDir, nodeDirName)
}

func genesisFolder(nodeDir string) string {
	return filepath.Join(nodeDir, GenesisFolderName)
}

func configFilePath(nodeDir string) string {
	return filepath.Join(genesisFolder(nodeDir), configFileName)
}

func genesisRunConfigFilePath(nodeDir string) string {
	return filepath.Join(genesisFolder(nodeDir), runConfigFileName)
}

func defaultUpgradeFolder(nodeDir string) string {
	return filepath.Join(nodeDir, DefaultUpgradeFolderName)
}

func defaultUpgradeRunConfigFilePath(nodeDir string) string {
	return filepath.Join(defaultUpgradeFolder(nodeDir), runConfigFileName)
}
