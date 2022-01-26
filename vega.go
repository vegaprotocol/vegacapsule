package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

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
	fmt.Println("init vega: ", pass)
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

type walletOutput struct {
	Name     string "json`name`"
	Version  int    "json`version`"
	FilePath string "json`filePath`"
}

type importVegaWalletOutput struct {
	Wallet walletOutput "json`wallet`"
}

func importVegaWallet(
	vegaBinaryPath string,
	homePath string,
	mnemonicFile string,
	walletPhraseFile string,
) (*importVegaWalletOutput, error) {
	args := []string{
		"wallet",
		"--output", "json",
		"--no-version-check",
		"import",
		"--home", homePath,
		"--recovery-phrase-file", mnemonicFile,
		"--passphrase-file", walletPhraseFile,
		"--wallet", "wallet_from_mnemonic",
	}

	log.Printf("Importing wallet: %v", args)

	wo := &importVegaWalletOutput{}
	if err := executeBinary(vegaBinaryPath, args, wo); err != nil {
		return nil, err
	}

	return wo, nil
}

type importVegaNodeWalletOutput struct {
	RegistryFilePath string "json`registryFilePath`"
	WalletFilePath   string "json`walletFilePath`"
}

func importVegaNodeWallet(
	vegaBinaryPath string,
	homePath string,
	nodeWalletPhraseFile string,
	walletPhraseFile string,
	walletPath string,
) (*importVegaNodeWalletOutput, error) {
	args := []string{
		"nodewallet",
		"--home", homePath,
		"--passphrase-file", nodeWalletPhraseFile,
		"import",
		"--output", "json",
		"--chain", "vega",
		"--wallet-passphrase-file", walletPhraseFile,
		"--wallet-path", walletPath,
	}

	log.Printf("Importing node wallet: %v", args)

	nwo := &importVegaNodeWalletOutput{}
	if err := executeBinary(vegaBinaryPath, args, nwo); err != nil {
		return nil, err
	}

	return nwo, nil
}

func initateVegaNode(
	vegaBinaryPath string,
	vegaDir string,
	prefix string,
	nodeDirPrefix string,
	tendermintNodePrefix string,
	vegaNodePrefix string,
	dataNodePrefix string,
	configOverride string,
	nodeMode string,
	ethKey keyPair,
	mnemonic walletMnemonic,
	id int,
) error {
	nodeDir := path.Join(vegaDir, fmt.Sprintf("node%d", id))

	if err := os.MkdirAll(nodeDir, os.ModePerm); err != nil {
		return err
	}

	// TODO These vars should come from config or be generated somehow....
	nodeWalletPass := "n0d3w4ll3t-p4ssphr4e3"
	vegaWalletPass := "w4ll3t-p4ssphr4e3"
	// chainEthereumWalletPass := "ch41nw4ll3t-3th3r3um-p4ssphr4e3"

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
		mnemonicFilePath := path.Join(nodeDir, "mnemonic.txt")
		walletPassFilePath := path.Join(nodeDir, "vega-wallet-pass.txt")
		nodeWalletPassFilePath := path.Join(nodeDir, "node-vega-wallet-pass.txt")

		if err := ioutil.WriteFile(mnemonicFilePath, []byte(mnemonic.Mnemonic), 0644); err != nil {
			return fmt.Errorf("failed to write mnemonic to file: %w", err)
		}

		if err := ioutil.WriteFile(walletPassFilePath, []byte(vegaWalletPass), 0644); err != nil {
			return fmt.Errorf("failed to write wallet passphrase to file: %w", err)
		}

		if err := ioutil.WriteFile(nodeWalletPassFilePath, []byte(nodeWalletPass), 0644); err != nil {
			return fmt.Errorf("failed to write node wallet passphrase to file: %w", err)
		}

		wo, err := importVegaWallet(vegaBinaryPath, nodeDir, mnemonicFilePath, walletPassFilePath)
		if err != nil {
			return fmt.Errorf("failed import vega wallet: %w", err)
		}

		// we need absolute vega wallet path to import it to node wallet
		absoluteWalletPath, err := filepath.Abs(wo.Wallet.FilePath)
		if err != nil {
			return err
		}

		nwo, err := importVegaNodeWallet(vegaBinaryPath, nodeDir, nodeWalletPassFilePath, walletPassFilePath, absoluteWalletPath)
		if err != nil {
			return fmt.Errorf("failed import vega node wallet: %w", err)
		}
		fmt.Println(nwo)

		// TODO import ethereum wallet here.. implement the geth in Go: "github.com/ethereum/go-ethereum/accounts/keystore"
		// generateEthereumWallet()
	}

	log.Printf("vega config initialised for node id %d, paths: %#v", id, vegaConfigs.Loader.ConfigFilePath())

	return nil
}

func generateVegaConfig(
	vegaBinaryPath string,
	vegaDir string,
	prefix string,
	nodeDirPrefix string,
	tendermintNodePrefix string,
	vegaNodePrefix string,
	dataNodePrefix string,
	nodeMode string,
	configOverride string,
) error {
	ethKeys, err := generateEthereumKeys(1)
	if err != nil {
		return err
	}

	mnemonics, err := generateWalletMnemonics(1, "DV")
	if err != nil {
		return err
	}

	return initateVegaNode(
		vegaBinaryPath,
		vegaDir,
		prefix,
		nodeDirPrefix,
		tendermintNodePrefix,
		vegaNodePrefix,
		dataNodePrefix,
		configOverride,
		nodeMode,
		ethKeys.Keys[0],
		mnemonics[0],
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
