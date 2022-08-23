package importer

import (
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type NodeKeysData struct {
	NodeIndex                     int    `json:"node_index"`
	EthereumPrivateKey            string `json:"ethereum_private_key"`
	TendermintValidatorPrivateKey string `json:"tendermint_validator_private_key"`
	TendermintNodePrivateKey      string `json:"tendermint_node_private_key"`

	VegaRecoveryPhrase string `json:"vega_recovery_phrase"`
}

type NodesKeysData []NodeKeysData

func createAndImportNodeWallets(nodeSet types.NodeSet, data NodeKeysData) (*types.NodeSet, error) {
	log.Printf("importing keys for the \"%s\" node", nodeSet.Name)

	tmKeys, err := decodeAndImportTendermintKeys(nodeSet, data.TendermintValidatorPrivateKey, data.TendermintNodePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode and import tendermint keys: %w", err)
	}

	importedKeystore, err := importEthereumKey(nodeSet, data.EthereumPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to import ethereum key: %w", err)
	}

	isolatedVegaWallet, err := createAndImportVegaWallet(nodeSet, data.VegaRecoveryPhrase)
	if err != nil {
		return nil, fmt.Errorf("failed to create and import isolated vega wallet: %w", err)
	}

	log.Println("... updating vegacapsule state")
	nodeSet.Vega.NodeWalletInfo.EthereumAddress = importedKeystore.KeyPair.Address
	nodeSet.Vega.NodeWalletInfo.EthereumPrivateKey = data.EthereumPrivateKey
	nodeSet.Vega.NodeWalletInfo.VegaWalletPublicKey = isolatedVegaWallet.WalletPublicKey
	nodeSet.Vega.NodeWalletInfo.VegaWalletRecoveryPhrase = data.VegaRecoveryPhrase
	nodeSet.Vega.NodeWalletInfo.VegaWalletName = defaultVegaIsolatedWalletName
	nodeSet.Tendermint.NodeID = tmKeys.nodeKey.NodeID
	log.Printf("importing data for the \"%s\" node finished", nodeSet.Name)

	return &nodeSet, nil
}

func getNodeWalletsDataForNodeIdx(nodeIdx int, keys NodesKeysData) *NodeKeysData {
	for _, nodeWalletData := range keys {
		if nodeIdx == nodeWalletData.NodeIndex {
			return &nodeWalletData
		}
	}

	return nil
}

func ImportKeysIntoValidatorsWallets(state state.NetworkState, keys NodesKeysData) (*state.NetworkState, error) {
	errs := utils.NewMultiError()

	for idx, nodeSet := range state.GeneratedServices.NodeSets {
		if nodeSet.Mode != types.NodeModeValidator {
			continue
		}

		nodeData := getNodeWalletsDataForNodeIdx(nodeSet.Index, keys)
		if nodeData == nil {
			errs.Add(fmt.Errorf("failed to import data for the %s node set: missing import data for the %d node in given set",
				nodeSet.Name,
				nodeSet.Index))

			continue
		}

		newNodeSet, err := createAndImportNodeWallets(nodeSet, *nodeData)
		if err != nil {
			errs.Add(fmt.Errorf("failed to create and import node wallets for the %s node set: %w", nodeSet.Name, err))
			continue
		}

		state.GeneratedServices.NodeSets[idx] = *newNodeSet
	}

	if err := regenerateTendermintConfig(state.GeneratedServices.NodeSets, state.Config); err != nil {
		errs.Add(fmt.Errorf("failed to update tendermint configuration files: %w", err))
	}

	if err := updateGenesis(state); err != nil {
		errs.Add(fmt.Errorf("failed to import genesis for network: %w", err))
	}

	if errs.HasAny() {
		return nil, errs
	}

	return &state, nil
}
