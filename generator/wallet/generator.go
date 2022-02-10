package wallet

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"

	vspaths "code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
	vwconfig "code.vegaprotocol.io/vegawallet/network"
	"github.com/imdario/mergo"
	"github.com/zannen/toml"
)

type ConfigTemplateContext struct {
	Prefix               string
	TendermintNodePrefix string
	VegaNodePrefix       string
	DataNodePrefix       string
	WalletPrefix         string
	Validators           []types.NodeSet
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Parse(templateRaw)
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
	homeDir, err := filepath.Abs(path.Join(conf.OutputDir, conf.WalletPrefix))
	if err != nil {
		return nil, err
	}

	return &ConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
	}, nil
}

type initateWalletOutput struct {
	ServiceConfigFilePath string `json:"serviceConfigFilePath"`
	RsaKeys               struct {
		PublicKeyFilePath  string `json:"publicKeyFilePath"`
		PrivateKeyFilePath string `json:"privateKeyFilePath"`
	} `json:"rsaKeys"`
}

func (cg *ConfigGenerator) initiateWallet(conf *config.WalletConfig) (*initateWalletOutput, error) {
	args := []string{"init", "--no-version-check", "--output", "json", "--home", cg.homeDir}

	log.Printf("Initiating wallet %q with: %v", conf.Name, args)

	out := &initateWalletOutput{}
	if _, err := utils.ExecuteBinary(conf.Binary, args, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (cg *ConfigGenerator) Initiate(conf *config.WalletConfig) (*types.Wallet, error) {
	if err := os.MkdirAll(cg.homeDir, os.ModePerm); err != nil {
		return nil, err
	}

	out, err := cg.initiateWallet(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate wallet %s: %w", conf.Name, err)
	}

	return &types.Wallet{
		HomeDir:               cg.homeDir,
		ServiceConfigFilePath: out.ServiceConfigFilePath,
		PublicKeyFilePath:     out.RsaKeys.PublicKeyFilePath,
		PrivateKeyFilePath:    out.RsaKeys.PrivateKeyFilePath,
	}, nil
}

func (cg ConfigGenerator) OverwriteConfig(wallet types.Wallet, validators []types.NodeSet, configTemplate *template.Template) error {
	templateCtx := ConfigTemplateContext{
		Prefix:               cg.conf.Prefix,
		TendermintNodePrefix: cg.conf.TendermintNodePrefix,
		VegaNodePrefix:       cg.conf.VegaNodePrefix,
		DataNodePrefix:       cg.conf.DataNodePrefix,
		WalletPrefix:         cg.conf.VegaNodePrefix,
		Validators:           validators,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	overrideConfig := vwconfig.Network{}

	if _, err := toml.DecodeReader(buff, &overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	configFilePath := wallet.ServiceConfigFilePath

	vegaConfig := vwconfig.Network{}
	if err := vspaths.ReadStructuredFile(configFilePath, &vegaConfig); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", configFilePath, err)
	}

	if err := mergo.Merge(&overrideConfig, vegaConfig); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := vspaths.WriteStructuredFile(configFilePath, overrideConfig); err != nil {
		return fmt.Errorf("failed to write configuration file at %s: %w", configFilePath, err)
	}

	return nil
}
