package importer

import (
	"fmt"
	"log"
	"strings"

	tmgen "code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/types"
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
	nodeKeyFilePath := tmgen.NodeKeyFilePath(tendermintHomePath)
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
	privateValidatorKeyFilePath := tmgen.PrivValidatorFilePath(tendermintHomePath)

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

func decodeAndImportTendermintKeys(nodeSet types.NodeSet, tmValidatorPrivateKey string, tmNodePrivateKey string) error {
	log.Println("... preparing tendermint validator keys")
	tendermintValidatorKey, err := decodeBase64TendermintPrivateKey(tmValidatorPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decode private tendermint validator key: %w", err)
	}

	log.Println("... preparing tendermint node keys")
	tendermintNodeKey, err := decodeBase64TendermintPrivateKey(tmNodePrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decode private tendermint node key: %w", err)
	}

	log.Println("... importing tendermint validator keys")
	err = importTendermintPrivateValidator(nodeSet.Tendermint.HomeDir, tendermintPrivateValidatorData{
		Address:    tendermintValidatorKey.NodeID,
		PublicKey:  tendermintValidatorKey.PublicKey,
		PrivateKey: tendermintValidatorKey.PrivateKey,
	})
	if err != nil {
		return fmt.Errorf("failed to import the tendermint validator private data")
	}

	log.Println("... importing tendermint node keys")
	if err := importTendermintNodeKey(nodeSet.Tendermint.HomeDir, tendermintNodeKey.PrivateKey); err != nil {
		return fmt.Errorf("failed to import the tendermint node private key")
	}

	log.Println("... verifying imported tendermint node")
	if err := verifyTendermintNode(nodeSet.Vega.BinaryPath, nodeSet.Tendermint.HomeDir, tendermintNodeKey.NodeID); err != nil {
		return fmt.Errorf("failed to verify imported tendermint node: %w", err)
	}

	log.Println("... adding tendermint wallet to the nodewallet")
	tmImportArgs := []string{
		"nodewallet", "import", "--force",
		"--home", nodeSet.Vega.HomeDir,
		"--chain", types.NodeWalletChainTypeTendermint,
		"--passphrase-file", nodeSet.Vega.NodeWalletPassFilePath,
		"--output", "json",
		"--tendermint-home", nodeSet.Tendermint.HomeDir,
	}
	if _, err := utils.ExecuteBinary(nodeSet.Vega.BinaryPath, tmImportArgs, nil); err != nil {
		return fmt.Errorf("failed to import tendermint to vega nodewallet: %w", err)
	}

	return nil
}
