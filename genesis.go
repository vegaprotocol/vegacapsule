package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"
	"text/template"

	"code.vegaprotocol.io/vega/genesis"
	vgtm "code.vegaprotocol.io/vega/tendermint"
	"github.com/imdario/mergo"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
)

type SmartContract struct {
	Ethereum string `json:"Ethereum"`
	Vega     string `json:"Vega"`
}

type GenesisTemplateContext struct {
	Addresses map[string]SmartContract
	ChainID   string
	NetworkID string
}

func NewGenesisTemplateContext(chainID, networkID string, addressesJSON []byte) (*GenesisTemplateContext, error) {
	addrs := map[string]SmartContract{}

	if err := json.Unmarshal(addressesJSON, &addrs); err != nil {
		return nil, fmt.Errorf("could not parse json smart contract addresses", addressesJSON)
	}

	return &GenesisTemplateContext{
		ChainID:   chainID,
		NetworkID: networkID,
		Addresses: addrs,
	}, nil
}

func (gc GenesisTemplateContext) GetEthContractAddr(contract string) string {
	sc, ok := gc.Addresses[contract]
	if !ok {
		log.Fatalf("could not find Ethereum smart contract %q", contract)
	}

	if sc.Ethereum == "" {
		log.Fatalf("could not find Ethereum smart contract %q", contract)
	}

	return sc.Ethereum
}

func (gc GenesisTemplateContext) GetVegaContractID(contract string) string {
	sc, ok := gc.Addresses[contract]
	if !ok {
		log.Fatalf("could not find Vega smart contract %q", contract)
	}

	if sc.Vega == "" {
		log.Fatalf("could not find Vega smart contract %q", contract)
	}

	return strings.Replace(sc.Vega, "0x", "", 1)
}

type GenesisGenerator struct {
	template    *template.Template
	templateCtx *GenesisTemplateContext
}

func NewGenesisGenerator(genTemplate string, templateCtx *GenesisTemplateContext) (*GenesisGenerator, error) {
	tpl, err := template.New("genesis.json").Parse(genTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse genesis override: %w", err)
	}

	return &GenesisGenerator{
		template:    tpl,
		templateCtx: templateCtx,
	}, nil
}

func (g *GenesisGenerator) executeTemplate() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	if err := g.template.Execute(buff, g.templateCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// try to unmarshal to validate that JSON template is valida genesis document
	genDocOverride := tmtypes.GenesisDoc{}
	if err := tmjson.Unmarshal(buff.Bytes(), &genDocOverride); err != nil {
		return nil, fmt.Errorf("failed to parse genesis templated genesis: %w", err)
	}

	return buff.Bytes(), nil
}

func (g *GenesisGenerator) Generate(tendermintDir string, nodes []*vegaNode) error {
	templatedOverride, err := g.executeTemplate()
	if err != nil {
		return err
	}

	var genDoc *tmtypes.GenesisDoc
	var genState *genesis.GenesisState
	var nValidators int

	for _, node := range nodes {
		if node.Genesis == nil {
			continue
		}

		doc, state, err := genesis.GenesisFromJSON(node.Genesis.RawOutput)
		if err != nil {
			return fmt.Errorf("failed to get genesis from JSON: %w", err)
		}

		nValidators++

		if genState == nil {
			genState = state
			genDoc = doc
			continue
		}

		// Add validators to shared state
		for _, v := range state.Validators {
			genState.Validators[v.TmPubKey] = v
		}
	}

	// Nothing to do, we can stop here
	if nValidators == 0 {
		return nil
	}

	// TODO clean up this genesis merging mess...
	if err := vgtm.AddAppStateToGenesis(genDoc, genState); err != nil {
		return err
	}

	genDocBytes, err := tmjson.Marshal(genDoc)
	if err != nil {
		return err
	}

	b, err := mergeJSON(genDocBytes, templatedOverride)
	if err != nil {
		return fmt.Errorf("failed to override genesis json: %w", err)
	}

	mergedGenDoc, _, err := genesis.GenesisFromJSON(b)
	if err != nil {
		return fmt.Errorf("failed to get merged config from json: %w", err)
	}

	for i := 0; i < nValidators; i++ {
		genesisFilePath := path.Join(tendermintDir, fmt.Sprintf("node%d", i), "config", "genesis.json")

		if err := mergedGenDoc.SaveAs(genesisFilePath); err != nil {
			return fmt.Errorf("failed to save genesis file: %w", err)
		}
	}

	return nil
}

func mergeJSON(dst, src []byte) ([]byte, error) {
	var dstM, srcM map[string]interface{}

	if err := json.Unmarshal(dst, &dstM); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dst: %w", err)
	}
	if err := json.Unmarshal(src, &srcM); err != nil {
		return nil, fmt.Errorf("failed to unmarshal src: %w", err)
	}

	if err := mergo.MergeWithOverwrite(&dstM, srcM); err != nil {
		return nil, fmt.Errorf("failed to merge maps: %w", err)
	}

	b, err := json.Marshal(dstM)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal destination map: %w", err)
	}

	return b, nil
}
