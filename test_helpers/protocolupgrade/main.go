package main

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	flag "github.com/spf13/pflag"
)

func submitProtocolUpgrage(
	val types.VegaNodeOutput,
	blockHeight uint64,
	vegaReleaseTag, dataNodeReleaseTag string,
) error {
	b, err := utils.ExecuteBinary("vegawallet", []string{
		"--home", val.HomeDir,
		"command", "send",
		"--wallet", "created-wallet",
		"--node-address", "localhost:3002",
		"--passphrase-file", val.NodeWalletInfo.VegaWalletPassFilePath,
		"--pubkey", val.NodeWalletInfo.VegaWalletPublicKey,
		fmt.Sprintf(
			`{"protocolUpgradeProposal": {"upgradeBlockHeight": "%d", "vegaReleaseTag": "%s", "dataNodeReleaseTag": "%s"}}`,
			blockHeight,
			vegaReleaseTag,
			dataNodeReleaseTag,
		),
	}, nil)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func templateVegaConfig(val types.VegaNodeOutput) error {
	args := []string{
		"template",
		"node-sets",
		"--nodeset-name", val.NomadJobName,
		"--type", "vega",
		"--path", "/Users/karel/work/vegacapsule/net_confs/node_set_templates/default/vega_full_visor_snap.tmpl",
		"--with-merge",
		"--update-network",
	}

	fmt.Printf("Calling %s %v", "vegacapsule", args)
	b, err := utils.ExecuteBinary("vegacapsule", args, nil)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

var (
	blockHeight        uint64
	vegaReleaseTag     string
	dataNodeReleaseTag string
)

func init() {
	flag.Uint64Var(&blockHeight, "height", 0, "Upgrade Block Height")
	flag.StringVar(&vegaReleaseTag, "vega-rt", "0.0.2", "Vega Release Tag")
	flag.StringVar(&dataNodeReleaseTag, "data-node-rt", "0.0.2", "Data Node Release Tag")
}

func main() {
	flag.Parse()

	if blockHeight == 0 {
		panic("--height must be bigger then 0")
	}

	validators := []types.VegaNodeOutput{}

	if _, err := utils.ExecuteBinary("vegacapsule", []string{"nodes", "ls-validators"}, &validators); err != nil {
		panic(err)
	}

	for _, val := range validators {
		if err := templateVegaConfig(val); err != nil {
			panic(err)
		}
		if err := submitProtocolUpgrage(val, blockHeight, vegaReleaseTag, dataNodeReleaseTag); err != nil {
			panic(err)
		}
	}
}
