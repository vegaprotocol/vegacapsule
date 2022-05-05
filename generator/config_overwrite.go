package generator

import (
	"fmt"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/datanode"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/generator/vega"
	"code.vegaprotocol.io/vegacapsule/types"
)

type configOverride struct {
	tendermintTmpl *template.Template
	vegaTmpl       *template.Template
	dataNodeTmpl   *template.Template
	gen            *Generator
}

func newConfigOverride(gen *Generator, n config.NodeConfig) (*configOverride, error) {
	tendermintTmpl, err := tendermint.NewConfigTemplate(n.ConfigTemplates.Tendermint)
	if err != nil {
		return nil, err
	}

	vegaTmpl, err := vega.NewConfigTemplate(n.ConfigTemplates.Vega)
	if err != nil {
		return nil, err
	}

	var dataNodeTmpl *template.Template
	if n.DataNodeBinary != "" {
		dataNodeTmpl, err = datanode.NewConfigTemplate(n.ConfigTemplates.DataNode)
		if err != nil {
			return nil, err
		}
	}

	return &configOverride{
		tendermintTmpl: tendermintTmpl,
		vegaTmpl:       vegaTmpl,
		dataNodeTmpl:   dataNodeTmpl,
		gen:            gen,
	}, nil
}

func (co *configOverride) Overwrite(nc config.NodeConfig, ns types.NodeSet, fc *types.Faucet) error {
	if co.tendermintTmpl != nil {
		if err := co.gen.tendermintGen.OverwriteConfig(ns, co.tendermintTmpl); err != nil {
			return fmt.Errorf("failed to overwrite Tendermit config for id %d: %w", ns.Index, err)
		}
	}
	if co.vegaTmpl != nil {
		if err := co.gen.vegaGen.OverwriteConfig(ns, fc, co.vegaTmpl); err != nil {
			return fmt.Errorf("failed to overwrite Vega config for id %d: %w", ns.Index, err)
		}
	}
	if co.dataNodeTmpl != nil && ns.DataNode != nil {
		if err := co.gen.dataNodeGen.OverwriteConfig(ns, co.dataNodeTmpl); err != nil {
			return fmt.Errorf("failed to overwrite Data Node config for id %d: %w", ns.Index, err)
		}
	}

	return nil
}
