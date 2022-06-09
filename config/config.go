package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Config struct {
	OutputDir            *string       `hcl:"-"`
	VegaBinary           *string       `hcl:"vega_binary_path"`
	Prefix               *string       `hcl:"prefix"`
	NodeDirPrefix        *string       `hcl:"node_dir_prefix"`
	TendermintNodePrefix *string       `hcl:"tendermint_node_prefix"`
	VegaNodePrefix       *string       `hcl:"vega_node_prefix"`
	DataNodePrefix       *string       `hcl:"data_node_prefix"`
	WalletPrefix         *string       `hcl:"wallet_prefix"`
	FaucetPrefix         *string       `hcl:"faucet_prefix"`
	Network              NetworkConfig `hcl:"network,block"`

	// Internal helper variables
	configDir string
}

type NetworkConfig struct {
	Name                string         `hcl:"name,label"`
	GenesisTemplate     *string        `hcl:"genesis_template"`
	GenesisTemplateFile *string        `hcl:"genesis_template_file"`
	Ethereum            EthereumConfig `hcl:"ethereum,block"`
	Wallet              *WalletConfig  `hcl:"wallet,block"`
	Faucet              *FaucetConfig  `hcl:"faucet,block"`

	PreStart  *PStartConfig `hcl:"pre_start,block"`
	PostStart *PStartConfig `hcl:"post_start,block"`

	Nodes []NodeConfig `hcl:"node_set,block"`

	SmartContractsAddresses     *string `hcl:"smart_contracts_addresses,optional"`
	SmartContractsAddressesFile *string `hcl:"smart_contracts_addresses_file,optional"`
}

type CommandRunner struct {
	PathsMapping ConfigRemoteNetworkPathsMapping `hcl:"paths_mapping,block"`

	Meta                 map[string]string `hcl:"meta,optional"`
	NomadJobTemplate     *string           `hcl:"nomad_job_template,optional"`
	NomadJobTemplateFile *string           `hcl:"nomad_job_template_file,optional"`
}

func (nc NetworkConfig) GetNodeConfig(name string) (*NodeConfig, error) {
	for _, nodeConf := range nc.Nodes {
		if nodeConf.Name == name {
			return &nodeConf, nil
		}
	}

	return nil, fmt.Errorf("node config with name %q not found", name)
}

type EthereumConfig struct {
	ChainID   string `hcl:"chain_id"`
	NetworkID string `hcl:"network_id"`
	Endpoint  string `hcl:"endpoint"`
}

type PStartConfig struct {
	Docker []DockerConfig `hcl:"docker_service,block"`
}

type StaticPort struct {
	To    int `hcl:"to,optional"`
	Value int `hcl:"value"`
}

type DockerConfig struct {
	Name         string            `hcl:"name,label"`
	Image        string            `hcl:"image"`
	Command      string            `hcl:"cmd"`
	Args         []string          `hcl:"args"`
	Env          map[string]string `hcl:"env,optional"`
	StaticPort   *StaticPort       `hcl:"static_port,block"`
	AuthSoftFail bool              `hcl:"auth_soft_fail,optional"`
}

type WalletConfig struct {
	Name     string `hcl:"name,label"`
	Binary   string `hcl:"binary"`
	Template string `hcl:"template,optional"`
}

type FaucetConfig struct {
	Name     string `hcl:"name,label"`
	Pass     string `hcl:"wallet_pass"`
	Template string `hcl:"template,optional"`
}

type NomadConfig struct {
	Name            string  `hcl:"name,label"`
	JobTemplate     *string `hcl:"job_template,optional"`
	JobTemplateFile *string `hcl:"job_template_file,optional"`
}

type PreGenerate struct {
	Nomad []NomadConfig `hcl:"nomad_job,block"`
}

type NodeConfig struct {
	Name  string `hcl:"name,label"`
	Mode  string `hcl:"mode"`
	Count int    `hcl:"count"`

	PreGenerate *PreGenerate `hcl:"pre_generate,block"`

	ClefWallet *ClefConfig `hcl:"clef_wallet,block" template:""`

	NodeWalletPass     string `hcl:"node_wallet_pass,optional" template:""`
	EthereumWalletPass string `hcl:"ethereum_wallet_pass,optional" template:""`
	VegaWalletPass     string `hcl:"vega_wallet_pass,optional" template:""`

	DataNodeBinary string `hcl:"data_node_binary,optional"`

	ConfigTemplates      ConfigTemplates `hcl:"config_templates,block"`
	NomadJobTemplate     *string         `hcl:"nomad_job_template,optional"`
	NomadJobTemplateFile *string         `hcl:"nomad_job_template_file,optional"`

	// We may want to call some commands on the remote nodes. For some scenarios, We could
	// use the network allocations to execute commands on the server — however, there's a
	// set of commands that require a stopped network before running them. An example of
	// the command that requires a stopped network is unsafe-reset-all. We may need a
	// runner/helper to run commands for this kind of command. If the below section is
	// not provided only local commands will be available :()
	RemoteCommandRunner *CommandRunner `hcl:"remote_command_runner,block"`
}

type ConfigRemoteNetworkPathsMapping struct {
	TendermintHome *string `hcl:"tendermint_home"`
	VegaHome       *string `hcl:"vega_home"`
	DataNodeHome   *string `hcl:"data_node_home,optional"`

	VegaBinary     *string `hcl:"vega_binary"`
	DataNodeBinary *string `hcl:"data_node_binary,optional"`
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

type ClefConfig struct {
	AccountAddress string `hcl:"ethereum_account_address" template:""`
	ClefRPCAddr    string `hcl:"clef_rpc_address" template:""`
}

type ConfigTemplates struct {
	Vega           *string `hcl:"vega,optional"`
	VegaFile       *string `hcl:"vega_file,optional"`
	Tendermint     *string `hcl:"tendermint,optional"`
	TendermintFile *string `hcl:"tendermint_file,optional"`
	DataNode       *string `hcl:"data_node,optional"`
	DataNodeFile   *string `hcl:"data_node_file,optional"`
}

func (c *Config) setAbsolutePaths() error {
	// Output directory
	if !filepath.IsAbs(*c.OutputDir) {
		absPath, err := filepath.Abs(*c.OutputDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for outputDir: %w", err)
		}
		*c.OutputDir = absPath
	}

	// Vega binary
	vegaBinPath, err := utils.BinaryAbsPath(*c.VegaBinary)
	if err != nil {
		return err
	}
	*c.VegaBinary = vegaBinPath

	// Wallet binary
	if c.Network.Wallet != nil {
		walletBinPath, err := utils.BinaryAbsPath(c.Network.Wallet.Binary)
		if err != nil {
			return err
		}
		c.Network.Wallet.Binary = walletBinPath
	}

	// Data nodes binaries
	for idx, nc := range c.Network.Nodes {
		if nc.DataNodeBinary == "" {
			continue
		}

		dataNodeBinPath, err := utils.BinaryAbsPath(nc.DataNodeBinary)
		if err != nil {
			return err
		}
		c.Network.Nodes[idx].DataNodeBinary = dataNodeBinPath
	}

	return nil
}

func (c *Config) Validate(configDir string) error {
	if err := c.setAbsolutePaths(); err != nil {
		return fmt.Errorf("failed to set absolute paths: %w", err)
	}

	c.configDir = configDir

	if err := c.validateCommandRunnerConfig(); err != nil {
		return fmt.Errorf("failed to validate command runner config: %w", err)
	}

	if err := c.loadRemoteCommandRunnerConfig(); err != nil {
		return fmt.Errorf("failed to load command runner config files: %w", err)
	}

	if err := c.loadAndValidateGenesis(); err != nil {
		return fmt.Errorf("failed to validate genesis: %w", err)
	}

	if err := c.loadAndValidateNodeSets(); err != nil {
		return fmt.Errorf("failed to validate node configs: %w", err)
	}

	if err := c.loadAndValidatSetSmartContractsAddresses(); err != nil {
		return fmt.Errorf("invalid configuration for smart contrtacts addresses: %w", err)
	}

	return nil
}

func (c *Config) loadAndValidateNodeSets() error {
	mErr := utils.NewMultiError()

	for i, nc := range c.Network.Nodes {
		updatedNc, err := c.loadAndValidateNomadJobTemplates(nc)
		if err != nil {
			mErr.Add(fmt.Errorf("failed to validate nomad job template for %s: %w", nc.Name, err))
			continue
		}

		updatedCt, err := c.loadAndValidateConfigTemplates(nc.ConfigTemplates)
		if err != nil {
			mErr.Add(fmt.Errorf("failed to validate node set config templates: %w", err))
			continue
		}

		updatedNc.ConfigTemplates = *updatedCt

		if nc.PreGenerate != nil {
			updatedPreGen, err := c.loadAndValidatePreGenerate(*nc.PreGenerate)
			if err != nil {
				mErr.Add(fmt.Errorf("failed to validate node set pre generate templates: %w", err))
				return err
			}

			updatedNc.PreGenerate = updatedPreGen
		}

		c.Network.Nodes[i] = *updatedNc
	}

	if mErr.HasAny() {
		return mErr
	}

	return nil
}

func (c Config) loadAndValidatePreGenerate(preGen PreGenerate) (*PreGenerate, error) {
	mErr := utils.NewMultiError()

	for i, nc := range preGen.Nomad {
		if nc.JobTemplate == nil && nc.JobTemplateFile != nil {
			tmpl, err := c.loadConfigTemplateFile(*nc.JobTemplateFile)
			if err != nil {
				mErr.Add(fmt.Errorf("failed to load pre generate nomad template file for %s: %w", nc.Name, err))

				continue
			}

			nc.JobTemplate = &tmpl
			nc.JobTemplateFile = nil

			preGen.Nomad[i] = nc
		}
	}

	if mErr.HasAny() {
		return nil, mErr
	}

	return &preGen, nil
}

func (c Config) loadAndValidateConfigTemplates(ct ConfigTemplates) (*ConfigTemplates, error) {
	mErr := utils.NewMultiError()

	if ct.Vega == nil && ct.VegaFile != nil {
		tmpl, err := c.loadConfigTemplateFile(*ct.VegaFile)
		if err != nil {
			mErr.Add(fmt.Errorf("failed to load Vega config template: %w", err))
		} else {
			ct.Vega = &tmpl
			ct.VegaFile = nil
		}
	}

	if ct.Tendermint == nil && ct.TendermintFile != nil {
		tmpl, err := c.loadConfigTemplateFile(*ct.TendermintFile)
		if err != nil {
			mErr.Add(fmt.Errorf("failed to load Tendermint config template: %w", err))
		} else {
			ct.Tendermint = &tmpl
			ct.TendermintFile = nil
		}
	}

	if ct.DataNode == nil && ct.DataNodeFile != nil {
		tmpl, err := c.loadConfigTemplateFile(*ct.DataNodeFile)
		if err != nil {
			mErr.Add(fmt.Errorf("failed to load Data Node config template: %w", err))
		} else {
			ct.DataNode = &tmpl
			ct.DataNodeFile = nil
		}
	}

	if mErr.HasAny() {
		return nil, mErr
	}

	return &ct, nil
}

func (c Config) loadConfigTemplateFile(path string) (string, error) {
	templateFile, err := utils.AbsPathWithPrefix(c.configDir, path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute file path %q: %w", path, err)
	}

	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", templateFile, err)
	}

	return string(template), nil
}

func (c Config) loadAndValidateNomadJobTemplates(nc NodeConfig) (*NodeConfig, error) {
	if nc.NomadJobTemplate != nil {
		return &nc, nil
	}

	if nc.NomadJobTemplateFile == nil {
		return &nc, nil
	}

	templateFile, err := utils.AbsPathWithPrefix(c.configDir, *nc.NomadJobTemplateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute file path %q: %w", *nc.NomadJobTemplateFile, err)
	}

	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", templateFile, err)
	}

	str := string(template)
	nc.NomadJobTemplate = &str
	nc.NomadJobTemplateFile = nil

	return &nc, nil
}

func (c *Config) loadAndValidateGenesis() error {
	if c.Network.GenesisTemplate != nil {
		return nil
	}

	if c.Network.GenesisTemplateFile != nil {
		genTemplateFile, err := utils.AbsPathWithPrefix(c.configDir, *c.Network.GenesisTemplateFile)
		if err != nil {
			return fmt.Errorf("failed to get absolute file path %q: %w", genTemplateFile, err)
		}

		genTemplate, err := ioutil.ReadFile(genTemplateFile)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", genTemplateFile, err)
		}

		genTemplateStr := string(genTemplate)
		// set file content as template a set file path to nil
		c.Network.GenesisTemplate = &genTemplateStr
		c.Network.GenesisTemplateFile = nil

		return nil
	}

	return fmt.Errorf("missing genesis file template: either genesis_template or genesis_template_file must be defined")
}

func (c *Config) loadAndValidatSetSmartContractsAddresses() error {
	if c.Network.SmartContractsAddresses == nil {
		if c.Network.SmartContractsAddressesFile == nil {
			return fmt.Errorf("missing smart contracts file: either smart_contracts_addresses or smart_contracts_addresses_file must be defined")
		}

		smartContractsFile, err := utils.AbsPathWithPrefix(c.configDir, *c.Network.SmartContractsAddressesFile)
		if err != nil {
			return fmt.Errorf("failed to get absolute file path %q: %w", smartContractsFile, err)
		}

		smartContracts, err := ioutil.ReadFile(smartContractsFile)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", smartContractsFile, err)
		}

		smartContractsStr := string(smartContracts)

		c.Network.SmartContractsAddresses = &smartContractsStr
		c.Network.SmartContractsAddressesFile = nil
	}

	_, err := c.SmartContractsInfo()

	if err != nil {
		return fmt.Errorf("failed to check smart contract addreses: %w", err)
	}

	return nil
}

func (c Config) validateCommandRunnerConfig() error {
	remoteConfigProvided := false

	mErr := utils.NewMultiError()

	for _, nodeConf := range c.Network.Nodes {
		if remoteConfigProvided && nodeConf.RemoteCommandRunner == nil {
			mErr.Add(fmt.Errorf("missing the `remote_command_runner` configuration for the \"%s\" node set:  the `remote_command_runner` configuration must be specified either for none or all nodes sets' configuration", nodeConf.Name))
		}

		if nodeConf.RemoteCommandRunner != nil {
			remoteConfigProvided = true
		}
	}

	// No remode command runner provided, no required further falidation. Network will not support remote commands
	if !remoteConfigProvided {
		return nil
	}

	for _, nodeConf := range c.Network.Nodes {
		if err := nodeConf.RemoteCommandRunner.Validate(nodeConf.DataNodeBinary != ""); err != nil {
			mErr.Add(fmt.Errorf("failed to validate the `remote_command_runner` for the \"%s\" node set: %w", nodeConf.Name, err))
		}
	}

	if mErr.HasAny() {
		return mErr
	}

	return nil
}

func (conf *Config) loadRemoteCommandRunnerConfig() error {
	mErr := utils.NewMultiError()

	for idx, nodeConf := range conf.Network.Nodes {
		if utils.EmptyStrPoint(nodeConf.RemoteCommandRunner.NomadJobTemplate) && !utils.EmptyStrPoint(nodeConf.RemoteCommandRunner.NomadJobTemplateFile) {
			tmpl, err := conf.loadConfigTemplateFile(*nodeConf.RemoteCommandRunner.NomadJobTemplateFile)

			if err != nil {
				mErr.Add(fmt.Errorf("failed to load remote command runner nomad template: %w", err))
				continue
			}
			conf.Network.Nodes[idx].RemoteCommandRunner.NomadJobTemplate = &tmpl
			conf.Network.Nodes[idx].RemoteCommandRunner.NomadJobTemplateFile = nil
		}
	}

	if mErr.HasAny() {
		return mErr
	}

	return nil
}

func (c Config) SmartContractsInfo() (*types.SmartContractsInfo, error) {
	smartcontracts := &types.SmartContractsInfo{}

	if err := json.Unmarshal([]byte(*c.Network.SmartContractsAddresses), &smartcontracts); err != nil {
		return nil, fmt.Errorf("failed to get smart contracts info: config.network.smart_contracts_addresses format is wrong: %w", err)
	}

	return smartcontracts, nil
}

func (c *Config) Persist() error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(*c, f.Body())
	return ioutil.WriteFile(filepath.Join(*c.OutputDir, "config.hcl"), f.Bytes(), 0644)
}

func DefaultConfig() (*Config, error) {
	outputDir, err := DefaultNetworkHome()
	if err != nil {
		return nil, err
	}

	return &Config{
		OutputDir:            &outputDir,
		Prefix:               utils.StrPoint("st-local"),
		NodeDirPrefix:        utils.StrPoint("node"),
		TendermintNodePrefix: utils.StrPoint("tendermint"),
		VegaNodePrefix:       utils.StrPoint("vega"),
		DataNodePrefix:       utils.StrPoint("data"),
		WalletPrefix:         utils.StrPoint("wallet"),
		FaucetPrefix:         utils.StrPoint("faucet"),
		VegaBinary:           utils.StrPoint("vega"),
	}, nil
}

func DefaultNetworkHome() (string, error) {
	capsuleHome, err := utils.CapsuleHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(capsuleHome, "testnet"), nil
}

func (runner CommandRunner) Validate(withDataNode bool) error {
	if runner.NomadJobTemplate == nil && runner.NomadJobTemplateFile == nil {
		return fmt.Errorf("either `nomad_job_template` or `nomad_job_template_file` must be specified for the remote command runner")
	}

	errorMsg := "the `remote_network_path_mapping` configuration is not fully specified: the \"%s\" field is missing"

	if utils.EmptyStrPoint(runner.PathsMapping.VegaHome) {
		return fmt.Errorf(errorMsg, "vega_home")
	}

	if utils.EmptyStrPoint(runner.PathsMapping.TendermintHome) {
		return fmt.Errorf(errorMsg, "tendermint_home")
	}

	if utils.EmptyStrPoint(runner.PathsMapping.VegaBinary) {
		return fmt.Errorf(errorMsg, "vega_binary")
	}

	// No data node is started, no need to provide data node paths mapping
	if !withDataNode {
		return nil
	}

	if utils.EmptyStrPoint(runner.PathsMapping.DataNodeHome) {
		return fmt.Errorf(fmt.Sprintf("%s: field is required for non validator node", errorMsg), "data_node_home")
	}

	if utils.EmptyStrPoint(runner.PathsMapping.DataNodeBinary) {
		return fmt.Errorf(fmt.Sprintf("%s: field is required for non validator node", errorMsg), "data_node_binary")
	}

	return nil
}
