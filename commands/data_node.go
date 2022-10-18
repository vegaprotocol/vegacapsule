package commands

import (
	"log"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type lastBlockOutput struct {
	LastBlock int `json:"last_block"`
}

func LastBlock(binary string, homeDir string) (int, error) {
	args := []string{
		config.DataNodeSubCmd,
		"last-block",
		"--output", "json",
		"--home", homeDir,
	}
	lastBlockOut := &lastBlockOutput{}

	out, err := utils.ExecuteBinary(binary, args, lastBlockOut)
	if err != nil {
		log.Println(out)
		return 0, err
	}

	return lastBlockOut.LastBlock, nil
}
