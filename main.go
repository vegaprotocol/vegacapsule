package main

import (
	"log"
	"path"
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

// type Variable struct {
// 	Type  string `hcl:"type"`
// 	Value string `hcl:"value"`
// }

// type Config struct {
// 	IOMode  string        `hcl:"io_mode"`
// 	Service ServiceConfig `hcl:"service,block"`
// }

// type ServiceConfig struct {
// 	Protocol   string          `hcl:"protocol,label"`
// 	Type       string          `hcl:"type,label"`
// 	ListenAddr string          `hcl:"listen_addr"`
// 	Processes  []ProcessConfig `hcl:"process,block"`
// }

// type ProcessConfig struct {
// 	Type    string   `hcl:"type,label"`
// 	Command []string `hcl:"command"`
// }

// func simpleParser() {
// 	var config Config
// 	err := hclsimple.DecodeFile("config.hcl", nil, &config)
// 	if err != nil {
// 		log.Fatalf("Failed to load configuration: %s", err)
// 	}
// 	log.Printf("Configuration is %#v", config)

// 	for _, p := range config.Service.Processes {
// 		cmd := exec.Command(p.Command[0], p.Command[1:]...)
// 		cmd.Stderr = os.Stderr
// 		cmd.Stdout = os.Stdout

// 		log.Printf("cmd running %q\n", p.Type)

// 		if err := cmd.Run(); err != nil {
// 			log.Printf("cmd %q failed: %s \n", p.Type, err)
// 		}
// 	}
// }

// func contextParser() {
// 	parser := hclparse.NewParser()
// 	f, diags := parser.ParseHCLFile("config.hcl")

// 	if diags.HasErrors() {
// 		wr := hcl.NewDiagnosticTextWriter(
// 			os.Stdout,      // writer to send messages to
// 			parser.Files(), // the parser's file cache, for source snippets
// 			78,             // wrapping width
// 			true,           // generate colored/highlighted output
// 		)
// 		wr.WriteDiagnostics(diags)
// 		return
// 	}

// 	ctx := &hcl.EvalContext{
// 		Variables: map[string]cty.Value{
// 			"pid": cty.NumberIntVal(int64(os.Getpid())),
// 		},
// 	}

// 	var c Config
// 	moreDiags := gohcl.DecodeBody(f.Body, ctx, &c)
// 	diags = append(diags, moreDiags...)
// 	fmt.Printf("%#v\n", c)
// }

// TODO should come from config
const defaultGanacheMnemonic = "cherry manage trip absorb logic half number test shed logic purpose rifle"

func main() {

	nomadRunner, err := NewNomadRunner(nil)
	if err != nil {
		log.Fatal(err)
	}

	// ganacheCheck()
	//deregisterJobs(nomadClient)
	outputDir := "./testnet"
	prefix := "st-local"
	nodeDirPrefix := "node"
	tendermintNodePrefix := "tendermint-node"
	vegaNodePrefix := "vega-node"
	dataNodePrefix := "data-node"
	nodeMode := "validator"
	vegaBinaryPath := "/Users/karelmoravec/go/bin/vega"
	vegaDir := path.Join(outputDir, "vega")
	tendermintDir := path.Join(outputDir, "tendermint")
	chainID := "1440"
	networkID := "1441"

	nValidators := 2
	nNonValidators := 0

	if err := generateTendermintConfig(tendermintDir, prefix, nodeDirPrefix, tendermintNodePrefix, vegaNodePrefix, defaultTendermintOverride, nValidators, nNonValidators); err != nil {
		panic(err)
	}

	nodes, err := generateVegaConfig(vegaBinaryPath, vegaDir, tendermintDir, prefix, nodeDirPrefix, tendermintNodePrefix, vegaNodePrefix, dataNodePrefix, nodeMode, defaultVegaOverride, defaultGenesisOverride)
	if err != nil {
		panic(err)
	}

	tmplCtx, err := NewGenesisTemplateContext(chainID, networkID, []byte(defaultSmartContractsAddresses))
	if err != nil {
		log.Fatalf("failed create genesis template context: %w", err)
	}

	genGenerator, err := NewGenesisGenerator(defaultGenesisTemplate, tmplCtx)
	if err != nil {
		log.Fatalf("failed to crate new genesis generator: %w", err)
	}

	if err := genGenerator.Generate(tendermintDir, nodes); err != nil {
		log.Fatalf("failed to generate genesis: %w", err)
	}

	// nomadRunner.Run(job *api.Job)
}
