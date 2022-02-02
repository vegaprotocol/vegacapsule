package main

import (
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

func start() {
	log.Println("starting network")

	outputDir := "./testnet"
	prefix := "st-local"
	nodeDirPrefix := "node"
	tendermintNodePrefix := "tendermint-node"
	vegaNodePrefix := "vega-node"
	dataNodePrefix := "data-node"
	vegaBinaryPath := "/Users/karelmoravec/go/bin/vega"

	vegaDir, err := filepath.Abs(path.Join(outputDir, "vega"))
	if err != nil {
		log.Fatal(err)
	}

	tendermintDir, err := filepath.Abs(path.Join(outputDir, "tendermint"))
	if err != nil {
		log.Fatal(err)
	}

	chainID := "1440"
	networkID := "1441"

	nValidators := 2
	nNonValidators := 0

	tendermintNodes, err := generateTendermintConfigs(tendermintDir, prefix, nodeDirPrefix, tendermintNodePrefix, vegaNodePrefix, defaultTendermintOverride, nValidators, nNonValidators)
	if err != nil {
		panic(err)
	}

	vegaNodes, err := generateVegaConfig(vegaBinaryPath, vegaDir, prefix, tendermintNodes, nodeDirPrefix, tendermintNodePrefix, vegaNodePrefix, dataNodePrefix, defaultVegaOverride, defaultGenesisTemplate)
	if err != nil {
		panic(err)
	}

	tmplCtx, err := NewGenesisTemplateContext(chainID, networkID, []byte(defaultSmartContractsAddresses))
	if err != nil {
		log.Fatalf("failed create genesis template context: %s", err)
	}

	genGenerator, err := NewGenesisGenerator(defaultGenesisTemplate, tmplCtx)
	if err != nil {
		log.Fatalf("failed to crate new genesis generator: %s", err)
	}

	if err := genGenerator.Generate(tendermintDir, vegaNodes); err != nil {
		log.Fatalf("failed to generate genesis: %s", err)
	}

	nomadRunner, err := NewNomadRunner(nil)
	if err != nil {
		log.Fatal(err)
	}

	runner := NewRunner(nomadRunner)

	if err := runner.StartRawNetwork(vegaBinaryPath, tendermintNodes, vegaNodes); err != nil {
		log.Fatalf("failed to start nomad network: %s", err)
	}
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
}

func main() {
	if len(os.Args) < 2 {
		log.Println("missing command")
		os.Exit(1)
	}
	arg := os.Args[1]
	switch arg {
	case "start":
		start()
	case "stop":
		stop()
	default:
		log.Printf("unknown command %s", arg)
		os.Exit(1)
	}

}

func strPoint(s string) *string {
	return &s
}

func intPoint(i int) *int {
	return &i
}
