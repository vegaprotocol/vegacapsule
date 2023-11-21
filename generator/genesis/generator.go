package genesis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"code.vegaprotocol.io/vega/core/genesis"
	vgtm "code.vegaprotocol.io/vega/core/tendermint"
	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	"github.com/Masterminds/sprig"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/imdario/mergo"
)

type updateGenesisOutput struct {
	RawOutput json.RawMessage
}

type Generator struct {
	vegaBinary  string
	template    *template.Template
	templateCtx *TemplateContext
}

func NewGenerator(conf *config.Config, templateRaw string) (*Generator, error) {
	tpl, err := template.New("genesis.json").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse genesis override: %w", err)
	}

	templateContext, err := NewTemplateContext(conf.Network.Ethereum.ChainID, conf.Network.Ethereum.NetworkID, []byte(*conf.Network.SmartContractsAddresses))
	if err != nil {
		return nil, err
	}

	return &Generator{
		vegaBinary:  *conf.VegaBinary,
		template:    tpl,
		templateCtx: templateContext,
	}, nil
}

func (g *Generator) ExecuteTemplate() (*bytes.Buffer, error) {
	buff := bytes.NewBuffer([]byte{})

	if err := g.template.Execute(buff, g.templateCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// try to unmarshal to validate that JSON template is valida genesis document
	genDocOverride := tmtypes.GenesisDoc{}
	if err := tmjson.Unmarshal(buff.Bytes(), &genDocOverride); err != nil {
		return nil, fmt.Errorf("failed to parse genesis templated genesis: %w", err)
	}

	return buff, nil
}

func (g *Generator) GenerateAndSave(chainID *string, validatorsSets []types.NodeSet, nonValidatorsSets []types.NodeSet, genValidators []tmtypes.GenesisValidator) error {
	genDoc, err := g.generate(validatorsSets, genValidators, chainID)
	if err != nil {
		return err
	}

	for _, ns := range append(validatorsSets, nonValidatorsSets...) {
		if err := genDoc.SaveAs(ns.Tendermint.GenesisFilePath); err != nil {
			return fmt.Errorf("failed to save genesis file: %w", err)
		}
	}

	return nil
}

func (g *Generator) Generate(validatorsSets []types.NodeSet, genValidators []tmtypes.GenesisValidator, chainID *string) (*bytes.Buffer, error) {
	genDoc, err := g.generate(validatorsSets, genValidators, chainID)
	if err != nil {
		return nil, err
	}

	tempFileName := "genesis.json"

	f, err := os.CreateTemp("", tempFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary genesis file %q: %w", tempFileName, err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	if err := genDoc.SaveAs(f.Name()); err != nil {
		return nil, fmt.Errorf("failed to save genesis file: %w", err)
	}

	buffOut := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffOut, f); err != nil {
		return nil, fmt.Errorf("failed to copy content of config file %q: %w", f.Name(), err)
	}

	return buffOut, nil
}

func (g *Generator) generate(nodeSets []types.NodeSet, genValidators []tmtypes.GenesisValidator, chainID *string) (*tmtypes.GenesisDoc, error) {
	templatedOverride, err := g.ExecuteTemplate()
	if err != nil {
		return nil, err
	}

	var genDoc *tmtypes.GenesisDoc
	var genState *genesis.State

	for _, ns := range nodeSets {
		updatedGenesis, err := g.updateGenesis(
			ns.Vega.HomeDir,
			ns.Tendermint.HomeDir,
			ns.Vega.NodeWalletPassFilePath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update genesis for %q from %q: %w", ns.Tendermint.HomeDir, ns.Vega.HomeDir, err)
		}

		doc, state, err := genesisFromJSON(updatedGenesis.RawOutput)
		if err != nil {
			return nil, fmt.Errorf("failed to get genesis from JSON: %w", err)
		}

		if chainID != nil {
			doc.ChainID = *chainID
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
		return nil, fmt.Errorf("failed to generate genesis for empty NodeSets")
	}

	// TODO should this be inside of template???
	genDoc.Validators = genValidators

	// TODO clean up this genesis merging mess...
	if err := vgtm.AddAppStateToGenesis(genDoc, genState); err != nil {
		return nil, err
	}

	genDocBytes, err := tmjson.Marshal(genDoc)
	if err != nil {
		return nil, err
	}

	b, err := mergeJSON(genDocBytes, templatedOverride.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to override genesis json: %w", err)
	}

	mergedGenDoc, _, err := genesisFromJSON(b)
	if err != nil {
		return nil, fmt.Errorf("failed to get merged config from json: %w", err)
	}

	return mergedGenDoc, nil
}

func (g *Generator) updateGenesis(
	vegaHomePath, tendermintHomePath,
	nodeWalletPhraseFile string,
) (*updateGenesisOutput, error) {
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

func genesisFromJSON(data []byte) (*tmtypes.GenesisDoc, *genesis.State, error) {
	doc := &tmtypes.GenesisDoc{}
	err := tmjson.Unmarshal(data, doc)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't unmarshal the genesis document: %w", err)
	}

	state := &genesis.State{}

	if len(doc.AppState) != 0 {
		if err := json.Unmarshal(doc.AppState, state); err != nil {
			return nil, nil, fmt.Errorf("couldn't unmarshal genesis state: %w", err)
		}
	}

	return doc, state, nil
}

func ConfigFilePath(nodeDir string) string {
	return filepath.Join(nodeDir, "config", "genesis.json")
}
