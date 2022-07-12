package importer

import (
	"fmt"
	"path/filepath"
	"strings"

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
)

func verifyTendermintNode(vegaBinary, tendermintHomePath, expectedNodeID string) error {
	out, err := utils.ExecuteBinary(vegaBinary, []string{"tm", "show-node-id", "--home", tendermintHomePath}, nil)
	if err != nil {
		return fmt.Errorf("failed to get the tendermint node id: %w", err)
	}

	tendermintNodeID := strings.ToLower(strings.Trim(string(out), " \t\n"))

	if tendermintNodeID != strings.ToLower(expectedNodeID) {
		return fmt.Errorf("tendermint node is invalid: expected \"%s\", got \"%s\"", expectedNodeID, tendermintNodeID)
	}

	return nil
}

func importTendermintNodeKey(tendermintHomePath, privateKey string) error {
	nodeKeyFilePath := filepath.Join(tendermintHomePath, "config", "node_key.json")
	rootNode, err := dasel.NewFromFile(nodeKeyFilePath, "json")
	if err != nil {
		return fmt.Errorf("failed to load private validator file: %w", err)
	}
	if err := rootNode.Put(".priv_key.value", privateKey); err != nil {
		return fmt.Errorf("failed to update address in she priv_validator_key.json file: %w", err)
	}

	if err := rootNode.WriteToFile(nodeKeyFilePath, "json", []storage.ReadWriteOption{}); err != nil {
		return fmt.Errorf("failed to write the node_key.json file: %w", err)
	}

	return nil
}

type tendermintPrivateValidatorData struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

func importTendermintPrivateValidator(tendermintHomePath string, data tendermintPrivateValidatorData) error {
	privateValidatorKeyFilePath := filepath.Join(tendermintHomePath, "config", "priv_validator_key.json")

	rootNode, err := dasel.NewFromFile(privateValidatorKeyFilePath, "json")
	if err != nil {
		return fmt.Errorf("failed to load private validator file: %w", err)
	}
	if err := rootNode.Put(".address", data.Address); err != nil {
		return fmt.Errorf("failed to update address in the priv_validator_key.json file: %w", err)
	}
	if err := rootNode.Put(".pub_key.value", data.PublicKey); err != nil {
		return fmt.Errorf("failed to update public key in the priv_validator_key.json file: %w", err)
	}
	if err := rootNode.Put(".priv_key.value", data.PrivateKey); err != nil {
		return fmt.Errorf("failed to update private key in the priv_validator_key.json file: %w", err)
	}
	if err := rootNode.WriteToFile(privateValidatorKeyFilePath, "json", []storage.ReadWriteOption{}); err != nil {
		return fmt.Errorf("failed to write the priv_validator_key.json file: %w", err)
	}

	return nil
}
