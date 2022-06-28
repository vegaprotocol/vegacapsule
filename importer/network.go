package importer

import (
	"fmt"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

const (
	defaultVegaIsolatedWalletName = "imported-wallet"
)

type NodeImportData struct {
	NodeIndex string `json:"node_index"`

	EthereumPrivateKey string `json:"ethereum_private_key"`

	TendermintValidatorPrivateKey string `json:"tendermint_validator_private_key"`
	TendermintValidatorPublicKey  string `json:"tendermint_validator_public_key"`
	TendermintValidatorAddress    string `json:"tendermint_validator_address"`
	TendermintValidatorNodeID     string `json:"tendermint_validator_node_id"`

	TendermintNodePrivateKey string `json:"tendermint_node_private_key"`

	VegaRecoveryPhrase string `json:"vega_recovery_phrase"`
}

type NetworkImportdata []NodeImportData

func importNodeData(nodeSet types.NodeSet, data NodeImportData) (*types.NodeSet, error) {
	ethereumKeystorePath := filepath.Join(nodeSet.Vega.HomeDir, "wallets", "ethereum")
	ethereumKeystorePassFilePath := filepath.Join(nodeSet.Vega.HomeDir, "ethereum-vega-wallet-pass.txt")
<<<<<<< HEAD
=======

	log.Println("... preparing tendermint validator keys")
	tendermintValidatorKey, err := decodeTendermintPrivateKey(data.TendermintValidatorPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private tendermint validator key: %w", err)
	}

	log.Println("... preparing tendermint node keys")
	tendermintNodeKey, err := decodeTendermintPrivateKey(data.TendermintNodePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private tendermint validator key: %w", err)
	}

	log.Println("... importing ethereum private wallet")
>>>>>>> 6af27fc (feat: improve description for network import command)
	importedEthereumKey, err := importPrivateKeyIntoKeystore(data.EthereumPrivateKey, ethereumKeystorePassFilePath, ethereumKeystorePath)
	if err != nil {
		return nil, fmt.Errorf("failed to import private ethereum key into node keystore")
	}

	err = importTendermintPrivateValidator(nodeSet.Tendermint.HomeDir, tendermintPrivateValidatorData{
		Address:    data.TendermintValidatorAddress,
		PublicKey:  data.TendermintValidatorPublicKey,
		PrivateKey: data.TendermintValidatorPrivateKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to import the tendermint validator private data")
	}

	if err := importTendermintNodeKey(nodeSet.Tendermint.HomeDir, data.TendermintNodePrivateKey); err != nil {
		return nil, fmt.Errorf("failed to import the tendermint node private key")
	}

	if err := verifyTendermintNode(nodeSet.Vega.BinaryPath, nodeSet.Tendermint.HomeDir, data.TendermintValidatorNodeID); err != nil {
		return nil, fmt.Errorf("failed to verify imported tendermint node: %w", err)
	}

	recoveryPhraseTempFilePath, err := createTempFile(data.VegaRecoveryPhrase)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file for vega recovery phrase: %w", err)
	}
	defer os.Remove(recoveryPhraseTempFilePath)

	vegaWalletPassFilePath := filepath.Join(nodeSet.Vega.HomeDir, "vega-wallet-pass.txt")
	vegaWalletFilePath, err := createIsolatedVegaWallet("vega", isolatedVegaWallet{
		VegaHomePath:           nodeSet.Vega.HomeDir,
		RecoveryPhraseFilePath: recoveryPhraseTempFilePath,
		VegaWalletPassFilePath: vegaWalletPassFilePath,
		IsolatedWalletName:     defaultVegaIsolatedWalletName,
	}, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create isolated vega wallet: %w", err)
	}

	nodeWalletPassFile := filepath.Join(nodeSet.Vega.HomeDir, "node-vega-wallet-pass.txt")
	importNodeWalletData := importNodeWalletInput{
		VegaHomePath:       nodeSet.Vega.HomeDir,
		TendermintHomePath: nodeSet.Tendermint.HomeDir,
		PassphraseFilePath: nodeWalletPassFile,

		EthKeystoreFilePath:     importedEthereumKey.keystoreFilePath,
		EthKeystorePassFilePath: ethereumKeystorePassFilePath,

		VegaWalletFilePath:     vegaWalletFilePath,
		VegaWalletPassFilePath: vegaWalletPassFilePath,
	}

	if err := importVegaNodeWallet("vega", importNodeWalletData); err != nil {
		return nil, fmt.Errorf("failed to add imported wallets to node wallet: %w", err)
	}

	nodeSet.Vega.NodeWalletInfo.EthereumAddress = importedEthereumKey.ethereumAddress
	nodeSet.Vega.NodeWalletInfo.EthereumPrivateKey = data.EthereumPrivateKey
	// nodeSet.Vega.NodeWalletInfo.VegaWalletPublicKey = ... // TODO: Fix it
	nodeSet.Vega.NodeWalletInfo.VegaWalletRecoveryPhrase = data.VegaRecoveryPhrase

	return &nodeSet, nil
}

func getNodeImportDataForNodeIdx(nodeIdx int, networkData NetworkImportdata) *NodeImportData {
	for _, nodeData := range networkData {
		if fmt.Sprintf("%d", nodeIdx) == nodeData.NodeIndex {
			return &nodeData
		}
	}

	return nil
}

func ImportDataIntoNetworkValidators(state state.NetworkState, networkData NetworkImportdata) (*state.NetworkState, error) {
	errs := utils.NewMultiError()

	for idx, nodeSet := range state.GeneratedServices.NodeSets {
		if nodeSet.Mode != types.NodeModeValidator {
			continue
		}

		nodeData := getNodeImportDataForNodeIdx(nodeSet.Index, networkData)
		if nodeData == nil {
			errs.Add(fmt.Errorf("failed to import data for the %s node set: missing import data for the %d node in given set",
				nodeSet.Name,
				nodeSet.Index))

			continue
		}

		newNodeSet, err := importNodeData(nodeSet, *nodeData)
		if err != nil {
			errs.Add(fmt.Errorf("error importing new data for the %s node set: %w", nodeSet.Name, err))
			continue
		}

		state.GeneratedServices.NodeSets[idx] = *newNodeSet
	}

	if errs.HasAny() {
		return nil, errs
	}

	return &state, nil
}