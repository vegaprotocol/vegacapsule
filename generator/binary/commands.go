package binary

import (
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type initBinaryOutput struct {
	BinaryConfigFilePath string `json:"binaryConfigFilePath"`
}

func (cg *ConfigGenerator) initiateBinary(conf *config.BinaryConfig) (*initBinaryOutput, error) {
	if conf == nil {
		return nil, fmt.Errorf("binary config is nil")
	}
	log.Printf("Initiating binary with: %s %v", *conf.BinaryFile, conf.Args)

	out := &initBinaryOutput{}
	if _, err := utils.ExecuteBinary(*conf.BinaryFile, conf.Args, out); err != nil {
		return nil, err
	}

	return out, nil
}
