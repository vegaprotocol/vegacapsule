package main

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

func main() {
	validators := []types.VegaNodeOutput{}

	if _, err := utils.ExecuteBinary("vegacapsule", []string{"nodes", "ls-validators"}, &validators); err != nil {
		panic(err)
	}

	for _, val := range validators {
		b, err := utils.ExecuteBinary("vegawallet", []string{
			"--home", val.HomeDir,
			"command", "send",
			"--wallet", "created-wallet",
			"--node-address", "localhost:3002",
			"--passphrase-file", val.NodeWalletInfo.VegaWalletPassFilePath,
			"--pubkey", val.NodeWalletInfo.VegaWalletPublicKey,
			`{"protocolUpgradeProposal": {"upgradeBlockHeight": "50", "vegaReleaseTag": "0.0.2", "dataNodeReleaseTag": "0.0.2"}}`,
		}, nil)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(b))
	}
}
