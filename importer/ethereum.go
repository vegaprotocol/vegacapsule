package importer

import (
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/ethereum"
	vegagen "code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

func importEthereumKey(nodeSet types.NodeSet, ethereumPrivateKey string) (*ethereum.KeyStore, error) {
	log.Println("... importing ethereum private key")
	importedKeystore, err := ethereum.ImportPrivateKeyIntoKeystore(ethereumPrivateKey, nodeSet.Vega.NodeWalletInfo.EthereumPassFilePath, vegagen.EthereumWalletPath(nodeSet.Vega.HomeDir))
	if err != nil {
		return nil, fmt.Errorf("failed to import private ethereum key into node keystore: %w", err)
	}

	log.Println("... adding ethereum wallet to the nodewallet")
	ethImportArgs := []string{
		"nodewallet", "import", "--force",
		"--home", nodeSet.Vega.HomeDir,
		"--chain", types.NodeWalletChainTypeEthereum,
		"--passphrase-file", nodeSet.Vega.NodeWalletPassFilePath,
		"--output", "json",
		"--wallet-passphrase-file", nodeSet.Vega.NodeWalletInfo.EthereumPassFilePath,
		"--wallet-path", importedKeystore.FilePath,
	}
	if _, err := utils.ExecuteBinary(nodeSet.Vega.BinaryPath, ethImportArgs, nil); err != nil {
		return nil, fmt.Errorf("failed to import ethereum into vega node wallet: %w", err)
	}

	return importedKeystore, nil
}
