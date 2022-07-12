package importer

import (
	"fmt"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type isolatedVegaWallet struct {
	VegaHomePath           string
	RecoveryPhraseFilePath string
	VegaWalletPassFilePath string
	IsolatedWalletName     string
}

type isolatedWalletOutput struct {
	WalletFilePath  string
	WalletPublicKey string
}

func createIsolatedVegaWallet(vegaBinary string, data isolatedVegaWallet, force bool) (*isolatedWalletOutput, error) {
	vegaWalletFilePath := filepath.Join(data.VegaHomePath, "data", "wallets", data.IsolatedWalletName)

	if force {
		if err := os.RemoveAll(vegaWalletFilePath); err != nil {
			return nil, fmt.Errorf("failed to remove existing vega wallet: %w", err)
		}
	}

	args := []string{
		"wallet", "import",
		"--home", data.VegaHomePath,
		"--no-version-check",
		"--output", "json",
		"--recovery-phrase-file", data.RecoveryPhraseFilePath,
		"--passphrase-file", data.VegaWalletPassFilePath,
		"--wallet", data.IsolatedWalletName,
	}

	importOut := &struct {
		Key struct {
			Public string `json:"publicKey"`
		} `json:"key"`
	}{}

	if _, err := utils.ExecuteBinary(vegaBinary, args, importOut); err != nil {
		return nil, fmt.Errorf("failed to create isolated vega wallet: %w", err)
	}

	return &isolatedWalletOutput{
		WalletFilePath:  vegaWalletFilePath,
		WalletPublicKey: importOut.Key.Public,
	}, nil
}

type importNodeWalletInput struct {
	VegaHomePath       string
	TendermintHomePath string
	PassphraseFilePath string

	EthKeystoreFilePath     string
	EthKeystorePassFilePath string

	VegaWalletFilePath     string
	VegaWalletPassFilePath string
}

func importVegaNodeWallet(vegaBinary string, data importNodeWalletInput) error {
	tmImportArgs := []string{
		"nodewallet", "import", "--force",
		"--home", data.VegaHomePath,
		"--chain", types.NodeWalletChainTypeTendermint,
		"--passphrase-file", data.PassphraseFilePath,
		"--output", "json",
		"--tendermint-home", data.TendermintHomePath,
	}
	if _, err := utils.ExecuteBinary(vegaBinary, tmImportArgs, nil); err != nil {
		return fmt.Errorf("failed to import tendermint to vega nodewallet: %w", err)
	}

	ethImportArgs := []string{
		"nodewallet", "import", "--force",
		"--home", data.VegaHomePath,
		"--chain", types.NodeWalletChainTypeEthereum,
		"--passphrase-file", data.PassphraseFilePath,
		"--output", "json",
		"--wallet-passphrase-file", data.EthKeystorePassFilePath,
		"--wallet-path", data.EthKeystoreFilePath,
	}
	if _, err := utils.ExecuteBinary(vegaBinary, ethImportArgs, nil); err != nil {
		return fmt.Errorf("failed to import ethereum to vega nodewallet: %w", err)
	}

	vegaImpotyArgs := []string{
		"nodewallet", "import", "--force",
		"--home", data.VegaHomePath,
		"--chain", types.NodeWalletChainTypeVega,
		"--passphrase-file", data.PassphraseFilePath,
		"--output", "json",
		"--wallet-passphrase-file", data.VegaWalletPassFilePath,
		"--wallet-path", data.VegaWalletFilePath,
	}
	if _, err := utils.ExecuteBinary(vegaBinary, vegaImpotyArgs, nil); err != nil {
		return fmt.Errorf("failed to import ethereum to vega nodewallet: %w", err)
	}

	return nil
}
