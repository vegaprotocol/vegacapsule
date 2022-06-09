package generator

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/nomad"
	"code.vegaprotocol.io/vegacapsule/generator/wallet"
	"code.vegaprotocol.io/vegacapsule/types"
)

func (g *Generator) initiateNodeSet(index int, n config.NodeConfig) (*types.NodeSet, error) {
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
		n.EthereumWalletAddressClef,
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
	validatorsSet := []types.NodeSet{}
	nonValidatorsSet := []types.NodeSet{}

	var index int
	for _, n := range g.conf.Network.Nodes {
		for i := 0; i < n.Count; i++ {
			templates, err := g.templatePreGenerateJobs(n.PreGenerate, index)
			if err != nil {
				return nil, fmt.Errorf("failed to template pre generate jobs for node set %q-%q: %w", n.Name, index, err)
			}
			preGenJobIDs, err := g.startNomadJobs(templates)
			if err != nil {
				return nil, fmt.Errorf("failed to start pre generate jobs for node set %s-%d: %w", n.Name, index, err)
			}

			nodeSet, err := g.initiateNodeSet(index, n)
			if err != nil {
				return nil, err
			}

			nodeSet.PreGenerateJobsIDs = preGenJobIDs

			if n.Mode == types.NodeModeValidator {
				validatorsSet = append(validatorsSet, *nodeSet)
			} else {
				nonValidatorsSet = append(nonValidatorsSet, *nodeSet)
			}

			index++
		}
	}

	return &nodeSets{
		validators:    validatorsSet,
		nonValidators: nonValidatorsSet,
	}, nil
}

func (g *Generator) templatePreGenerateJobs(preGenConf *config.PreGenerate, index int) ([]string, error) {
	if preGenConf == nil {
		return []string{}, nil
	}

	jobTemplates := make([]string, 0, len(preGenConf.Nomad))
	for _, nc := range preGenConf.Nomad {
		if nc.JobTemplate == nil {
			continue
		}

		template, err := nomad.GeneratePreGenerateTemplate(*nc.JobTemplate, nomad.PreGenerateTemplateCtx{
			Name:  nc.Name,
			Index: index,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to template nomad job for pre generate %q: %w", nc.Name, err)
		}

		jobTemplates = append(jobTemplates, template.String())
	}

	return jobTemplates, nil
}

func (g *Generator) startNomadJobs(rawNomadJobs []string) ([]string, error) {
	if len(rawNomadJobs) == 0 {
		return rawNomadJobs, nil
	}

	jobs, err := g.jobRunner.RunRawNomadJobs(context.Background(), rawNomadJobs)
	if err != nil {
		return nil, fmt.Errorf("failed to run node set pre generate job: %w", err)
	}

	jobIDs := make([]string, 0, len(jobs))
	for _, j := range jobs {
		jobIDs = append(jobIDs, *j.ID)
	}

	return jobIDs, nil
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
