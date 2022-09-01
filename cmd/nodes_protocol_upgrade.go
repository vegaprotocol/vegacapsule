package cmd

import (
	"fmt"
	"strconv"

	"code.vegaprotocol.io/vegacapsule/commands"
	"code.vegaprotocol.io/vegacapsule/generator/visor"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"
)

var (
	upgradeBlockHeight           int64
	upgradeReleaseTag            string
	upgradeRunConfigTemplateFile string
	upgradePropose               bool
	upgradeForce                 bool
	upgradeInclude               []string
	upgradeExclude               []string
)

var nodesProtocolUpgradeCmd = &cobra.Command{
	Use:   "protocol-upgrade",
	Short: "Prepares protocol upgrade for all running nodes and send transaction to network if allowed",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("protocol-upgrade")
		}

		if upgradeBlockHeight == 0 {
			return fmt.Errorf("parameter height can not be zero")
		}

		if len(upgradeInclude) != 0 && len(upgradeExclude) != 0 {
			return fmt.Errorf("combining flags include-nodes and exclude-nodes is not allowed")
		}

		visorGen, err := visor.NewGenerator(netState.Config)
		if err != nil {
			return fmt.Errorf("failed to create new visor generator: %w", err)
		}

		runTemplateRaw, err := netState.Config.LoadConfigTemplateFile(upgradeRunConfigTemplateFile)
		if err != nil {
			return err
		}

		visorRunTmpl, err := visor.NewConfigTemplate(runTemplateRaw)
		if err != nil {
			return err
		}

		nodeSets, err := filtertUpgradeNodeSet(*netState.GeneratedServices, upgradeInclude, upgradeExclude)
		if err != nil {
			return err
		}

		for _, ns := range nodeSets {
			if err := visorGen.PrepareUpgrade(ns.Index, upgradeReleaseTag, ns, visorRunTmpl, upgradeForce); err != nil {
				return err
			}
		}

		if !upgradePropose {
			return nil
		}

		for _, ns := range nodeSets {
			b, err := commands.VegaProtocolUpgradeProposal(
				*netState.Config.VegaBinary,
				ns.Vega.HomeDir,
				upgradeReleaseTag,
				strconv.FormatInt(upgradeBlockHeight, 10),
				ns.Vega.NodeWalletPassFilePath,
			)
			if err != nil {
				return fmt.Errorf("failed to submit protocol upgrade proposal to node %q: %w", ns.Name, err)
			}

			fmt.Printf("applied protocol upgrade for node set %q: %s \n", ns.Name, b)
		}

		return nil
	},
}

func init() {
	nodesProtocolUpgradeCmd.Flags().Int64Var(&upgradeBlockHeight,
		"height",
		0,
		"The block height at which the upgrade should be made",
	)
	nodesProtocolUpgradeCmd.Flags().StringVar(&upgradeReleaseTag,
		"release-tag",
		"",
		"A valid vega core release tag for the upgrade proposal",
	)
	nodesProtocolUpgradeCmd.Flags().StringVar(&upgradeRunConfigTemplateFile,
		"template-path",
		"",
		"Run config template to be applied",
	)
	nodesProtocolUpgradeCmd.Flags().BoolVar(&upgradePropose,
		"propose",
		false,
		"Automatically sends protocol upgrade proposal transaction to network",
	)
	nodesProtocolUpgradeCmd.Flags().StringSliceVar(&upgradeInclude,
		"include-nodes",
		nil,
		"IDs of node that should be included in the upgrade. Can not be combined with include-nodes",
	)
	nodesProtocolUpgradeCmd.Flags().StringSliceVar(&upgradeExclude,
		"exclude-nodes",
		nil,
		"IDs of node that should be excluded from the upgrade. Can not be combined with exclude-nodes",
	)
	nodesProtocolUpgradeCmd.Flags().BoolVar(&upgradeForce,
		"force",
		false,
		"Forces to run upgrade",
	)
	nodesProtocolUpgradeCmd.MarkFlagRequired("height")
	nodesProtocolUpgradeCmd.MarkFlagRequired("release-tag")
	nodesProtocolUpgradeCmd.MarkFlagRequired("template-path")
}

func filtertUpgradeNodeSet(genS types.GeneratedServices, upgradeInclude, upgradeExclude []string) ([]types.NodeSet, error) {
	if len(upgradeInclude) != 0 {
		var nodeSets []types.NodeSet

		for _, nodeSetName := range upgradeInclude {
			ns, err := genS.GetNodeSet(nodeSetName)
			if err != nil {
				return nil, fmt.Errorf("failed to get requested node set: %w", err)
			}
			nodeSets = append(nodeSets, *ns)
		}

		return nodeSets, nil
	}

	if len(upgradeExclude) != 0 {
		var nodeSets []types.NodeSet

		for _, ns := range genS.NodeSets {
			for _, excludeNodeName := range upgradeExclude {
				if ns.Name == excludeNodeName {
					continue
				}
			}

			nodeSets = append(nodeSets, ns)
		}

		return nodeSets, nil
	}

	return genS.NodeSets.ToSlice(), nil
}
