package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/zannen/toml"

	"github.com/imdario/mergo"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vega/config"
)

type VegaTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	DataNodePrefix       string
	ETHEndpoint          string
	Type                 string
	NodeNumber           int
}

// ClientAddr = "tcp://{{.Prefix}}-{{.TendermintNodePrefix}}{{.NodeNumber}}:26657"

var defaultVegaOverride = `
[API]
	Port = 30{{.NodeNumber}}2
	[API.REST]
   		Port = 30{{.NodeNumber}}3

[Blockchain]
    [Blockchain.Tendermint]
        ClientAddr = "tcp://127.0.0.1:266{{.NodeNumber}}7"
        ServerAddr = "0.0.0.0"
		ServerPort = 266{{.NodeNumber}}8
	[Blockchain.Null]
		Port = 31{{.NodeNumber}}1

[EvtForward]
    Level = "Info"
    RetryRate = "1s"

[NodeWallet]
    [NodeWallet.ETH]
        Address = "{{.ETHEndpoint}}"

[Processor]
    [Processor.Ratelimit]
        Requests = 10000
        PerNBlocks = 1
{{if eq .Type "full"}}
[Broker]
    [Broker.Socket]
        Address = "{{.Prefix}}-{{.DataNodePrefix}}{{.NodeNumber}}"
        Port = 3005
        Enabled = true
{{end}}
`

type nodeMode string

const (
	NodeModeValidator nodeMode = "validator"
	NodeModeFull      nodeMode = "full"
)

type initateNodeOutput struct {
	ConfigFilePath           string `json:"configFilePath"`
	NodeWalletConfigFilePath string `json:"nodeWalletConfigFilePath"`
}

func initiateNode(
	vegaBinaryPath string,
	homePath string,
	nodeWalletPhraseFile string,
	nodeMode nodeMode,
) (*initateNodeOutput, error) {
	args := []string{
		"init",
		"--home", homePath,
		"--nodewallet-passphrase-file", nodeWalletPhraseFile,
		"--output", "json",
		string(nodeMode),
	}

	log.Printf("Initiating node %q wallet with: %v", nodeMode, args)

	out := &initateNodeOutput{}
	if _, err := executeBinary(vegaBinaryPath, args, out); err != nil {
		return nil, err
	}

	return out, nil
}

func overrideVegaConfig(
	tmplCtx VegaTemplateContext,
	configFilePath string,
	configOverride string,
) error {
	t, err := template.New("config.toml").Parse(configOverride)
	if err != nil {
		return fmt.Errorf("failed to parse config override: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	if err := t.Execute(buff, tmplCtx); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	overrideConfig := config.Config{}

	if _, err := toml.DecodeReader(buff, &overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	vegaConfig := config.NewDefaultConfig()
	if err := paths.ReadStructuredFile(configFilePath, &vegaConfig); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", configFilePath, err)
	}

	if err := mergo.Merge(&overrideConfig, vegaConfig); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := paths.WriteStructuredFile(configFilePath, overrideConfig); err != nil {
		return fmt.Errorf("failed to write configuration file at %s: %w", configFilePath, err)
	}

	return nil
}

type generateNodeWalletOutput struct {
	Mnemonic         string `json:"mnemonic,omitempty"`
	RegistryFilePath string `json:"registryFilePath"`
	WalletFilePath   string `json:"walletFilePath"`
}

type walletChainType string

const (
	NodeWalletChainTypeVega     walletChainType = "vega"
	NodeWalletChainTypeEthereum walletChainType = "ethereum"
)

func generateNodeWallet(
	vegaBinaryPath string,
	homePath string,
	nodeWalletPhraseFile string,
	walletPhraseFile string,
	walletType walletChainType,
) (*generateNodeWalletOutput, error) {
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
	if _, err := executeBinary(vegaBinaryPath, args, out); err != nil {
		return nil, err
	}

	return out, nil
}

type importNodeWalletOutput struct {
	RegistryFilePath string `json:"registryFilePath"`
	TendermintPubkey string `json:"tendermintPubkey"`
}

func generateTendermintNodeWallet(
	vegaBinaryPath string,
	homePath string,
	nodeWalletPhraseFile string,
	tendermintHomePath string,
) (*importNodeWalletOutput, error) {
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
	if _, err := executeBinary(vegaBinaryPath, args, nwo); err != nil {
		return nil, err
	}

	return nwo, nil
}

type updateGenesisOutput struct {
	RawOutput json.RawMessage
}

func updateGenesis(
	vegaBinaryPath string,
	vegaHomePath string,
	nodeWalletPhraseFile string,
	tendermintHomePath string,
) (*updateGenesisOutput, error) {
	args := []string{
		"genesis",
		"--home", vegaHomePath,
		"--passphrase-file", nodeWalletPhraseFile,
		"update",
		"--tm-home", tendermintHomePath,
		"--dry-run",
	}

	log.Printf("Updating genesis: %v", args)

	rawOut, err := executeBinary(vegaBinaryPath, args, nil)
	if err != nil {
		return nil, err
	}

	return &updateGenesisOutput{RawOutput: rawOut}, nil
}

type vegaNode struct {
	NodeMode               string
	NodeHome               nodeMode
	WalletPassFilePath     string
	NodeWalletPassFilePath string
	EthereumPassFilePath   string
	VegaWallet             *generateNodeWalletOutput
	EthereumWallet         *generateNodeWalletOutput
	TendermintWallet       *importNodeWalletOutput
	Genesis                *updateGenesisOutput
}

func initateVegaNode(
	vegaBinaryPath string,
	vegaDir string,
	tendermintNode tendermintNode,
	prefix string,
	nodeDirPrefix string,
	tendermintNodePrefix string,
	vegaNodePrefix string,
	dataNodePrefix string,
	configOverride string,
	id int,
) (*vegaNode, error) {
	nodeDir := path.Join(vegaDir, fmt.Sprintf("node%d", id))

	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return nil, err
	}

	// TODO These vars should come from config or be generated somehow....
	nodeWalletPass := "n0d3w4ll3t-p4ssphr4e3"
	vegaWalletPass := "w4ll3t-p4ssphr4e3"
	ethereumWalletPass := "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

	nodeWalletPassFilePath := path.Join(nodeDir, "node-vega-wallet-pass.txt")
	if err := ioutil.WriteFile(nodeWalletPassFilePath, []byte(nodeWalletPass), 0644); err != nil {
		return nil, fmt.Errorf("failed to write node wallet passphrase to file: %w", err)
	}

	nodeMode := NodeModeFull
	if tendermintNode.IsValidator {
		nodeMode = NodeModeValidator
	}

	initOut, err := initiateNode(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, nodeMode)
	if err != nil {
		return nil, fmt.Errorf("failed to initate vega node: %w", err)
	}

	tmplCtx := VegaTemplateContext{
		Prefix:               prefix,
		TendermintNodePrefix: tendermintNodePrefix,
		VegaNodePrefix:       vegaNodePrefix,
		DataNodePrefix:       dataNodePrefix,
		Type:                 string(nodeMode),
		ETHEndpoint:          "http://192.168.1.102:8545/", // TODO this should come from config...
		NodeNumber:           id,
	}

	if err := overrideVegaConfig(tmplCtx, initOut.ConfigFilePath, configOverride); err != nil {
		return nil, fmt.Errorf("failed to override Vega config: %w", err)
	}

	// TODO this condition should really come from config...
	if tendermintNode.IsValidator {
		walletPassFilePath := path.Join(nodeDir, "vega-wallet-pass.txt")
		nodeWalletPassFilePath := path.Join(nodeDir, "node-vega-wallet-pass.txt")
		ethereumPassFilePath := path.Join(nodeDir, "ethereum-vega-wallet-pass.txt")

		if err := ioutil.WriteFile(walletPassFilePath, []byte(vegaWalletPass), 0644); err != nil {
			return nil, fmt.Errorf("failed to write wallet passphrase to file: %w", err)
		}

		if err := ioutil.WriteFile(nodeWalletPassFilePath, []byte(nodeWalletPass), 0644); err != nil {
			return nil, fmt.Errorf("failed to write node wallet passphrase to file: %w", err)
		}

		if err := ioutil.WriteFile(ethereumPassFilePath, []byte(ethereumWalletPass), 0644); err != nil {
			return nil, fmt.Errorf("failed to write ethereum wallet passphrase to file: %w", err)
		}

		vegaOut, err := generateNodeWallet(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, walletPassFilePath, NodeWalletChainTypeVega)
		if err != nil {
			return nil, fmt.Errorf("failed to generate vega wallet: %w", err)
		}

		log.Printf("node wallet out: %#v", vegaOut)

		ethOut, err := generateNodeWallet(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, ethereumPassFilePath, NodeWalletChainTypeEthereum)
		if err != nil {
			return nil, fmt.Errorf("failed to generate vega wallet: %w", err)
		}

		log.Printf("ethereum wallet out: %#v", ethOut)

		tmtOut, err := generateTendermintNodeWallet(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, tendermintNode.Home)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tenderming wallet: %w", err)
		}

		log.Printf("tendermint wallet out: %#v", tmtOut)

		genesisOut, err := updateGenesis(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, tendermintNode.Home)
		if err != nil {
			return nil, fmt.Errorf("failed to update genesis file: %w", err)
		}

		log.Printf("updated genesis file: %q\n", tendermintNode.GenesisPath)

		return &vegaNode{
			NodeMode:               nodeDir,
			NodeHome:               nodeMode,
			WalletPassFilePath:     walletPassFilePath,
			NodeWalletPassFilePath: nodeWalletPassFilePath,
			EthereumPassFilePath:   ethereumPassFilePath,
			VegaWallet:             vegaOut,
			EthereumWallet:         ethOut,
			TendermintWallet:       tmtOut,
			Genesis:                genesisOut,
		}, nil
	}

	log.Printf("vega config initialized for node id %d, paths: %#v", id, initOut.ConfigFilePath)

	return &vegaNode{
		NodeMode: nodeDir,
		NodeHome: nodeMode,
	}, nil
}

func generateVegaConfig(
	vegaBinaryPath string,
	vegaDir string,
	prefix string,
	tendermintNodes []tendermintNode,
	nodeDirPrefix string,
	tendermintNodePrefix string,
	vegaNodePrefix string,
	dataNodePrefix string,
	configOverride string,
	genesisOverrideStr string,
) ([]*vegaNode, error) {
	var nodes []*vegaNode

	for i, tn := range tendermintNodes {
		node, err := initateVegaNode(
			vegaBinaryPath,
			vegaDir,
			tn,
			prefix,
			nodeDirPrefix,
			tendermintNodePrefix,
			vegaNodePrefix,
			dataNodePrefix,
			configOverride,
			i,
		)
		if err != nil {
			return nil, err
		}

		fmt.Printf("initiated vega node: %#v \n", node)

		nodes = append(nodes, node)
	}

	return nodes, nil
}
