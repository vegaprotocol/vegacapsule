package cmd

import (
	"fmt"
	"runtime/debug"

	vgjson "code.vegaprotocol.io/vega/libs/json"
	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/utils"

	"github.com/spf13/cobra"
)

var (
	cLIVersion     = "v0.72.2+dev"
	cLIVersionHash = ""
	withDeps       bool
)

type versionOutput struct {
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

type versionWithNameOutput struct {
	Name string `json:"name"`
	Path string `json:"path,omitempty"`
	versionOutput
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display software version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !withDeps {
			vgjson.Print(versionOutput{
				Version: cLIVersion,
				Hash:    cLIVersionHash,
			})
			return nil
		}

		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Config == nil {
			return fmt.Errorf("failed to display versions of with dependency binaries: missing network configuration")
		}

		versions := []*versionWithNameOutput{
			{
				Name: "vegacapsule",
				versionOutput: versionOutput{
					Version: cLIVersion,
					Hash:    cLIVersionHash,
				},
			},
		}

		vegaVersion, err := getBinaryVersion(netState.Config.GetVegaBinary(), "vega", "")
		if err != nil {
			return err
		}
		versions = append(versions, vegaVersion)

		if netState.GeneratedServices.Wallet != nil {
			walletVersion, err := getBinaryVersion(
				netState.GeneratedServices.Wallet.BinaryPath,
				netState.GeneratedServices.Wallet.Name,
				config.WalletSubCmd,
			)
			if err != nil {
				return err
			}
			versions = append(versions, walletVersion)
		}

		for _, ns := range netState.GeneratedServices.NodeSets {
			vegaVersion, err := getBinaryVersion(
				ns.Vega.BinaryPath,
				fmt.Sprintf("%s %s", ns.Name, ns.Vega.Name),
				"",
			)
			if err != nil {
				return err
			}
			versions = append(versions, vegaVersion)

			if ns.DataNode != nil {
				dataNodeVersion, err := getBinaryVersion(
					ns.DataNode.BinaryPath,
					fmt.Sprintf("%s %s", ns.Name, ns.DataNode.Name),
					config.DataNodeSubCmd,
				)
				if err != nil {
					return err
				}
				versions = append(versions, dataNodeVersion)
			}
		}

		vgjson.PrettyPrint(versions)

		return nil
	},
}

func init() {
	versionCmd.Flags().BoolVar(&withDeps,
		"with-deps",
		false,
		"Allows to print versions of currently used vega, data-node and vegawallet binaries.",
	)

	setVersionHash()
}

func getBinaryVersion(path, name, subCmd string) (*versionWithNameOutput, error) {
	args := []string{"version", "--output", "json"}

	if subCmd != "" {
		args = append([]string{subCmd}, args...)
	}

	var version versionOutput
	if _, err := utils.ExecuteBinary(path, args, &version); err != nil {
		return nil, err
	}

	return &versionWithNameOutput{
		Name:          name,
		Path:          path,
		versionOutput: version,
	}, nil
}

func setVersionHash() {
	info, _ := debug.ReadBuildInfo()
	modified := false

	for _, v := range info.Settings {
		if v.Key == "vcs.revision" {
			cLIVersionHash = v.Value
		}
		if v.Key == "vcs.modified" && v.Value == "true" {
			modified = true
		}
	}
	if modified {
		cLIVersionHash += "-modified"
	}
}
