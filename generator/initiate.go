package generator

import (
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/nomad"
	"code.vegaprotocol.io/vegacapsule/generator/wallet"
	"code.vegaprotocol.io/vegacapsule/types"
)

func (g *Generator) initiateNodeSet(absoluteIndex, relativeIndex, groupIndex int, nc config.NodeConfig) (*types.NodeSet, error) {
	n, err := config.TemplateNodeConfig(config.NodeConfigTemplateContext{NodeNumber: absoluteIndex}, nc)
	if err != nil {
		return nil, fmt.Errorf("failed to execute node config templates for %s: %w", nc.Name, err)
	}

	initTNode, err := g.tendermintGen.Initiate(absoluteIndex, n.Mode, n.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate Tendermint node id %d for node set %s: %w", absoluteIndex, n.Name, err)
	}

	initVNode, err := g.vegaGen.Initiate(
		absoluteIndex,
		n.VegaBinary,
		n.Mode,
		initTNode.HomeDir,
		n.NodeWalletPass,
		n.VegaWalletPass,
		n.EthereumWalletPass,
		n.ClefWallet,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate Vega node id %d for node set %s: %w", absoluteIndex, n.Name, err)
	}

	var initDNode *types.DataNode
	if n.UseDataNode {
		node, err := g.dataNodeGen.Initiate(absoluteIndex, g.vegaChainID(), nc.VegaBinary)
		if err != nil {
			return nil, fmt.Errorf("failed to initiate Data node id %d for node set %s: %w", absoluteIndex, n.Name, err)
		}

		initDNode = node
	}

	var initVisor *types.Visor
	// if data node binary is defined it is assumed that data-node should be deployed
	if n.VisorBinary != "" {
		node, err := g.visorGen.Initiate(absoluteIndex, n.VisorBinary, initVNode, initDNode)
		if err != nil {
			return nil, fmt.Errorf("failed to initiate Visor id %d for node set %s: %w", absoluteIndex, n.Name, err)
		}

		initVisor = node
	}

	nodeSet := &types.NodeSet{
		GroupName:     n.Name,
		Index:         absoluteIndex,
		RelativeIndex: relativeIndex,
		GroupIndex:    groupIndex,
		Name:          fmt.Sprintf("%s-nodeset-%s-%d-%s", g.conf.Network.Name, n.Name, absoluteIndex, n.Mode),
		Mode:          n.Mode,
		Vega:          *initVNode,
		Tendermint:    *initTNode,
		DataNode:      initDNode,
		Visor:         initVisor,
		PreStartProbe: n.PreStartProbe,
	}

	if n.NomadJobTemplate != nil {
		nodeJob, err := nomad.GenerateNodeSetTemplate(*n.NomadJobTemplate, *nodeSet)
		if err != nil {
			return nil, err
		}

		rawJob := nodeJob.String()
		nodeSet.NomadJobRaw = &rawJob
	}

	return nodeSet, nil
}

func (g *Generator) initiateNodeSets() (*nodeSets, error) {
	var mut sync.Mutex
	validatorsSet := []types.NodeSet{}
	nonValidatorsSet := []types.NodeSet{}

	var eg errgroup.Group
	var absIndex int
	for groupIndex, n := range g.conf.Network.Nodes {
		for relativeIndex := 0; relativeIndex < n.Count; relativeIndex++ {
			nc, err := n.Clone()
			if err != nil {
				return nil, fmt.Errorf("failed to clone node config for %q: %w", n.Name, err)
			}
			absIndexC := absIndex
			groupIndexC := groupIndex
			relativeIndexC := relativeIndex

			eg.Go(func() error {
				preGenJobs, err := g.startPreGenerateJobs(*nc, absIndexC)
				if err != nil {
					return err
				}

				nodeSet, err := g.initiateNodeSet(absIndexC, relativeIndexC, groupIndexC, *nc)
				if err != nil {
					return err
				}

				nodeSet.PreGenerateJobs = preGenJobs

				mut.Lock()
				if nc.Mode == types.NodeModeValidator {
					validatorsSet = append(validatorsSet, *nodeSet)
				} else {
					nonValidatorsSet = append(nonValidatorsSet, *nodeSet)
				}
				mut.Unlock()

				return nil
			})

			absIndex++
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &nodeSets{
		validators:    validatorsSet,
		nonValidators: nonValidatorsSet,
	}, nil
}

func (g *Generator) initAndConfigureFaucet(conf *config.FaucetConfig) (*types.Faucet, error) {
	initFaucet, err := g.faucetGen.InitiateAndConfigure(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate faucet: %w", err)
	}

	return initFaucet, nil
}

func (g *Generator) initAndConfigureBinary(conf *config.BinaryConfig) (*types.Binary, error) {
	initBinary, err := g.binaryGen.InitiateAndConfigure(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate binary: %w", err)
	}

	return initBinary, nil
}

func (g *Generator) initAndConfigureWallet(conf *config.WalletConfig, validatorsSet, nonValidatorSet []types.NodeSet) (*types.Wallet, error) {
	walletConfTemplate, err := wallet.NewConfigTemplate(g.conf.Network.Wallet.Template)
	if err != nil {
		return nil, err
	}

	initWallet, err := g.walletGen.InitiateWithNetworkConfig(g.conf.Network.Wallet, validatorsSet, nonValidatorSet, walletConfTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate wallet: %w", err)
	}

	if conf.TokenPassphraseFile != nil {
		initWallet.TokenPassphrasePath = *conf.TokenPassphraseFile
	}

	return initWallet, nil
}
