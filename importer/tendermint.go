package importer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"code.vegaprotocol.io/vegacapsule/config"
	tmgen "code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
)

func verifyTendermintNode(vegaBinary, tendermintHomePath, expectedNodeID string) error {
	expectedNodeID = strings.ToLower(expectedNodeID)

	out, err := utils.ExecuteBinary(vegaBinary, []string{"tm", "show-node-id", "--home", tendermintHomePath}, nil)
	if err != nil {
		return fmt.Errorf("failed to get the tendermint node id: %w", err)
	}

	tendermintNodeID := strings.ToLower(strings.Trim(string(out), " \t\n"))

	if tendermintNodeID != expectedNodeID {
		return fmt.Errorf("tendermint node is invalid: expected \"%s\", got \"%s\"", expectedNodeID, tendermintNodeID)
	}

	nodeKeyFileBytes, err := os.ReadFile(tmgen.NodeKeyFilePath(tendermintHomePath))
	if err != nil {
		return fmt.Errorf("failed to read content of the node_key.json file: %w", err)
	}

	nodeKeyData := &struct {
		ID string `json:"id"`
	}{}

	if err := json.Unmarshal(nodeKeyFileBytes, nodeKeyData); err != nil {
		return fmt.Errorf("failed to unmarshal node_key.json file: %w", err)
	}

	collectedNodeID := strings.ToLower(nodeKeyData.ID)
	if collectedNodeID != expectedNodeID {
		return fmt.Errorf("node_key.json file is not updated properly: ID in the file is invalid: expected `%v`, got `%v`", expectedNodeID, collectedNodeID)
	}

	return nil
}

func importTendermintNodeKey(tendermintHomePath, privateKey string, nodeID string) error {
	nodeKeyFilePath := tmgen.NodeKeyFilePath(tendermintHomePath)
	rootNode, err := dasel.NewFromFile(nodeKeyFilePath, "json")
	if err != nil {
		return fmt.Errorf("failed to load private validator file: %w", err)
	}
	if err := rootNode.Put(".priv_key.value", privateKey); err != nil {
		return fmt.Errorf("failed to update private key in the node_key.json file: %w", err)
	}

	if err := rootNode.Put(".id", nodeID); err != nil {
		return fmt.Errorf("failed to update node id in the node_key.json file: %w", err)
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

type tendermintNodeKeys struct {
	nodeKey      tendermintKey
	validatorKey tendermintKey
}

func decodeAndImportTendermintKeys(nodeSet types.NodeSet, tmValidatorPrivateKey string, tmNodePrivateKey string) (*tendermintNodeKeys, error) {
	log.Println("... preparing tendermint validator keys")
	tendermintValidatorKey, err := decodeBase64TendermintPrivateKey(tmValidatorPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private tendermint validator key: %w", err)
	}

	log.Println("... preparing tendermint node keys")
	tendermintNodeKey, err := decodeBase64TendermintPrivateKey(tmNodePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private tendermint node key: %w", err)
	}

	log.Println("... importing tendermint validator keys")
	err = importTendermintPrivateValidator(nodeSet.Tendermint.HomeDir, tendermintPrivateValidatorData{
		Address:    tendermintValidatorKey.NodeID,
		PublicKey:  tendermintValidatorKey.PublicKey,
		PrivateKey: tendermintValidatorKey.PrivateKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to import the tendermint validator private data")
	}

	log.Println("... importing tendermint node keys")
	if err := importTendermintNodeKey(nodeSet.Tendermint.HomeDir, tendermintNodeKey.PrivateKey, tendermintNodeKey.NodeID); err != nil {
		return nil, fmt.Errorf("failed to import the tendermint node private key")
	}

	log.Println("... verifying imported tendermint node")
	if err := verifyTendermintNode(nodeSet.Vega.BinaryPath, nodeSet.Tendermint.HomeDir, tendermintNodeKey.NodeID); err != nil {
		return nil, fmt.Errorf("failed to verify imported tendermint node: %w", err)
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
		return nil, fmt.Errorf("failed to import tendermint to vega nodewallet: %w", err)
	}

	return &tendermintNodeKeys{
		nodeKey:      *tendermintNodeKey,
		validatorKey: *tendermintValidatorKey,
	}, nil
}

func regenerateTendermintConfig(generatedNodeSets types.NodeSetMap, config *config.Config) error {
	errs := &utils.MultiError{}

	for _, nodeSetGroup := range config.Network.Nodes {
		log.Printf("updating tendermint configuration for the %s group\n", nodeSetGroup.Name)

		if nodeSetGroup.ConfigTemplates.Tendermint == nil {
			errs.Add(fmt.Errorf("tendermint template for the `%s` node set group is nil", nodeSetGroup.Name))
			continue
		}

		tmpl, err := tmgen.NewConfigTemplate(*nodeSetGroup.ConfigTemplates.Tendermint)
		if err != nil {
			errs.Add(fmt.Errorf("failed to create tendermint template for the `%s` node set group: %w", nodeSetGroup.Name, err))
			continue
		}

		gen, err := tmgen.NewConfigGenerator(config, generatedNodeSets.ToSlice())
		if err != nil {
			errs.Add(fmt.Errorf("failed to create tendermint generator for the `%s` node set group: %w", nodeSetGroup.Name, err))
			continue
		}

		for _, nodeSet := range generatedNodeSets {
			if nodeSet.GroupName != nodeSetGroup.Name {
				continue
			}
			log.Printf("... overwriting config.toml for the %s node set\n", nodeSet.Name)

			if err := gen.OverwriteConfig(nodeSet, tmpl); err != nil {
				errs.Add(fmt.Errorf("failed to update tendermint configuration for the `%s` node set: %w", nodeSet.Name, err))
				continue
			}
		}
		log.Println("done")
	}

	if errs.HasAny() {
		return errs
	}
	return nil
}
