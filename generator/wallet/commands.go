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

func (cg *ConfigGenerator) initiateWallet(conf *config.WalletConfig) error {
	args := []string{config.WalletSubCmd, "init", "--output", "json", "--home", cg.homeDir}

	log.Printf("Initiating wallet %q with: %v", conf.Name, args)

	vegaBinary := *cg.conf.VegaBinary
	if conf.VegaBinary != nil {
		vegaBinary = *conf.VegaBinary
	}

	if _, err := utils.ExecuteBinary(vegaBinary, args, nil); err != nil {
		return err
	}

	// If the user has configured a token pass phrase file, we should initialise the token storage
	if conf.TokenPassphraseFile != nil && len(*conf.TokenPassphraseFile) > 0 {
		args = []string{config.WalletSubCmd, "api-token", "init", "--home", cg.homeDir, "--passphrase-file", *conf.TokenPassphraseFile}
		log.Printf("Initiating api-token wallet %q with: %v", conf.Name, args)
		if _, err := utils.ExecuteBinary(vegaBinary, args, nil); err != nil {
			return err
		}
	}

	return nil
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
