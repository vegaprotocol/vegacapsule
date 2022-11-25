package generator

import (
	"context"
	"fmt"
	"log"
	"os"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/datanode"
	"code.vegaprotocol.io/vegacapsule/generator/faucet"
	"code.vegaprotocol.io/vegacapsule/generator/genesis"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/generator/visor"
	"code.vegaprotocol.io/vegacapsule/generator/wallet"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type nodeSets struct {
	validators    []types.NodeSet
	nonValidators []types.NodeSet
}

func (ns nodeSets) GetAllByGroupName(groupName string) []types.NodeSet {
	var out []types.NodeSet

	for _, n := range append(ns.validators, ns.nonValidators...) {
		if n.GroupName == groupName {
			out = append(out, n)
		}
	}

	return out
}

type jobRunner interface {
	RunRawNomadJobs(ctx context.Context, rawJobs []string) ([]types.RawJobWithNomadJob, error)
	StopNetwork(ctx context.Context, jobs *types.NetworkJobs, nodesOnly bool) ([]string, error)
}

type Generator struct {
	conf          *config.Config
	tendermintGen *tendermint.ConfigGenerator
	vegaGen       *vega.ConfigGenerator
	dataNodeGen   *datanode.ConfigGenerator
	genesisGen    *genesis.Generator
	walletGen     *wallet.ConfigGenerator
	faucetGen     *faucet.ConfigGenerator
	visorGen      *visor.Generator
	jobRunner     jobRunner
}

func New(conf *config.Config, genServices types.GeneratedServices, jobRunner jobRunner) (*Generator, error) {
	tendermintGen, err := tendermint.NewConfigGenerator(conf, genServices.NodeSets.ToSlice())
	if err != nil {
		return nil, fmt.Errorf("failed to create new tendermint config generator: %w", err)
	}
	vegaGen, err := vega.NewConfigGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new vega config generator: %w", err)
	}
	genesisGen, err := genesis.NewGenerator(conf, *conf.Network.GenesisTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to create new genesis generator: %w", err)
	}
	dataNodeGen, err := datanode.NewConfigGenerator(conf, genServices.NodeSets.ToSlice())
	if err != nil {
		return nil, fmt.Errorf("failed to create new data node generator: %w", err)
	}
	walletGen, err := wallet.NewConfigGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new wallet generator: %w", err)
	}
	faucetGen, err := faucet.NewConfigGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new faucet generator: %w", err)
	}
	visorGen, err := visor.NewGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new visor generator: %w", err)
	}

	return &Generator{
		conf:          conf,
		tendermintGen: tendermintGen,
		vegaGen:       vegaGen,
		genesisGen:    genesisGen,
		dataNodeGen:   dataNodeGen,
		walletGen:     walletGen,
		faucetGen:     faucetGen,
		visorGen:      visorGen,
		jobRunner:     jobRunner,
	}, nil
}

func (g *Generator) configureNodeSets(nss *nodeSets, fc *types.Faucet) error {
	for _, nc := range g.conf.Network.Nodes {
		co, err := newConfigOverride(g, nc)
		if err != nil {
			return err
		}

		for _, ns := range nss.GetAllByGroupName(nc.Name) {
			if err := co.Overwrite(nc, ns, fc); err != nil {
				return err
			}
		}

	}

	return nil
}

func (g *Generator) vegaChainID() string {
	return g.conf.Network.Name
}

func (g *Generator) Generate() (genSvc *types.GeneratedServices, err error) {
	var fc *types.Faucet
	if g.conf.Network.Faucet != nil {
		initFaucet, err := g.initAndConfigureFaucet(g.conf.Network.Faucet)
		if err != nil {
			return nil, err
		}

		fc = initFaucet
	}

	ns, err := g.initiateNodeSets()
	if err != nil {
		return nil, err
	}

	defer func() {
		// Stop pre-generate jobs
		if err != nil {
			if err := g.stopNomadJobs(); err != nil {
				log.Printf("Failed to stop Nomad jobs: %s", err)
			}
		}
	}()

	if err := g.configureNodeSets(ns, fc); err != nil {
		return nil, err
	}

	if err := g.genesisGen.GenerateAndSave(utils.ToPoint(g.vegaChainID()), ns.validators, ns.nonValidators, g.tendermintGen.GenesisValidators()); err != nil {
		return nil, fmt.Errorf("failed to generate genesis: %w", err)
	}

	var wl *types.Wallet
	if g.conf.Network.Wallet != nil {
		initWallet, err := g.initAndConfigureWallet(g.conf.Network.Wallet, ns.validators, ns.nonValidators)
		if err != nil {
			return nil, err
		}

		wl = initWallet
	}

	return types.NewGeneratedServices(wl, fc, append(ns.validators, ns.nonValidators...)), nil
}

func (g *Generator) AddNodeSet(absoluteIndex, relativeIndex, groupIndex int, nc config.NodeConfig, ns types.NodeSet, fc *types.Faucet) (*types.NodeSet, error) {
	preGenJobs, err := g.startPreGenerateJobs(nc, absoluteIndex)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Stop pre-generate jobs
		if err != nil {
			if err := g.stopNomadJobs(); err != nil {
				log.Printf("Failed to stop Nomad jobs: %s", err)
			}
		}
	}()

	cnc, err := nc.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clode node config for %q: %w", nc.Name, err)
	}

	initNodeSet, err := g.initiateNodeSet(absoluteIndex, relativeIndex, groupIndex, *cnc)
	if err != nil {
		return nil, err
	}

	initNodeSet.PreGenerateJobs = preGenJobs

	co, err := newConfigOverride(g, *cnc)
	if err != nil {
		return nil, err
	}

	if err := co.Overwrite(*cnc, *initNodeSet, fc); err != nil {
		return nil, err
	}

	if err := utils.CopyFile(ns.Tendermint.GenesisFilePath, initNodeSet.Tendermint.GenesisFilePath); err != nil {
		return nil, fmt.Errorf("failed to copy genesis file: %w", err)
	}

	log.Printf("Added new node set with id %q", initNodeSet.Name)

	return initNodeSet, nil
}

func (g *Generator) RemoveNodeSet(ns types.NodeSet) error {
	if err := os.RemoveAll(ns.Tendermint.HomeDir); err != nil {
		return fmt.Errorf("failed to remove Tendermint directory %q: %w", ns.Tendermint.HomeDir, err)
	}

	if err := os.RemoveAll(ns.Vega.HomeDir); err != nil {
		return fmt.Errorf("failed to remove Vega directory %q: %w", ns.Vega.HomeDir, err)
	}

	if ns.Visor != nil {
		if err := os.RemoveAll(ns.Visor.HomeDir); err != nil {
			return fmt.Errorf("failed to remove Visor directory %q: %w", ns.Vega.HomeDir, err)
		}
	}

	if ns.DataNode != nil {
		if err := os.RemoveAll(ns.DataNode.HomeDir); err != nil {
			return fmt.Errorf("failed to remove DataNode directory %q: %w", ns.DataNode.HomeDir, err)
		}
	}

	return nil
}
