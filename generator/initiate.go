package generator

import (
	"fmt"
	"sync"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/wallet"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
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
		cmdRunner, err := g.initRemoteCommandRunner(*n.RemoteCommandRunner, nodeSet)
		if err != nil {
			return nil, fmt.Errorf("failed initialising command runner: %w", err)
		}
		nodeSet.RemoteCommandRunner = cmdRunner
	}

	if n.NomadJobTemplate != nil {
		nodeJob, err := utils.GenerateTemplate(*n.NomadJobTemplate, *nodeSet)
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

func (g Generator) initRemoteCommandRunner(commandRunnerConfig config.CommandRunner, nodeSet *types.NodeSet) (*types.CommandRunner, error) {
	mErr := utils.NewMultiError()

	vegaBinaryPath, err := utils.GenerateTemplate(*commandRunnerConfig.PathsMapping.VegaBinary, *nodeSet)
	if err != nil {
		mErr.Add(fmt.Errorf("failed to template vega binary path in the command runner config: %w", err))
	}

	vegaHomePath, err := utils.GenerateTemplate(*commandRunnerConfig.PathsMapping.VegaHome, *nodeSet)
	if err != nil {
		mErr.Add(fmt.Errorf("failed to template vega home path in the command runner config: %w", err))
	}

	tendermintHomePath, err := utils.GenerateTemplate(*commandRunnerConfig.PathsMapping.TendermintHome, *nodeSet)
	if err != nil {
		mErr.Add(fmt.Errorf("failed to template tendermint home path in the command runner config: %w", err))
	}

	remoteCommandRunner := &types.CommandRunner{
		Name: fmt.Sprintf("%s-cmd-runner", nodeSet.Name),
		PathsMapping: types.NetworkPathsMapping{
			VegaBinary:     vegaBinaryPath.String(),
			VegaHome:       vegaHomePath.String(),
			TendermintHome: tendermintHomePath.String(),
		},
	}

	if commandRunnerConfig.PathsMapping.DataNodeBinary != nil {
		dataNodeBinaryPath, err := utils.GenerateTemplate(*commandRunnerConfig.PathsMapping.DataNodeBinary, *nodeSet)
		if err != nil {
			mErr.Add(fmt.Errorf("failed to template data-node binary path in the command runner config: %w", err))
		}
		remoteCommandRunner.PathsMapping.DataNodeBinary = utils.StrPoint(dataNodeBinaryPath.String())
	}

	if commandRunnerConfig.PathsMapping.DataNodeHome != nil {
		dataNodeHomePath, err := utils.GenerateTemplate(*commandRunnerConfig.PathsMapping.DataNodeHome, *nodeSet)
		if err != nil {
			mErr.Add(fmt.Errorf("failed to template data-node home path in the command runner config: %w", err))
		}
		remoteCommandRunner.PathsMapping.DataNodeHome = utils.StrPoint(dataNodeHomePath.String())

	}

	// set remote command runner to pass the most recent informations to the template generator
	nodeSet.RemoteCommandRunner = remoteCommandRunner
	if commandRunnerConfig.Nomad.JobTemplate != nil {
		rawTemplate, err := utils.GenerateTemplate(*commandRunnerConfig.Nomad.JobTemplate, nodeSet)

		remoteCommandRunner.NomadJobRaw = rawTemplate.String()

		if err != nil {
			mErr.Add(fmt.Errorf("failed to generate nomad template for remote command runner: %w", err))
		}
	}

	if mErr.HasAny() {
		return nil, mErr
	}

	return remoteCommandRunner, nil
}
