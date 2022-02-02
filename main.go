package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/hashicorp/nomad/api"
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
	outputDir := "./testnet"
	prefix := "st-local"
	nodeDirPrefix := "node"
	tendermintNodePrefix := "tendermint-node"
	vegaNodePrefix := "vega-node"
	dataNodePrefix := "data-node"
	vegaBinaryPath := "/Users/karelmoravec/go/bin/vega"
	vegaDir := path.Join(outputDir, "vega")
	tendermintDir := path.Join(outputDir, "tendermint")
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

	vegaDirAbs, err := filepath.Abs(vegaDir)
	if err != nil {
		log.Fatal(err)
	}

	tendermintDirAbs, err := filepath.Abs(tendermintDir)
	if err != nil {
		log.Fatal(err)
	}

	tasks := []*api.Task{}

	for i := 0; i < nValidators+nNonValidators; i++ {
		tasks = append(tasks,
			&api.Task{
				Name:   fmt.Sprintf("vega-node-%d", i),
				Driver: "raw_exec",
				Config: map[string]interface{}{
					"command": vegaBinaryPath,
					"args": []string{
						"node",
						"--home", fmt.Sprintf("%s/node%d", vegaDirAbs, i),
						"--nodewallet-passphrase-file", fmt.Sprintf("%s/node%d/node-vega-wallet-pass.txt", vegaDirAbs, i),
					},
				},
				Resources: &api.Resources{
					CPU:      intPoint(500),
					MemoryMB: intPoint(256),
				},
			},
			&api.Task{
				Name:   fmt.Sprintf("tendermint-node-%d", i),
				Driver: "raw_exec",
				Config: map[string]interface{}{
					"command": vegaBinaryPath,
					"args": []string{
						"tm",
						"node",
						"--home", fmt.Sprintf("%s/node%d", tendermintDirAbs, i),
					},
				},
				Resources: &api.Resources{
					CPU:      intPoint(500),
					MemoryMB: intPoint(256),
				},
			},
		)
	}

	j := &api.Job{
		ID:          strPoint("test-vega-network-1"),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				Name:  strPoint("vega"),
				Tasks: tasks,
			},
		},
	}

	if _, err := nomadRunner.Run(j); err != nil {
		log.Fatal(err)
	}
}

func strPoint(s string) *string {
	return &s
}

func intPoint(i int) *int {
	return &i
}
