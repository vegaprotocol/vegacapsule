package cmd

import (
	"fmt"
	"strconv"

	"code.vegaprotocol.io/vegacapsule/commands"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var (
	upgradeBlockHeight           int64
	upgradeReleaseTag            string
	upgradeRunConfigTemplateFile string
)

var nodesProtocolUpgradeCmd = &cobra.Command{
	Use:   "protocol-upgrade",
	Short: "Schedules protocol upgrade proposal for all running nodes",
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

		// visorGen, err := visor.NewGenerator(netState.Config)
		// if err != nil {
		// 	return fmt.Errorf("failed to create new visor generator: %w", err)
		// }

		// visorRunTmpl, err := visor.NewConfigTemplate(upgradeRunConfigTemplateFile)
		// if err != nil {
		// 	return err
		// }

		// for _, ns := range netState.GeneratedServices.NodeSets {
		// 	if err := visorGen.PrepareUpgrade(ns.Index, upgradeReleaseTag, ns, visorRunTmpl); err != nil {
		// 		return err
		// 	}
		// }

		for _, ns := range netState.GeneratedServices.NodeSets {
			b, err := commands.VegaProtocolUpgradeProposal(
				*netState.Config.VegaBinary,
				ns.Tendermint.HomeDir,
				upgradeReleaseTag,
				strconv.FormatInt(upgradeBlockHeight, 10),
				ns.Vega.NodeWalletPassFilePath,
			)
			if err != nil {
				return fmt.Errorf("failed to restore node %q from checkpoint: %w", ns.Name, err)
			}

			fmt.Printf("applied protocol upgrade for node set %q: %s", ns.Name, b)
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
	nodesProtocolUpgradeCmd.Flags().StringVar(&upgradeReleaseTag,
		"template-path",
		"",
		"Run config template to be applied",
	)
	nodesProtocolUpgradeCmd.MarkFlagRequired("height")
	nodesProtocolUpgradeCmd.MarkFlagRequired("release-tag")
	// nodesProtocolUpgradeCmd.MarkFlagRequired("template-path")
}
