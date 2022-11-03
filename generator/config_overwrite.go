package generator

import (
	"fmt"
	"log"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/datanode"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/generator/visor"
	"code.vegaprotocol.io/vegacapsule/types"
)

type configOverride struct {
	tendermintTmpl *template.Template
	vegaTmpl       *template.Template
	dataNodeTmpl   *template.Template
	visorRunTmpl   *template.Template
	visorConfTmpl  *template.Template
	gen            *Generator
}

func newConfigOverride(gen *Generator, n config.NodeConfig) (*configOverride, error) {
	var err error

	var tendermintTmpl *template.Template
	if n.ConfigTemplates.Tendermint != nil {
		tendermintTmpl, err = tendermint.NewConfigTemplate(*n.ConfigTemplates.Tendermint)
		if err != nil {
			return nil, err
		}
	}

	var vegaTmpl *template.Template
	if n.ConfigTemplates.Vega != nil {
		vegaTmpl, err = vega.NewConfigTemplate(*n.ConfigTemplates.Vega)
		if err != nil {
			return nil, err
		}
	}

	var dataNodeTmpl *template.Template
	if n.UseDataNode && n.ConfigTemplates.DataNode != nil {
		dataNodeTmpl, err = datanode.NewConfigTemplate(*n.ConfigTemplates.DataNode)
		if err != nil {
			return nil, err
		}
	}

	var visorRunTmpl *template.Template
	if n.VisorBinary != "" && n.ConfigTemplates.VisorRunConf != nil {
		visorRunTmpl, err = visor.NewConfigTemplate(*n.ConfigTemplates.VisorRunConf)
		if err != nil {
			return nil, err
		}
	}

	var visorConfTmpl *template.Template
	if n.VisorBinary != "" && n.ConfigTemplates.VisorConf != nil {
		visorConfTmpl, err = visor.NewConfigTemplate(*n.ConfigTemplates.VisorConf)
		if err != nil {
			return nil, err
		}
	}

	return &configOverride{
		tendermintTmpl: tendermintTmpl,
		vegaTmpl:       vegaTmpl,
		dataNodeTmpl:   dataNodeTmpl,
		visorRunTmpl:   visorRunTmpl,
		visorConfTmpl:  visorConfTmpl,
		gen:            gen,
	}, nil
}

func (co *configOverride) Overwrite(nc config.NodeConfig, ns types.NodeSet, fc *types.Faucet) error {
	if co.tendermintTmpl != nil {
		log.Printf("Overwriting Tendermint config for nodeset %s", ns.Name)
		if err := co.gen.tendermintGen.OverwriteConfig(ns, co.tendermintTmpl); err != nil {
			return fmt.Errorf("failed to overwrite Tendermit config for id %d: %w", ns.Index, err)
		}
	}
	if co.vegaTmpl != nil {
		log.Printf("Overwriting Vega config for nodeset %s", ns.Name)
		if err := co.gen.vegaGen.OverwriteConfig(ns, fc, co.vegaTmpl); err != nil {
			return fmt.Errorf("failed to overwrite Vega config for id %d: %w", ns.Index, err)
		}
	}
	if ns.DataNode != nil && co.dataNodeTmpl != nil {
		log.Printf("Overwriting Data Node config for nodeset %s", ns.Name)
		if err := co.gen.dataNodeGen.OverwriteConfig(ns, co.dataNodeTmpl); err != nil {
			return fmt.Errorf("failed to overwrite Data Node config for id %d: %w", ns.Index, err)
		}
	}
	if ns.Visor != nil && co.visorConfTmpl != nil {
		log.Printf("Overwriting Visor config for nodeset %s", ns.Name)
		if err := co.gen.visorGen.OverwriteConfig(ns, co.visorConfTmpl); err != nil {
			return fmt.Errorf("failed to overwrite Visor config for id %d: %w", ns.Index, err)
		}
	}
	if ns.Visor != nil && co.visorRunTmpl != nil {
		log.Printf("Overwriting Visor genesis run config for nodeset %s", ns.Name)
		if err := co.gen.visorGen.OverwriteRunConfig(ns, co.visorRunTmpl, ""); err != nil {
			return fmt.Errorf("failed to overwrite Visor genesis run config for id %d: %w", ns.Index, err)
		}
	}

	return nil
}
