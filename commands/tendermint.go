package commands

import "code.vegaprotocol.io/vegacapsule/utils"

func TendermintUnsafeResetAll(binary, homeDir string) ([]byte, error) {
	args := []string{
		"tm",
		"unsafe_reset_all",
		"--home", homeDir,
	}

	b, err := utils.ExecuteBinary(binary, args, nil)
	if err != nil {
		return nil, err
	}

	return b, nil
}
