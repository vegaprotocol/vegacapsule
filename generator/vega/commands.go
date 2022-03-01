package vega

import (
	"log"

	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type initateNodeOutput struct {
	ConfigFilePath           string `json:"configFilePath"`
	NodeWalletConfigFilePath string `json:"nodeWalletConfigFilePath"`
}

type generateNodeWalletOutput struct {
	Mnemonic         string `json:"mnemonic,omitempty"`
	RegistryFilePath string `json:"registryFilePath"`
	WalletFilePath   string `json:"walletFilePath"`
}

type importNodeWalletOutput struct {
	RegistryFilePath string `json:"registryFilePath"`
	TendermintPubkey string `json:"tendermintPubkey"`
	WalletFilePath   string `json:"walletFilePath"`
}

type createWalletOutput struct {
	Wallet struct {
		Name           string `json:"name"`
		FilePath       string `json:"filePath"`
		RecoveryPhrase string `json:"recoveryPhrase"`
	} `json:"wallet"`
	Key struct {
		Public string `json:"publicKey"`
	} `json:"key"`
}

func (vg ConfigGenerator) initiateNode(homePath string, nodeWalletPhraseFile string, nodeMode string) (*initateNodeOutput, error) {
	args := []string{
		"init",
		"--home", homePath,
		"--nodewallet-passphrase-file", nodeWalletPhraseFile,
		"--output", "json",
		string(nodeMode),
	}

	log.Printf("Initiating node %q with: %v", nodeMode, args)

	out := &initateNodeOutput{}
	if _, err := utils.ExecuteBinary(vg.conf.VegaBinary, args, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (vg ConfigGenerator) createWallet(homePath, name, walletPhraseFilePath string) (*createWalletOutput, error) {
	args := []string{
		"wallet",
		"--home", homePath,
		"create",
		"--output", "json",
		"--passphrase-file", walletPhraseFilePath,
		"--wallet", name,
	}

	log.Printf("Creating vega wallet with: %v", args)

	out := &createWalletOutput{}
	if _, err := utils.ExecuteBinary(vg.conf.VegaBinary, args, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (vg ConfigGenerator) generateNodeWallet(homePath string, nodeWalletPhraseFile string, walletPhraseFile string, walletType string) (*generateNodeWalletOutput, error) {
	args := []string{
		"nodewallet",
		"--home", homePath,
		"--passphrase-file", nodeWalletPhraseFile,
		"generate",
		"--output", "json",
		"--chain", walletType,
		"--wallet-passphrase-file", walletPhraseFile,
	}

	log.Printf("Generating node %q wallet with: %v", walletType, args)

	out := &generateNodeWalletOutput{}
	if _, err := utils.ExecuteBinary(vg.conf.VegaBinary, args, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (vg ConfigGenerator) importTendermintNodeWallet(homePath string, nodeWalletPhraseFile string, tendermintHomePath string) (*importNodeWalletOutput, error) {
	args := []string{
		"nodewallet",
		"--home", homePath,
		"--passphrase-file", nodeWalletPhraseFile,
		"import",
		"--output", "json",
		"--chain", "tendermint",
		"--tendermint-home", tendermintHomePath,
	}

	log.Printf("Importing tenderming wallet: %v", args)

	nwo := &importNodeWalletOutput{}
	if _, err := utils.ExecuteBinary(vg.conf.VegaBinary, args, nwo); err != nil {
		return nil, err
	}

	return nwo, nil
}

func (vg ConfigGenerator) importVegaNodeWallet(homePath, nodeWalletPhraseFile, walletPhraseFile, walletFilePath string) (*importNodeWalletOutput, error) {
	args := []string{
		"nodewallet",
		"--home", homePath,
		"--passphrase-file", nodeWalletPhraseFile,
		"import",
		"--output", "json",
		"--chain", types.NodeWalletChainTypeVega,
		"--wallet-passphrase-file", walletPhraseFile,
		"--wallet-path", walletFilePath,
	}

	log.Printf("Importing node vega wallet with: %v", args)

	out := &importNodeWalletOutput{}
	if _, err := utils.ExecuteBinary(vg.conf.VegaBinary, args, out); err != nil {
		return nil, err
	}

	return out, nil
}
