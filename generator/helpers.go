package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/types"
)

const (
	nodesRegistryFileName = "nodes-wallets.json"
)

func (g Generator) persistNodesWalletsInfo(homePath string, validators []types.NodeSet, nonValidators []types.NodeSet) error {
	registry := types.NodeWallets{
		Wallets: make([]types.NodeWalletInfo, len(validators)+len(nonValidators)),
	}

	for idx, node := range append(validators, nonValidators...) {
		registry.Wallets[idx] = node.Vega.NodeWalletInfo
	}

	registryJson, err := json.MarshalIndent(registry, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to marshal node wallets info: %w", err)
	}

	filePath := filepath.Join(homePath, nodesRegistryFileName)
	if err := ioutil.WriteFile(filePath, registryJson, 0644); err != nil {
		return fmt.Errorf("failed to write nodes wallets info to file '%s': %w", filePath, err)
	}

	return nil
}
