package config

import (
	"encoding/json"

	"code.vegaprotocol.io/vegacapsule/types"
)

type NodeConfig struct {
	Name  string `hcl:"name,label" cty:"name"`
	Mode  string `hcl:"mode" cty:"mode"`
	Count int    `hcl:"count" cty:"count"`

	NodeWalletPass     string `hcl:"node_wallet_pass,optional" template:"" cty:"node_wallet_pass"`
	EthereumWalletPass string `hcl:"ethereum_wallet_pass,optional" template:"" cty:"ethereum_wallet_pass"`
	VegaWalletPass     string `hcl:"vega_wallet_pass,optional" template:"" cty:"vega_wallet_pass"`

	VegaBinary  *string `hcl:"vega_binary_path,optional"`
	UseDataNode bool    `hcl:"use_data_node,optional" cty:"use_data_node"`
	VisorBinary string  `hcl:"visor_binary,optional"`

	ConfigTemplates ConfigTemplates `hcl:"config_templates,block"`

	PreGenerate *PreGenerate `hcl:"pre_generate,block"`

	PreStartProbe *types.ProbesConfig `hcl:"pre_start_probe,block" template:""`

	ClefWallet *ClefConfig `hcl:"clef_wallet,block" template:""`

	NomadJobTemplate     *string `hcl:"nomad_job_template,optional"`
	NomadJobTemplateFile *string `hcl:"nomad_job_template_file,optional"`
}

type PreGenerate struct {
	Nomad []NomadConfig `hcl:"nomad_job,block"`
}

type ClefConfig struct {
	AccountAddresses []string `hcl:"ethereum_account_addresses" template:""`
	ClefRPCAddr      string   `hcl:"clef_rpc_address" template:""`
}

type ConfigTemplates struct {
	Vega             *string `hcl:"vega,optional"`
	VegaFile         *string `hcl:"vega_file,optional"`
	Tendermint       *string `hcl:"tendermint,optional"`
	TendermintFile   *string `hcl:"tendermint_file,optional"`
	DataNode         *string `hcl:"data_node,optional"`
	DataNodeFile     *string `hcl:"data_node_file,optional"`
	VisorRunConf     *string `hcl:"visor_run_conf,optional"`
	VisorRunConfFile *string `hcl:"visor_run_conf_file,optional"`
	VisorConf        *string `hcl:"visor_conf,optional"`
	VisorConfFile    *string `hcl:"visor_conf_file,optional"`
}

func (nc NodeConfig) Clone() (*NodeConfig, error) {
	origJSON, err := json.Marshal(nc)
	if err != nil {
		return nil, err
	}

	clone := NodeConfig{}
	if err = json.Unmarshal(origJSON, &clone); err != nil {
		return nil, err
	}

	return &clone, nil
}
