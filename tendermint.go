package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/spf13/viper"

	tmcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/p2p"
)

type TendermingTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	NodeNumber           int
	NodesCount           int
	NodeIDs              []string
}

// proxy-app = "tcp://{{.Prefix}}-{{.VegaNodePrefix}}{{.NodeNumber}}:26658"

var defaultTendermintOverride = `
log_level = "error"

proxy_app = "tcp://127.0.0.1:266{{.NodeNumber}}8"
moniker = "{{.Prefix}}-{{.TendermintNodePrefix}}"

[rpc]
laddr = "tcp://0.0.0.0:266{{.NodeNumber}}7"
unsafe = true

[p2p]
laddr = "tcp://0.0.0.0:266{{.NodeNumber}}6"
addr_book_strict = true
max_packet_msg_payload_size = 4096
pex = false
allow_duplicate_ip = false
persistent_peers = "{{range $i, $v := .NodeIDs}}{{if ne $i 0}},{{end}}{{$v}}@127.0.0.1:266{{$i}}6{{end}}"

[mempool]
size = 10000
cache_size = 20000

[consensus]
skip_timeout_commit = false
`

type tendermintNode struct {
	Home        string
	ConfigPath  string
	GenesisPath string
	IsValidator bool
}

var defaultTemplateFuncs = template.FuncMap{
	"loop": func(n int) []int {
		var arr = make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = i
		}
		return arr
	},
}

// TODO perhaps find more flexible way of doing this??
func getNodeIDs(outputDir, nodeDirPrefix string, nValidators int) ([]string, error) {
	config := cfg.DefaultConfig()
	nodeIDs := make([]string, 0, nValidators)
	for i := 0; i < nValidators; i++ {
		nodeDir := filepath.Join(outputDir, fmt.Sprintf("%s%d", nodeDirPrefix, i))
		config.SetRoot(nodeDir)
		nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
		if err != nil {
			return nil, err
		}

		nodeIDs = append(nodeIDs, string(nodeKey.ID()))
	}

	return nodeIDs, nil
}

func generateTendermintConfigs(
	outputDir string,
	prefix string,
	nodeDirPrefix string,
	tendermintNodePrefix string,
	vegaNodePrefix string,
	configOverride string,
	nValidators int,
	nNonValidators int,
) ([]tendermintNode, error) {
	t, err := template.New("config.toml").Funcs(defaultTemplateFuncs).Parse(configOverride)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config override: %w", err)
	}

	os.Args = []string{"",
		"--o", outputDir,
		"--v", strconv.Itoa(nValidators),
		"--n", strconv.Itoa(nNonValidators),
		"--node-dir-prefix", nodeDirPrefix,
		"--hostname-prefix", fmt.Sprintf("%s-%s", prefix, tendermintNodePrefix),
		"--populate-persistent-peers",
	}

	log.Printf("Calling Tendermint testnet command with: %s", os.Args[1:])

	if err := tmcmd.TestnetFilesCmd.Execute(); err != nil {
		return nil, fmt.Errorf("failed to execute tendermint config generation: %w", err)
	}

	nodeIDs, err := getNodeIDs(outputDir, nodeDirPrefix, nValidators+nNonValidators)
	if err != nil {
		return nil, fmt.Errorf("failed to get Tendermint node ids: %w", err)
	}

	config := cfg.DefaultConfig()
	buff := bytes.NewBuffer([]byte{})

	tct := TendermingTemplateContext{
		Prefix:               prefix,
		TendermintNodePrefix: tendermintNodePrefix,
		VegaNodePrefix:       vegaNodePrefix,
		NodeNumber:           0,
		NodesCount:           nValidators + nNonValidators,
		NodeIDs:              nodeIDs,
	}

	nodes := make([]tendermintNode, 0, nValidators+nNonValidators)

	// Overwrite default config.
	for i := 0; i < nValidators+nNonValidators; i++ {
		homeDir := filepath.Join(outputDir, fmt.Sprintf("%s%d", nodeDirPrefix, i))
		configFilePath := filepath.Join(homeDir, "config", "config.toml")
		genesisFilePath := filepath.Join(homeDir, "config", "genesis.json")

		viper.SetConfigFile(configFilePath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file %q: %w", configFilePath, err)
		}

		buff.Reset()
		tct.NodeNumber = i
		err := t.Execute(buff, tct)
		if err != nil {
			return nil, fmt.Errorf("failed to execute template: %w", err)
		}

		if err := viper.MergeConfig(buff); err != nil {
			return nil, fmt.Errorf("failed to merge config override with config file %q: %w", configFilePath, err)
		}
		if err := viper.Unmarshal(config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal merged config file %q: %w", configFilePath, err)
		}
		if err := config.ValidateBasic(); err != nil {
			return nil, fmt.Errorf("failed to validated merged config file %q: %w", configFilePath, err)
		}

		// config.P2P.PersistentPeers

		config.SetRoot(homeDir)
		cfg.WriteConfigFile(configFilePath, config)

		var isValidator bool
		if i < nValidators {
			isValidator = true
		}

		nodes = append(nodes, tendermintNode{
			Home:        homeDir,
			ConfigPath:  configFilePath,
			GenesisPath: genesisFilePath,
			IsValidator: isValidator,
		})
	}

	return nodes, nil
}
