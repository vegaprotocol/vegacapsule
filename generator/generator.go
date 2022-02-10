package generator

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/datanode"
	"code.vegaprotocol.io/vegacapsule/generator/genesis"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/types"
)

type Generator struct {
	conf          *config.Config
	tendermintGen *tendermint.ConfigGenerator
	vegaGen       *vega.ConfigGenerator
	dataNodeGen   *datanode.ConfigGenerator
	genesisGen    *genesis.Generator
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
	dataNode, err := datanode.NewConfigGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new genesis generator: %w", err)
	}

	return &Generator{
		conf:          conf,
		tendermintGen: tendermintGen,
		vegaGen:       vegaGen,
		genesisGen:    genesisGen,
		dataNodeGen:   dataNode,
	}, nil
}

func (g *Generator) Generate() ([]types.NodeSet, error) {
	validatorsSet := []types.NodeSet{}
	nonValidatorsSet := []types.NodeSet{}

	var index int
	// Init phase
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

	index = 0
	// Override phase
	for _, n := range g.conf.Network.Nodes {
		tendermintConfTemplate, err := tendermint.NewConfigTemplate(n.Templates.Tendermint)
		if err != nil {
			return nil, err
		}

		vegaConfTemplate, err := vega.NewConfigTemplate(n.Templates.Vega)
		if err != nil {
			return nil, err
		}

		for i := 0; i < n.Count; i++ {
			if tendermintConfTemplate != nil {
				if err := g.tendermintGen.OverwriteConfig(index, tendermintConfTemplate); err != nil {
					return nil, fmt.Errorf("failed to overwrite Tendermit config for id %d: %w", index, err)
				}
			}
			if vegaConfTemplate != nil {
				if err := g.vegaGen.OverwriteConfig(index, n.Mode, vegaConfTemplate); err != nil {
					return nil, fmt.Errorf("failed to overwrite Vega config for id %d: %w", index, err)
				}
			}

			index++
		}
	}

	if err := g.genesisGen.Generate(validatorsSet, nonValidatorsSet, g.tendermintGen.GenesisValidators()); err != nil {
		return nil, fmt.Errorf("failed to generate genesis: %w", err)
	}

	return append(validatorsSet, nonValidatorsSet...), nil
}
