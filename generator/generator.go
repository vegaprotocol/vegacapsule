package generator

import (
	"fmt"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/datanode"
	"code.vegaprotocol.io/vegacapsule/generator/faucet"
	"code.vegaprotocol.io/vegacapsule/generator/genesis"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/generator/wallet"
	"code.vegaprotocol.io/vegacapsule/types"
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

func New(conf *config.Config) (*Generator, error) {
	tendermintGen, err := tendermint.NewConfigGenerator(conf)
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

func (g *Generator) initiateNodeSet() (*nodeSets, error) {
	validatorsSet := []types.NodeSet{}
	nonValidatorsSet := []types.NodeSet{}

	var index int
	for _, n := range g.conf.Network.Nodes {
		for i := 0; i < n.Count; i++ {
			initTNode, err := g.tendermintGen.Initiate(index, n.Mode)
			if err != nil {
				return nil, fmt.Errorf("failed to initiate Tendermit node id %d: %w", index, err)
			}

			initVNode, err := g.vegaGen.Initiate(index, n.Mode, initTNode.HomeDir, n.NodeWalletPass, n.VegaWalletPass, n.EthereumWalletPass)
			if err != nil {
				return nil, fmt.Errorf("failed to initiate Vega node id %d: %w", index, err)
			}

			var initDNode *types.DataNode
			// if data node binary is defined it is assumed that data-node should be deployed
			if n.DataNodeBinary != "" {
				node, err := g.dataNodeGen.Initiate(index, n.DataNodeBinary)
				if err != nil {
					return nil, fmt.Errorf("failed to initiate Data node id %d: %w", index, err)
				}

				initDNode = node
			}

			if n.Mode == types.NodeModeValidator {
				validatorsSet = append(validatorsSet, types.NodeSet{
					Mode:       n.Mode,
					Vega:       *initVNode,
					Tendermint: *initTNode,
					DataNode:   initDNode,
				})
			} else {
				nonValidatorsSet = append(nonValidatorsSet, types.NodeSet{
					Mode:       n.Mode,
					Vega:       *initVNode,
					Tendermint: *initTNode,
					DataNode:   initDNode,
				})
			}

			index++
		}
	}

	return &nodeSets{
		validators:    validatorsSet,
		nonValidators: nonValidatorsSet,
	}, nil
}

func (g *Generator) configureNodeSets(fc *types.Faucet) error {
	var index int
	for _, n := range g.conf.Network.Nodes {
		tendermintConfTemplate, err := tendermint.NewConfigTemplate(n.Templates.Tendermint)
		if err != nil {
			return err
		}

		vegaConfTemplate, err := vega.NewConfigTemplate(n.Templates.Vega)
		if err != nil {
			return err
		}

		var dataNodeConfTemplate *template.Template
		if n.DataNodeBinary != "" {
			dataNodeConfTemplate, err = datanode.NewConfigTemplate(n.Templates.DataNode)

			if err != nil {
				return err
			}
		}

		for i := 0; i < n.Count; i++ {
			if tendermintConfTemplate != nil {
				if err := g.tendermintGen.OverwriteConfig(index, tendermintConfTemplate); err != nil {
					return fmt.Errorf("failed to overwrite Tendermit config for id %d: %w", index, err)
				}
			}
			if vegaConfTemplate != nil {
				if err := g.vegaGen.OverwriteConfig(index, n.Mode, fc, vegaConfTemplate); err != nil {
					return fmt.Errorf("failed to overwrite Vega config for id %d: %w", index, err)
				}
			}
			if dataNodeConfTemplate != nil {
				if err := g.dataNodeGen.OverwriteConfig(index, dataNodeConfTemplate); err != nil {
					return fmt.Errorf("failed to overwrite Data Node config for id %d: %w", index, err)
				}
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

	ns, err := g.initiateNodeSet()
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

	if err := g.persistNodesWalletsInfo(g.conf.OutputDir, ns.validators, ns.nonValidators); err != nil {
		return nil, fmt.Errorf("failed to write wallets info: %w", err)
	}

	return &types.GeneratedServices{
		Wallet:   wl,
		Faucet:   fc,
		NodeSets: append(ns.validators, ns.nonValidators...),
	}, nil
}

func (g *Generator) initAndConfigureFaucet(conf *config.FaucetConfig) (*types.Faucet, error) {
	initFaucet, err := g.faucetGen.InitiateAndConfigure(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to initate faucet: %w", err)
	}

	return initFaucet, nil
}

func (g *Generator) initAndConfigureWallet(conf *config.WalletConfig, validatorsSet []types.NodeSet) (*types.Wallet, error) {
	walletConfTemplate, err := wallet.NewConfigTemplate(g.conf.Network.Wallet.Template)
	if err != nil {
		return nil, err
	}

	initWallet, err := g.walletGen.InitiateWithNetworkConfig(g.conf.Network.Wallet, validatorsSet, walletConfTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to initate wallet: %w", err)
	}

	return initWallet, nil
}
