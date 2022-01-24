package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"

	tmcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	cfg "github.com/tendermint/tendermint/config"
)

type TendermingTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	NodeNumber           int
}

var defaultTendermintOverride = `
log_level = "error"

proxy-app = "tcp://{{.Prefix}}-{{.VegaNodePrefix}}{{.NodeNumber}}:26658"
moniker = "{{.Prefix}}-{{.TendermintNodePrefix}}"

[rpc]
laddr = "tcp://0.0.0.0:26657"
unsafe = true

[p2p]
laddr = "tcp://0.0.0.0:26656"
addr-book-strict = true
max-packet-msg-payload-size = 4096
pex = false
allow-duplicate-ip = false

[mempool]
size = 10000
cache-size = 20000

[consensus]
skip-timeout-commit = false
`

func generateTendermintConfig(
	outputDir string,
	prefix string,
	nodeDirPrefix string,
	tendermintNodePrefix string,
	vegaNodePrefix string,
	configOverride string,
	nValidators int,
	nNonValidators int,
) error {
	t, err := template.New("config.toml").Parse(configOverride)
	if err != nil {
		return fmt.Errorf("failed to parse config override: %w", err)
	}

	os.Args = []string{"",
		"--o", outputDir,
		"--v", strconv.Itoa(nValidators),
		"--n", strconv.Itoa(1),
		"--node-dir-prefix", nodeDirPrefix,
		"--hostname-prefix", fmt.Sprintf("%s-%s", prefix, tendermintNodePrefix),
		"--populate-persistent-peers",
	}

	log.Printf("Calling Tendermint testnet command with: %s", os.Args[1:])

	if err := tmcmd.TestnetFilesCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute tendermint config generation: %w", err)
	}

	config := cfg.DefaultConfig()
	buff := bytes.NewBuffer([]byte{})

	tct := TendermingTemplateContext{
		Prefix:               prefix,
		TendermintNodePrefix: tendermintNodePrefix,
		VegaNodePrefix:       vegaNodePrefix,
		NodeNumber:           0,
	}

	// Overwrite default config.
	for i := 0; i < nValidators+nNonValidators; i++ {
		nodeDir := filepath.Join(outputDir, fmt.Sprintf("%s%d", nodeDirPrefix, i))
		configFile := filepath.Join(nodeDir, "config", "config.toml")

		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config file %q: %w", configFile, err)
		}

		buff.Reset()
		tct.NodeNumber = i
		err := t.Execute(buff, tct)
		if err != nil {
			panic(err)
		}

		if err := viper.MergeConfig(buff); err != nil {
			return fmt.Errorf("failed to merge config override with config file %q: %w", configFile, err)
		}
		if err := viper.Unmarshal(config); err != nil {
			return fmt.Errorf("failed to unmarshal merged config file %q: %w", configFile, err)
		}
		if err := config.ValidateBasic(); err != nil {
			return fmt.Errorf("failed to validated merged config file %q: %w", configFile, err)
		}

		config.SetRoot(nodeDir)
		cfg.WriteConfigFile(nodeDir, config)

		// if err := ; err != nil {
		// 	return fmt.Errorf("failed to write config file %q: %w", configFile, err)
		// }
	}

	return nil
}
