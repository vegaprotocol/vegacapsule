package faucet

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vega/core/faucet"
	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
)

type ConfigTemplateContext struct {
	Prefix    string
	HomeDir   string
	PublicKey string
}

func NewConfigTemplate(templateRaw string) (*template.Template, error) {
	t, err := template.New("config.toml").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the faucet's template config: %w", err)
	}

	return t, nil
}

type ConfigGenerator struct {
	conf    *config.Config
	homeDir string
}

func NewConfigGenerator(conf *config.Config) (*ConfigGenerator, error) {
	homeDir, err := filepath.Abs(path.Join(*conf.OutputDir, *conf.FaucetPrefix))
	if err != nil {
		return nil, err
	}

	return &ConfigGenerator{
		conf:    conf,
		homeDir: homeDir,
	}, nil
}

func (cg *ConfigGenerator) InitiateAndConfigure(conf *config.FaucetConfig) (*types.Faucet, error) {
	initFaucet, err := cg.Initiate(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to init faucet for %q: %w", conf.Name, err)
	}

	if conf.Template == "" {
		return initFaucet, nil
	}

	configTemplate, err := NewConfigTemplate(conf.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to create new template config for %q: %w", conf.Name, err)
	}

	if err := cg.OverwriteConfig(initFaucet, configTemplate); err != nil {
		return nil, fmt.Errorf("failed to override config for %q: %w", conf.Name, err)
	}

	return initFaucet, nil
}

func (cg *ConfigGenerator) Initiate(conf *config.FaucetConfig) (*types.Faucet, error) {
	if err := os.MkdirAll(cg.homeDir, os.ModePerm); err != nil {
		return nil, err
	}

	walletPassFilePath := path.Join(cg.homeDir, "faucet-wallet-pass.txt")
	if err := ioutil.WriteFile(walletPassFilePath, []byte(conf.Pass), 0644); err != nil {
		return nil, fmt.Errorf("failed to write faucet wallet passphrase to file: %w", err)
	}

	initOut, err := cg.initiateFaucet(cg.homeDir, walletPassFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initiated faucet %q: %w", conf.Name, err)
	}

	return &types.Faucet{
		GeneratedService: types.GeneratedService{
			Name:           fmt.Sprintf("%s-faucet", cg.conf.Network.Name),
			ConfigFilePath: initOut.FaucetConfigFilePath,
			HomeDir:        cg.homeDir,
		},
		PublicKey:          initOut.PublicKey,
		WalletFilePath:     initOut.FaucetWalletFilePath,
		WalletPassFilePath: walletPassFilePath,
		BinaryPath:         *cg.conf.VegaBinary,
	}, nil
}

func (cg ConfigGenerator) OverwriteConfig(fc *types.Faucet, configTemplate *template.Template) error {
	templateCtx := ConfigTemplateContext{
		Prefix:    *cg.conf.Prefix,
		HomeDir:   cg.homeDir,
		PublicKey: fc.PublicKey,
	}

	buff := bytes.NewBuffer([]byte{})

	if err := configTemplate.Execute(buff, templateCtx); err != nil {
		return fmt.Errorf("failed to execute template for faucet: %w", err)
	}

	overrideConfig := faucet.Config{}

	if _, err := toml.NewDecoder(buff).Decode(&overrideConfig); err != nil {
		return fmt.Errorf("failed decode override config: %w", err)
	}

	configFilePath := fc.ConfigFilePath

	defaultConfig := faucet.NewDefaultConfig()
	if err := paths.ReadStructuredFile(configFilePath, &defaultConfig); err != nil {
		return fmt.Errorf("failed to read configuration file at %s: %w", configFilePath, err)
	}

	if err := mergo.Merge(&overrideConfig, defaultConfig); err != nil {
		return fmt.Errorf("failed to merge configs: %w", err)
	}

	if err := paths.WriteStructuredFile(configFilePath, overrideConfig); err != nil {
		return fmt.Errorf("failed to write configuration file for faucet at %s: %w", configFilePath, err)
	}

	return nil
}
