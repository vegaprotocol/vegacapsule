package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/zannen/toml"

	"github.com/imdario/mergo"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vega/config"
	"code.vegaprotocol.io/vega/config/encoding"
	"code.vegaprotocol.io/vega/nodewallets"
)

type VegaConfig struct {
	Loader                   *config.Loader
	NodeWalletConfigFilePath string
}

type VegaTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	DataNodePrefix       string
	ETHEndpoint          string
	Type                 string
	NodeNumber           int
}

var defaultVegaOverride = `
[Blockchain]
    [Blockchain.Tendermint]
        ClientAddr = "tcp://{{.Prefix}}-{{.TendermintNodePrefix}}{{.NodeNumber}}:26657"
        ServerAddr = "0.0.0.0"

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
{{if eq .Type "validator"}}
[Broker]
    [Broker.Socket]
        Address = "{{.Prefix}}-{{.DataNodePrefix}}{{.NodeNumber}}"
        Port = 3005
        Enabled = true
{{end}}
`

// copied from Vega core
func initVegaConfig(modeS, dir, pass string) (*VegaConfig, error) {
	mode, err := encoding.NodeModeFromString(modeS)
	if err != nil {
		return nil, err
	}

	vegaPaths := paths.New(dir)

	// a nodewallet will be required only for a validator node
	var nwRegistry *nodewallets.RegistryLoader
	if mode == encoding.NodeModeValidator {
		nwRegistry, err = nodewallets.NewRegistryLoader(vegaPaths, pass)
		if err != nil {
			return nil, err
		}
	}

	cfgLoader, err := config.InitialiseLoader(vegaPaths)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise configuration loader: %w", err)
	}

	configExists, err := cfgLoader.ConfigExists()
	if err != nil {
		return nil, fmt.Errorf("couldn't verify configuration presence: %w", err)
	}

	if configExists {
		cfgLoader.Remove()
	}

	cfg := config.NewDefaultConfig()
	cfg.NodeMode = mode

	if err := cfgLoader.Save(&cfg); err != nil {
		return nil, fmt.Errorf("couldn't save configuration file: %w", err)
	}

	return &VegaConfig{
		Loader:                   cfgLoader,
		NodeWalletConfigFilePath: nwRegistry.RegistryFilePath(),
	}, nil
}

func overrideVegaConfig(
	tmplCtx VegaTemplateContext,
	loader *config.Loader,
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

	cfg := config.Config{}

	if _, err := toml.DecodeReader(buff, &cfg); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	vegaConfig, err := loader.Get()
	if err != nil {
		return fmt.Errorf("failed to get generated config: %w", err)
	}

	if err := mergo.Merge(&cfg, vegaConfig); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := loader.Save(&cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
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
	if err := executeBinary(vegaBinaryPath, args, out); err != nil {
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
	if err := executeBinary(vegaBinaryPath, args, nwo); err != nil {
		return nil, err
	}

	return nwo, nil
}

func updateGenesisFile(
	vegaBinaryPath string,
	vegaHomePath string,
	nodeWalletPhraseFile string,
	tendermintHomePath string,
) (*importNodeWalletOutput, error) {
	args := []string{
		"genesis",
		"--home", vegaHomePath,
		"--passphrase-file", nodeWalletPhraseFile,
		"update",
		"--tm-home", tendermintHomePath,
	}

	log.Printf("Updating genesis file: %v", args)

	nwo := &importNodeWalletOutput{}
	if err := executeBinary(vegaBinaryPath, args, nwo); err != nil {
		return nil, err
	}

	return nwo, nil
}

type vegaNode struct {
	NodeMode         string
	VegaWallet       *generateNodeWalletOutput
	EthereumWallet   *generateNodeWalletOutput
	TendermintWallet *importNodeWalletOutput
}

func initateVegaNode(
	vegaBinaryPath string,
	vegaDir string,
	tendermintDir string,
	prefix string,
	nodeDirPrefix string,
	tendermintNodePrefix string,
	vegaNodePrefix string,
	dataNodePrefix string,
	configOverride string,
	nodeMode string,
	id int,
) error {
	nodeDir := path.Join(vegaDir, fmt.Sprintf("node%d", id))
	tendermintNodeDir := path.Join(tendermintDir, fmt.Sprintf("node%d", id))

	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return err
	}

	// TODO These vars should come from config or be generated somehow....
	nodeWalletPass := "n0d3w4ll3t-p4ssphr4e3"
	vegaWalletPass := "w4ll3t-p4ssphr4e3"
	ethereumWalletPass := "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

	vegaConfigs, err := initVegaConfig(nodeMode, nodeDir, nodeWalletPass)
	if err != nil {
		return err
	}

	tmplCtx := VegaTemplateContext{
		Prefix:               prefix,
		TendermintNodePrefix: tendermintNodePrefix,
		VegaNodePrefix:       vegaNodePrefix,
		DataNodePrefix:       dataNodePrefix,
		Type:                 nodeMode,
		ETHEndpoint:          "tcp://rubbish.com",
		NodeNumber:           id,
	}

	if err := overrideVegaConfig(tmplCtx, vegaConfigs.Loader, configOverride); err != nil {
		return fmt.Errorf("failed to override Vega config: %w", err)
	}

	// TODO this condition should really come from config...
	if nodeMode == "validator" {
		walletPassFilePath := path.Join(nodeDir, "vega-wallet-pass.txt")
		nodeWalletPassFilePath := path.Join(nodeDir, "node-vega-wallet-pass.txt")
		ethereumPassFilePath := path.Join(nodeDir, "ethereum-vega-wallet-pass.txt")

		if err := ioutil.WriteFile(walletPassFilePath, []byte(vegaWalletPass), 0644); err != nil {
			return fmt.Errorf("failed to write wallet passphrase to file: %w", err)
		}

		if err := ioutil.WriteFile(nodeWalletPassFilePath, []byte(nodeWalletPass), 0644); err != nil {
			return fmt.Errorf("failed to write node wallet passphrase to file: %w", err)
		}

		if err := ioutil.WriteFile(ethereumPassFilePath, []byte(ethereumWalletPass), 0644); err != nil {
			return fmt.Errorf("failed to write ethereum wallet passphrase to file: %w", err)
		}

		vegaOut, err := generateNodeWallet(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, walletPassFilePath, NodeWalletChainTypeVega)
		if err != nil {
			return fmt.Errorf("failed to generate vega wallet: %w", err)
		}

		log.Printf("node wallet out: %#v", vegaOut)

		ethOut, err := generateNodeWallet(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, ethereumPassFilePath, NodeWalletChainTypeEthereum)
		if err != nil {
			return fmt.Errorf("failed to generate vega wallet: %w", err)
		}

		log.Printf("ethereum wallet out: %#v", ethOut)

		genTmtOut, err := generateTendermintNodeWallet(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, tendermintNodeDir)
		if err != nil {
			return fmt.Errorf("failed to generate tenderming wallet: %w", err)
		}

		log.Printf("tendermint wallet out: %#v", genTmtOut)

		// _, err = updateGenesisFile(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, tendermintNodeDir)
		// if err != nil {
		// 	return fmt.Errorf("failed to update genesis file: %w", err)
		// }

		// log.Printf("updated genesis file wallet out: %#v", genTmtOut)
	}

	log.Printf("vega config initialised for node id %d, paths: %#v", id, vegaConfigs.Loader.ConfigFilePath())

	return nil
}

func generateVegaConfig(
	vegaBinaryPath string,
	vegaDir string,
	tendermintDir string,
	prefix string,
	nodeDirPrefix string,
	tendermintNodePrefix string,
	vegaNodePrefix string,
	dataNodePrefix string,
	nodeMode string,
	configOverride string,
) error {
	return initateVegaNode(
		vegaBinaryPath,
		vegaDir,
		tendermintDir,
		prefix,
		nodeDirPrefix,
		tendermintNodePrefix,
		vegaNodePrefix,
		dataNodePrefix,
		configOverride,
		nodeMode,
		0,
	)

	// out, err := generateWalletMnemonics(3, "DV")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// b, err := json.Marshal(out)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(b))

	return nil
}
