package vega

import (
	"log"

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
}

func (vg ConfigGenerator) initiateNode(homePath string, nodeWalletPhraseFile string, nodeMode string) (*initateNodeOutput, error) {
	args := []string{
		"init",
		"--home", homePath,
		"--nodewallet-passphrase-file", nodeWalletPhraseFile,
		"--output", "json",
		string(nodeMode),
	}

	log.Printf("Initiating node %q wallet with: %v", nodeMode, args)

	out := &initateNodeOutput{}
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
		"--chain", string(walletType),
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

	log.Printf("Generating tenderming wallet: %v", args)

	nwo := &importNodeWalletOutput{}
	if _, err := utils.ExecuteBinary(vg.conf.VegaBinary, args, nwo); err != nil {
		return nil, err
	}

	return nwo, nil
}
