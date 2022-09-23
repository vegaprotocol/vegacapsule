package importer

import (
	"fmt"
	"log"
	"os"

	vegagen "code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

const (
	defaultVegaIsolatedWalletName = "imported-wallet"
)

type createIsolatedVegaWalletOutput struct {
	WalletFilePath  string
	WalletPublicKey string
}

func createIsolatedVegaWallet(nodeSet types.NodeSet, recoveryPhraseFilePath string, isolatedWalletName string, force bool) (*createIsolatedVegaWalletOutput, error) {
	vegaWalletFilePath := vegagen.IsolatedWalletPath(nodeSet.Vega.HomeDir, isolatedWalletName)

	if force {
		if err := os.RemoveAll(vegaWalletFilePath); err != nil {
			return nil, fmt.Errorf("failed to remove existing vega wallet: %w", err)
		}
	}

	args := []string{
		"wallet", "import",
		"--home", nodeSet.Vega.HomeDir,
		"--output", "json",
		"--recovery-phrase-file", recoveryPhraseFilePath,
		"--passphrase-file", nodeSet.Vega.NodeWalletInfo.VegaWalletPassFilePath,
		"--wallet", isolatedWalletName,
	}

	importOut := &struct {
		Key struct {
			Public string `json:"publicKey"`
		} `json:"key"`
	}{}

	if _, err := utils.ExecuteBinary(nodeSet.Vega.BinaryPath, args, importOut); err != nil {
		return nil, fmt.Errorf("failed to create isolated vega wallet: %w", err)
	}

	return &createIsolatedVegaWalletOutput{
		WalletFilePath:  vegaWalletFilePath,
		WalletPublicKey: importOut.Key.Public,
	}, nil
}

func createAndImportVegaWallet(nodeSet types.NodeSet, recoveryPhrase string) (*createIsolatedVegaWalletOutput, error) {
	log.Println("... create isolated vega wallet from given recovery passphrase")
	recoveryPhraseTempFilePath, err := createTempFile(recoveryPhrase)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file for vega recovery phrase: %w", err)
	}
	defer os.Remove(recoveryPhraseTempFilePath)

	createIsolatedVegaWalletOutput, err := createIsolatedVegaWallet(nodeSet, recoveryPhraseTempFilePath, defaultVegaIsolatedWalletName, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create isolated vega wallet: %w", err)
	}

	log.Println("... adding isolated vega wallet to the nodewallet")
	vegaImpotyArgs := []string{
		"nodewallet", "import", "--force",
		"--home", nodeSet.Vega.HomeDir,
		"--chain", types.NodeWalletChainTypeVega,
		"--passphrase-file", nodeSet.Vega.NodeWalletPassFilePath,
		"--output", "json",
		"--wallet-passphrase-file", nodeSet.Vega.NodeWalletInfo.VegaWalletPassFilePath,
		"--wallet-path", createIsolatedVegaWalletOutput.WalletFilePath,
	}
	if _, err := utils.ExecuteBinary(nodeSet.Vega.BinaryPath, vegaImpotyArgs, nil); err != nil {
		return nil, fmt.Errorf("failed to import ethereum to vega nodewallet: %w", err)
	}

	return createIsolatedVegaWalletOutput, nil
}
