package wallet

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	vspaths "code.vegaprotocol.io/shared/paths"
	vwconfig "code.vegaprotocol.io/vega/wallet/network"
	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/Masterminds/sprig"
	"github.com/zannen/toml"
)

type ConfigTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	DataNodePrefix       string
	WalletPrefix         string
	Validators           []types.NodeSet
	NonValidators        []types.NodeSet
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	return t, nil
}

type ConfigGenerator struct {
	conf    *config.Config
	homeDir string
}

func NewConfigGenerator(conf *config.Config) (*ConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(*conf.OutputDir, *conf.WalletPrefix))
	if err != nil {
		return nil, err
	}

	return &ConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
	}, nil
}

func (cg *ConfigGenerator) InitiateWithNetworkConfig(conf *config.WalletConfig, validators, nonValidators []types.NodeSet, configTemplate *template.Template) (*types.Wallet, error) {
	if err := os.MkdirAll(cg.homeDir, os.ModePerm); err != nil {
		return nil, err
	}

	initOut, err := cg.initiateWallet(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate wallet %s: %w", conf.Name, err)
	}

	if err := cg.generateNetworkConfig(validators, nonValidators, configTemplate); err != nil {
		return nil, fmt.Errorf("failed to generate network config %q: %w", cg.configFilePath(), err)
	}

	importOut, err := cg.importNetworkConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to import network to wallet %s: %w", conf.Name, err)
	}

	return &types.Wallet{
		Name:                  fmt.Sprintf("%s-wallet", cg.conf.Network.Name),
		HomeDir:               cg.homeDir,
		ServiceConfigFilePath: cg.configFilePath(),
		Network:               importOut.Name,
		PublicKeyFilePath:     initOut.RsaKeys.PublicKeyFilePath,
		PrivateKeyFilePath:    initOut.RsaKeys.PrivateKeyFilePath,
	}, nil
}

func (cg ConfigGenerator) generateNetworkConfig(validators, nonValidators []types.NodeSet, configTemplate *template.Template) error {
	templateCtx := ConfigTemplateContext{
		Prefix:               *cg.conf.Prefix,
		TendermintNodePrefix: *cg.conf.TendermintNodePrefix,
		VegaNodePrefix:       *cg.conf.VegaNodePrefix,
		DataNodePrefix:       *cg.conf.DataNodePrefix,
		WalletPrefix:         *cg.conf.VegaNodePrefix,
		Validators:           validators,
		NonValidators:        nonValidators,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	overrideConfig := vwconfig.Network{}

	if _, err := toml.DecodeReader(buff, &overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	configFilePath := cg.configFilePath()

	if err := vspaths.WriteStructuredFile(configFilePath, overrideConfig); err != nil {
		return fmt.Errorf("failed to write configuration file at %s: %w", configFilePath, err)
	}

	return nil
}

func (cg ConfigGenerator) configFilePath() string {
	return filepath.Join(cg.homeDir, "config.toml")
}
