package generator

import (
	"bytes"
	"fmt"
	"sync"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/nomad"
	"code.vegaprotocol.io/vegacapsule/generator/wallet"
	"code.vegaprotocol.io/vegacapsule/types"
	"golang.org/x/sync/errgroup"
)

func (g *Generator) initiateNodeSet(index int, nc config.NodeConfig) (*types.NodeSet, error) {
	n, err := templateNodeConfig(index, nc)
	if err != nil {
		return nil, fmt.Errorf("failed to execute node config templates for %s: %w", nc.Name, err)
	}

	initTNode, err := g.tendermintGen.Initiate(index, n.Mode, n.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate Tendermit node id %d for node set %s: %w", index, n.Name, err)
	}

	initVNode, err := g.vegaGen.Initiate(
		index,
		n.Mode,
		initTNode.HomeDir,
		n.NodeWalletPass,
		n.VegaWalletPass,
		n.EthereumWalletPass,
		n.ClefWallet,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate Vega node id %d for node set %s: %w", index, n.Name, err)
	}

	var initDNode *types.DataNode
	// if data node binary is defined it is assumed that data-node should be deployed
	if n.DataNodeBinary != "" {
		node, err := g.dataNodeGen.Initiate(index, n.DataNodeBinary)
		if err != nil {
			return nil, fmt.Errorf("failed to initiate Data node id %d for node set %s: %w", index, n.Name, err)
		}

		initDNode = node
	}

	nodeSet := &types.NodeSet{
		GroupName:  n.Name,
		Index:      index,
		Name:       fmt.Sprintf("%s-nodeset-%s-%d-%s", g.conf.Network.Name, n.Name, index, n.Mode),
		Mode:       n.Mode,
		Vega:       *initVNode,
		Tendermint: *initTNode,
		DataNode:   initDNode,
	}

	// If mapping given, we want to create command runner
	if n.RemoteCommandRunner != nil {
		nodeSet.RemoteCommandRunner = &types.CommandRunner{
			Name: fmt.Sprintf("%s-cmd-runner", nodeSet.Name),
			PathsMapping: types.NetworkPathsMapping{
				VegaBinary:     *n.RemoteCommandRunner.PathsMapping.VegaBinary,
				VegaHome:       *n.RemoteCommandRunner.PathsMapping.VegaHome,
				TendermintHome: *n.RemoteCommandRunner.PathsMapping.TendermintHome,

				DataNodeBinary: n.RemoteCommandRunner.PathsMapping.DataNodeBinary,
				DataNodeHome:   n.RemoteCommandRunner.PathsMapping.DataNodeHome,
			},
		}

		var (
			rawTemplate *bytes.Buffer
			err         error
		)
		if n.RemoteCommandRunner.Nomad.JobTemplate != nil {
			rawTemplate, err = nomad.GenerateNodeSetTemplate(*n.RemoteCommandRunner.Nomad.JobTemplate, nodeSet)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to generate nomad template for remote command runner: %w", err)
		}

		nodeSet.RemoteCommandRunner.NomadJobRaw = rawTemplate.String()
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
	var index int
	for _, n := range g.conf.Network.Nodes {
		for i := 0; i < n.Count; i++ {
			nc, err := n.Clone()
			if err != nil {
				return nil, fmt.Errorf("failed to clode node config for %q: %w", n.Name, err)
			}
			indexc := index

			eg.Go(func() error {
				preGenJobs, err := g.startPreGenerateJobs(*nc, indexc)
				if err != nil {
					return err
				}

				nodeSet, err := g.initiateNodeSet(indexc, *nc)
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

			index++
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
		return nil, fmt.Errorf("failed to initate faucet: %w", err)
	}

	return initFaucet, nil
}

func (g *Generator) initAndConfigureWallet(conf *config.WalletConfig, validatorsSet, nonValidatorSet []types.NodeSet) (*types.Wallet, error) {
	walletConfTemplate, err := wallet.NewConfigTemplate(g.conf.Network.Wallet.Template)
	if err != nil {
		return nil, err
	}

	initWallet, err := g.walletGen.InitiateWithNetworkConfig(g.conf.Network.Wallet, validatorsSet, nonValidatorSet, walletConfTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to initate wallet: %w", err)
	}

	return initWallet, nil
}
