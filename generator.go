package main

import (
	"fmt"
)

type nodeSet struct {
	Mode       string
	Vega       InitVegaNode
	Tendermint InitTendermintNode
}

type Generator struct {
	conf          *Config
	tendermintGen *TendermintConfigGenerator
	vegaGen       *VegaConfigGenerator
	genesisGen    *GenesisGenerator
}

func NewGenerator(conf *Config) (*Generator, error) {
	tendermintGen, err := NewTendermintConfigGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new tendermint config generator: %w", err)
	}
	vegaGen, err := NewVegaConfigGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new vega config generator: %w", err)
	}
	genesisGen, err := NewGenesisGenerator(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new genesis generator: %w", err)
	}

	return &Generator{
		conf:          conf,
		tendermintGen: tendermintGen,
		vegaGen:       vegaGen,
		genesisGen:    genesisGen,
	}, nil
}

func (g *Generator) Generate() ([]nodeSet, error) {
	validatorsSet := []nodeSet{}
	nonValidatorsSet := []nodeSet{}

	var index int
	// Init phase
	for _, n := range g.conf.Network.Nodes {
		for i := 0; i < n.Count; i++ {
			initTNode, err := g.tendermintGen.Initiate(index, n.Mode)
			if err != nil {
				return nil, fmt.Errorf("failed to initiate Tendermit node: %w", err)
			}

			initVNode, err := g.vegaGen.Initiate(index, n.Mode, initTNode.HomeDir, n.NodeWalletPass, n.VegaWalletPass, n.EthereumWalletPass)
			if err != nil {
				return nil, fmt.Errorf("failed to initiate Tendermit node: %w", err)
			}

			if n.Mode == NodeModeValidator {
				validatorsSet = append(validatorsSet, nodeSet{
					Mode:       n.Mode,
					Vega:       *initVNode,
					Tendermint: *initTNode,
				})
			} else {
				nonValidatorsSet = append(nonValidatorsSet, nodeSet{
					Mode:       n.Mode,
					Vega:       *initVNode,
					Tendermint: *initTNode,
				})
			}

			index++
		}
	}

	index = 0
	// Override phase
	for _, n := range g.conf.Network.Nodes {
		// TODO finish override of Vega as well
		tendermintTemplate, err := g.tendermintGen.NewConfigTemplate(n.Templates.Tendermint)
		if err != nil {
			return nil, err
		}

		for i := 0; i < n.Count; i++ {
			if tendermintTemplate != nil {
				if err := g.tendermintGen.OverwriteConfig(index, tendermintTemplate); err != nil {
					return nil, fmt.Errorf("failed to overwrite Tendermit config: %w", err)
				}
			}

			index++
		}
	}

	if err := g.genesisGen.Generate(validatorsSet); err != nil {
		return nil, fmt.Errorf("failed to generate genesis: %w", err)
	}

	return append(validatorsSet, nonValidatorsSet...), nil
}
