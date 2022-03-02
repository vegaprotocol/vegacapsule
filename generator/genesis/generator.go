package genesis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	"code.vegaprotocol.io/vega/genesis"
	vgtm "code.vegaprotocol.io/vega/tendermint"
	"github.com/imdario/mergo"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/Masterminds/sprig"
)

type updateGenesisOutput struct {
	RawOutput json.RawMessage
}

type Generator struct {
	vegaBinary  string
	template    *template.Template
	templateCtx *TemplateContext
}

func NewGenerator(conf *config.Config) (*Generator, error) {
	tpl, err := template.New("genesis.json").Funcs(sprig.TxtFuncMap()).Parse(conf.Network.GenesisTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse genesis override: %w", err)
	}

	templateContext, err := NewTemplateContext(conf.Network.Ethereum.ChainID, conf.Network.Ethereum.NetworkID, []byte(defaultSmartContractsAddresses))
	if err != nil {
		return nil, err
	}

	return &Generator{
		vegaBinary:  conf.VegaBinary,
		template:    tpl,
		templateCtx: templateContext,
	}, nil
}

func (g *Generator) executeTemplate() ([]byte, error) {
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

func (g *Generator) updateGenesis(vegaHomePath, tendermintHomePath, nodeWalletPhraseFile string) (*updateGenesisOutput, error) {
	args := []string{
		"genesis",
		"--home", vegaHomePath,
		"--passphrase-file", nodeWalletPhraseFile,
		"update",
		"--tm-home", tendermintHomePath,
		"--dry-run",
	}

	log.Printf("Updating genesis with: %s %v", g.vegaBinary, args)

	rawOut, err := utils.ExecuteBinary(g.vegaBinary, args, nil)
	if err != nil {
		return nil, err
	}

	return &updateGenesisOutput{RawOutput: rawOut}, nil
}

func (g *Generator) Generate(validatorsSets []types.NodeSet, nonValidatorsSets []types.NodeSet, genValidators []tmtypes.GenesisValidator) error {
	templatedOverride, err := g.executeTemplate()
	if err != nil {
		return err
	}

	var genDoc *tmtypes.GenesisDoc
	var genState *genesis.GenesisState

	for _, ns := range validatorsSets {
		updatedGenesis, err := g.updateGenesis(ns.Vega.HomeDir, ns.Tendermint.HomeDir, ns.Vega.NodeWalletPassFilePath)
		if err != nil {
			return fmt.Errorf("failed to update genesis for %q from %q: %w", ns.Tendermint.HomeDir, ns.Vega.HomeDir, err)
		}

		doc, state, err := genesis.GenesisFromJSON(updatedGenesis.RawOutput)
		if err != nil {
			return fmt.Errorf("failed to get genesis from JSON: %w", err)
		}

		if genDoc == nil {
			genDoc = doc
			genState = state
			continue
		}

		// Add validators to shared state
		for _, v := range state.Validators {
			genState.Validators[v.TmPubKey] = v
		}
	}

	// Nothing to do, we can stop here
	if genDoc == nil {
		return nil
	}

	// TODO should this be inside of template???
	genDoc.Validators = genValidators

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

	for _, ns := range append(validatorsSets, nonValidatorsSets...) {
		if err := mergedGenDoc.SaveAs(ns.Tendermint.GenesisFilePath); err != nil {
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
