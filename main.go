package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

var (
	dockerGanacheImage        = "trufflesuite/ganache-cli:v6.12.2"
	dockerCipipelineImage     = "ghcr.io/vegaprotocol/devops-infra/cipipeline:latest"
	dockerEefImage            = "vegaprotocol/ethereum-event-forwarder:$eef_version"
	dockerPytoolsImage        = "ghcr.io/vegaprotocol/devops-infra/pytools:docker"
	dockerSmartcontractsImage = "ghcr.io/vegaprotocol/devops-infra/smartcontracts:docker"
	dockerVegaImage           = "ghcr.io/vegaprotocol/vega/vega:$vega_version"
	dockerVegawalletImage     = "vegaprotocol/vegawallet:$vegawallet_version"
	dockerDatanodeImage       = "ghcr.io/vegaprotocol/data-node/data-node:$datanode_version"
	dockerVegatoolsImage      = "vegaprotocol/vegatools:$vegatools_version"
	dockerClefImage           = "ghcr.io/vegaprotocol/devops-infra/clef:latest"
)

// TODO should come from config
const defaultGanacheMnemonic = "cherry manage trip absorb logic half number test shed logic purpose rifle"

var (
	outputDir            = "/Users/karelmoravec/vega/vegacomposer/testnet"
	prefix               = "st-local"
	nodeDirPrefix        = "node"
	tendermintNodePrefix = "tendermint-node"
	vegaNodePrefix       = "vega-node"
	dataNodePrefix       = "data-node"
	vegaBinaryPath       = "/Users/karelmoravec/go/bin/vega"
	chainID              = "1440"
	networkID            = "1441"
	nValidators          = 2
	nNonValidators       = 0
)

func start(configFilePath string) error {
	config, err := ParseConfigFile(configFilePath)
	if err != nil {
		return err
	}

	fmt.Println("---- config:", config)
	return nil
	log.Println("starting network")

	vegaDir, err := filepath.Abs(path.Join(outputDir, "vega"))
	if err != nil {
		return err
	}

	tendermintDir, err := filepath.Abs(path.Join(outputDir, "tendermint"))
	if err != nil {
		return err
	}

	tendermintNodes, err := generateTendermintConfigs(tendermintDir, prefix, nodeDirPrefix, tendermintNodePrefix, vegaNodePrefix, defaultTendermintOverride, nValidators, nNonValidators)
	if err != nil {
		return err
	}

	vegaNodes, err := generateVegaConfig(vegaBinaryPath, vegaDir, prefix, tendermintNodes, nodeDirPrefix, tendermintNodePrefix, vegaNodePrefix, dataNodePrefix, defaultVegaOverride, defaultGenesisTemplate)
	if err != nil {
		return err
	}

	tmplCtx, err := NewGenesisTemplateContext(chainID, networkID, []byte(defaultSmartContractsAddresses))
	if err != nil {
		return fmt.Errorf("failed create genesis template context: %s", err)
	}

	genGenerator, err := NewGenesisGenerator(defaultGenesisTemplate, tmplCtx)
	if err != nil {
		return fmt.Errorf("failed to crate new genesis generator: %s", err)
	}

	if err := genGenerator.Generate(tendermintDir, vegaNodes); err != nil {
		return fmt.Errorf("failed to generate genesis: %s", err)
	}

	nomadRunner, err := NewNomadRunner(nil)
	if err != nil {
		return err
	}

	runner := NewRunner(nomadRunner)

	if err := runner.StartRawNetwork(vegaBinaryPath, tendermintNodes, vegaNodes); err != nil {
		return fmt.Errorf("failed to start nomad network: %s", err)
	}

	log.Println("starting network success")
	return nil
}

func stop() {
	log.Println("stopping network")
	nomadRunner, err := NewNomadRunner(nil)
	if err != nil {
		log.Fatal(err)
	}

	runner := NewRunner(nomadRunner)

	if err := runner.StopRawNetwork(); err != nil {
		log.Fatalf("failed to start nomad network: %s", err)
	}
	log.Println("stopping network success")
}

func destroy() {
	log.Println("destroying network")
	stop()

	if err := os.RemoveAll(outputDir); err != nil {
		log.Fatalf("failed to destroy network: %s", err)
	}
	log.Println("destroying network success")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'start'|'stop'|'destroy' subcommands")
		os.Exit(1)
	}

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	configFilePath := startCmd.String("config-path", "", "enable")

	arg := os.Args[1]
	switch arg {
	case "start":
		if err := startCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
		if err := start(*configFilePath); err != nil {
			log.Fatal(err)
		}
	case "stop":
		stop()
	case "destroy":
		destroy()
	default:
		log.Printf("unknown subcommand %s: expected 'start'|'stop'|'destroy' subcommands", arg)
		os.Exit(1)
	}

}

func strPoint(s string) *string {
	return &s
}

func intPoint(i int) *int {
	return &i
}
