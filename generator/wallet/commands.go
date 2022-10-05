package wallet

import (
	"log"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type importNetworkOutput struct {
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
}

type initateWalletOutput struct {
	RsaKeys struct {
		PublicKeyFilePath  string `json:"publicKeyFilePath"`
		PrivateKeyFilePath string `json:"privateKeyFilePath"`
	} `json:"rsaKeys"`
}

func (cg *ConfigGenerator) initiateWallet(conf *config.WalletConfig) (*initateWalletOutput, error) {
	args := []string{config.WalletSubCmd, "init", "--output", "json", "--home", cg.homeDir}

	log.Printf("Initiating wallet %q with: %v", conf.Name, args)

	vegaBinary := *cg.conf.VegaBinary
	if conf.VegaBinary != nil {
		vegaBinary = *conf.VegaBinary
	}

	out := &initateWalletOutput{}
	if _, err := utils.ExecuteBinary(vegaBinary, args, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (cg *ConfigGenerator) importNetworkConfig(conf *config.WalletConfig) (*importNetworkOutput, error) {
	args := []string{config.WalletSubCmd, "network", "import", "--output", "json", "--home", cg.homeDir, "--from-file", cg.configFilePath()}

	log.Printf("Importing network to wallet %q with: %v", conf.Name, args)

	vegaBinary := *cg.conf.VegaBinary
	if conf.VegaBinary != nil {
		vegaBinary = *conf.VegaBinary
	}

	out := &importNetworkOutput{}
	if _, err := utils.ExecuteBinary(vegaBinary, args, out); err != nil {
		return nil, err
	}

	return out, nil
}
