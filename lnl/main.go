package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"

	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
)

const (
	genesisSourceURL    = "https://raw.githubusercontent.com/vegaprotocol/networks/master/mainnet1/genesis.json"
	genesisTempOutFile  = "/Users/daniel/www/vega/system-tests/vegacapsule/net_configs/mainnet/genesis.orig.json"
	genesisFinalOutFile = "/Users/daniel/www/vega/system-tests/vegacapsule/net_configs/mainnet/genesis.json"
)

const (
	vegatoolsPath             = "vegatools"
	vegacapsulePath           = "vegacapsule"
	checkpointFile            = "/Users/daniel/www/vega/checkpoint-store/Mainnet/20220914171226-1681774-c8140e6d43d050ee8cf8f79578372d3e42e25202485ca8b6ee3f127cd872d337.cp"
	updatedCheckpointJSONFile = "checkpoint.json"
)

type OverrideAction string

const (
	OverrideRemove OverrideAction = "Remove"
	OverrideUpdate OverrideAction = "Update"
)

type StructuredFileOverride struct {
	Selector string
	Action   OverrideAction
	Content  string
}

type StructuredFileOverrides []StructuredFileOverride

func OverrideMainnetGenesis() error {
	overrides := []StructuredFileOverride{
		{
			Selector: ".validators",
			Action:   OverrideRemove,
		},
		{
			Selector: ".app_state.validators",
			Action:   OverrideRemove,
		},
		{
			Selector: ".app_state.network_parameters.market\\.monitor\\.price\\.updateFrequency",
			Action:   OverrideRemove,
		},
		{
			Selector: ".app_state.checkpoint",
			Action:   OverrideRemove,
		},
	}

	if err := UpdateStructuredFile("json", genesisTempOutFile, genesisFinalOutFile, overrides); err != nil {
		return fmt.Errorf("failed to update genesis: %w", err)
	}

	return nil
}

func UpdateStructuredFile(fileType, inputFilePath, outputFilePath string, overrides StructuredFileOverrides) error {
	rootNode, err := dasel.NewFromFile(inputFilePath, fileType)
	if err != nil {
		return fmt.Errorf("failed to load the %s file: %w", inputFilePath, err)
	}

	// TODO: Use multierror
	for _, override := range overrides {
		switch override.Action {
		case OverrideRemove:
			if err := rootNode.Delete(override.Selector); err != nil {
				return fmt.Errorf("failed to update structured file: %w", err)
			}
		case OverrideUpdate:
			if err := rootNode.Put(override.Selector, override.Content); err != nil {
				return fmt.Errorf("failed to update structured file: %w", err)
			}
		default:
			return fmt.Errorf("action %s not supported", override.Action)
		}
	}

	if err := rootNode.WriteToFile(outputFilePath, fileType, []storage.ReadWriteOption{}); err != nil {
		return fmt.Errorf("failed to write the node_key.json file: %w", err)
	}

	return nil
}

func DownloadGenesis() error {
	if err := downloadFile(genesisTempOutFile, genesisSourceURL); err != nil {
		return fmt.Errorf("failed to download genesis: %w", err)
	}

	return nil
}

func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create destination file for genesis: %s", err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download given file: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download given file: bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file from buffer into destination file: %w", err)
	}

	return nil
}

type validatorInfo struct {
	Tendermint struct {
		ValidatorPublicKey string
		NodeID             string
	}
	NodeWalletInfo struct {
		EthereumAddress     string
		VegaWalletPublicKey string
	}
}

type validatorsInfo []validatorInfo

type mainnetCheckpointData struct {
	Block struct {
		Height string `json:"height"`
	} `json:"block"`
	Validators struct {
		ValidatorState []struct {
			ValidatorUpdate struct {
				NodeID          string `json:"nodeId"`
				VegaPubKey      string `json:"vegaPubKey"`
				EthereumAddress string `json:"ethereumAddress"`
				TmPubKey        string `json:"tmPubKey"`
			} `json:"validatorUpdate"`
		} `json:"validatorState"`
	} `json:"validators"`
}

type ValidatorData struct {
	NodeID              string
	VegaPubKey          string
	EthereumAddress     string
	TendermintPublicKey string
}

type ValidatorReplacement struct {
	Index   int
	OldData ValidatorData
	NewData *ValidatorData
}

func UpdateMainnetCheckpoint() error {
	cpExt := path.Ext(checkpointFile)
	checkpointJsonFilePath := checkpointFile[0:len(checkpointFile)-len(cpExt)] + ".json"

	if _, err := utils.ExecuteBinary(vegatoolsPath, []string{
		"checkpoint",
		"--file", checkpointFile,
		"--out", checkpointJsonFilePath,
	}, nil); err != nil {
		return fmt.Errorf("failed to convert binary checkpoint into JSON: %w", err)
	}

	mainnetCheckpoint := mainnetCheckpointData{}
	checkpointJSONBytes, err := ioutil.ReadFile(checkpointJsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read checkpoint JSON file: %w", err)
	}

	if err := json.Unmarshal(checkpointJSONBytes, &mainnetCheckpoint); err != nil {
		return fmt.Errorf("failed to unmarshal JSON checkpoint: %w", err)
	}

	// TODO: Update validators
	generatedValidators := validatorsInfo{}
	if _, err := utils.ExecuteBinary(vegacapsulePath, []string{
		"nodes",
		"ls-validators",
	}, &generatedValidators); err != nil {
		return fmt.Errorf("failed to get generated validators from vegacapsule: %w", err)
	}

	validatorsToReplace := []ValidatorReplacement{}
	for index, validator := range mainnetCheckpoint.Validators.ValidatorState {
		replacement := ValidatorReplacement{
			Index: index,
			OldData: ValidatorData{
				NodeID:              validator.ValidatorUpdate.NodeID,
				VegaPubKey:          validator.ValidatorUpdate.VegaPubKey,
				EthereumAddress:     validator.ValidatorUpdate.EthereumAddress,
				TendermintPublicKey: validator.ValidatorUpdate.TmPubKey,
			},
		}

		if len(generatedValidators) > index {
			replacement.NewData = &ValidatorData{
				NodeID:              generatedValidators[index].Tendermint.NodeID,
				TendermintPublicKey: generatedValidators[index].Tendermint.ValidatorPublicKey,
				VegaPubKey:          generatedValidators[index].NodeWalletInfo.VegaWalletPublicKey,
				EthereumAddress:     generatedValidators[index].NodeWalletInfo.EthereumAddress,
			}
		}

		validatorsToReplace = append(validatorsToReplace, replacement)
	}

	fmt.Printf("%v\n\n", validatorsToReplace)
	fmt.Printf("%v", generatedValidators)

	input, err := ioutil.ReadFile(checkpointJsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read json checkpoint file: %w", err)
	}

	for _, validator := range validatorsToReplace {
		if validator.NewData == nil {
			continue
		}

		input = bytes.Replace(input, []byte(validator.OldData.NodeID), []byte(validator.NewData.NodeID), -1)
		input = bytes.Replace(input, []byte(validator.OldData.TendermintPublicKey), []byte(validator.NewData.TendermintPublicKey), -1)
		input = bytes.Replace(input, []byte(validator.OldData.VegaPubKey), []byte(validator.NewData.VegaPubKey), -1)
		input = bytes.Replace(input, []byte(validator.OldData.EthereumAddress), []byte(validator.NewData.EthereumAddress), -1)
	}

	if err = ioutil.WriteFile(updatedCheckpointJSONFile, input, 0666); err != nil {
		return fmt.Errorf("failed to save updated mainnet checkpoint: %w", err)
	}

	checkpointOut, err := utils.ExecuteBinary(vegatoolsPath, []string{
		"checkpoint",
		"--generate",
		"--file", checkpointJsonFilePath,
		"--out", "checkpoint-tmp.cp",
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to convert json checkpoint into binary: %w", err)
	}

	newCheckpointHashRegex := regexp.MustCompile("[0-9a-f]{64}")
	newCheckpointHash := newCheckpointHashRegex.FindString(string(checkpointOut))

	newCheckpointName := fmt.Sprintf("%s-%s-%s",
		time.Now().Format("20060102030405"),
		mainnetCheckpoint.Block.Height,
		newCheckpointHash,
	)
	if err := os.Rename(checkpointFile, newCheckpointName); err != nil {
		return fmt.Errorf("failed to rename temporary binary checkpoint: %w", err)
	}

	return nil
}

func main() {
	if err := DownloadGenesis(); err != nil {
		panic(fmt.Errorf("failed to download genesis: %w", err))
	}

	if err := OverrideMainnetGenesis(); err != nil {
		panic(fmt.Errorf("failed to override mainnet genesis: %w", err))
	}

	if err := UpdateMainnetCheckpoint(); err != nil {
		panic(fmt.Errorf("failed to override mainnet genesis: %w", err))
	}
}
