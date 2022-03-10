package generator

import (
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/datanode"
	"code.vegaprotocol.io/vegacapsule/generator/faucet"
	"code.vegaprotocol.io/vegacapsule/generator/genesis"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/generator/wallet"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type nodeSets struct {
	validators    []types.NodeSet
	nonValidators []types.NodeSet
}

type Generator struct {
	conf          *config.Config
	tendermintGen *tendermint.ConfigGenerator
	vegaGen       *vega.ConfigGenerator
	dataNodeGen   *datanode.ConfigGenerator
	genesisGen    *genesis.Generator
	walletGen     *wallet.ConfigGenerator
	faucetGen     *faucet.ConfigGenerator
}

func New(conf *config.Config, genServices types.GeneratedServices) (*Generator, error) {
	tendermintGen, err := tendermint.NewConfigGenerator(conf, genServices.NodeSets)
	if err != nil {
		return nil, fmt.Errorf("failed to create new tendermint config generator: %w", err)
	}
	vegaGen, err := vega.NewConfigGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new vega config generator: %w", err)
	}
	genesisGen, err := genesis.NewGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new genesis generator: %w", err)
	}
	dataNodeGen, err := datanode.NewConfigGenerator(conf)
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

	return &Generator{
		conf:          conf,
		tendermintGen: tendermintGen,
		vegaGen:       vegaGen,
		genesisGen:    genesisGen,
		dataNodeGen:   dataNodeGen,
		walletGen:     walletGen,
		faucetGen:     faucetGen,
	}, nil
}

func (g *Generator) configureNodeSets(fc *types.Faucet) error {
	var index int
	for _, n := range g.conf.Network.Nodes {
		co, err := newConfigOverride(g, n)
		if err != nil {
			return err
		}

		for i := 0; i < n.Count; i++ {
			if err := co.Overwrite(index, n, fc); err != nil {
				return err
			}

			index++
		}
	}

	return nil
}

func (g *Generator) Generate() (*types.GeneratedServices, error) {
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

	if err := g.configureNodeSets(fc); err != nil {
		return nil, err
	}

	if err := g.genesisGen.Generate(ns.validators, ns.nonValidators, g.tendermintGen.GenesisValidators()); err != nil {
		return nil, fmt.Errorf("failed to generate genesis: %w", err)
	}

	var wl *types.Wallet
	if g.conf.Network.Wallet != nil {
		initWallet, err := g.initAndConfigureWallet(g.conf.Network.Wallet, ns.validators)
		if err != nil {
			return nil, err
		}

		wl = initWallet
	}

	return &types.GeneratedServices{
		Wallet:   wl,
		Faucet:   fc,
		NodeSets: append(ns.validators, ns.nonValidators...),
	}, nil
}

func (g *Generator) AddNodeSet(index int, nc config.NodeConfig, ns types.NodeSet, fc *types.Faucet) (*types.NodeSet, error) {
	initNodeSet, err := g.initiateNodeSet(index, nc)
	if err != nil {
		return nil, err
	}

	co, err := newConfigOverride(g, nc)
	if err != nil {
		return nil, err
	}

	if err := co.Overwrite(index, nc, fc); err != nil {
		return nil, err
	}

	if err := utils.CopyFile(ns.Tendermint.GenesisFilePath, initNodeSet.Tendermint.GenesisFilePath); err != nil {
		return nil, fmt.Errorf("failed to copy genesis file: %w", err)
	}

	log.Printf("Added new node set with id %q", initNodeSet.Name)

	return initNodeSet, nil
}
