package faucet

import (
	"log"

	"code.vegaprotocol.io/vegacapsule/utils"
)

type initFaucetOutput struct {
	PublicKey            string `json:"publicKey"`
	FaucetConfigFilePath string `json:"faucetConfigFilePath"`
	FaucetWalletFilePath string `json:"faucetWalletFilePath"`
}

func (cg ConfigGenerator) initiateFaucet(homePath string, phraseFile string) (*initFaucetOutput, error) {
	args := []string{
		"faucet",
		"init",
		"--home", homePath,
		"--passphrase-file", phraseFile,
		"--output", "json",
	}

	log.Printf("Initiating faucet with: %s %v", cg.conf.VegaBinary, args)

	out := &initFaucetOutput{}
	if _, err := utils.ExecuteBinary(*cg.conf.VegaBinary, args, out); err != nil {
		return nil, err
	}

	return out, nil
}
