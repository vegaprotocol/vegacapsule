package commands

import "code.vegaprotocol.io/vegacapsule/utils"

func VegaUnsafeResetAll(binary, homeDir string) ([]byte, error) {
	args := []string{
		"unsafe_reset_all",
		"--home", homeDir,
	}

	b, err := utils.ExecuteBinary(binary, args, nil)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func VegaRestoreCheckpoint(binary, homeDir, checkpointFile, nodeWalletPhraseFile string) ([]byte, error) {
	args := []string{
		"genesis",
		"load_checkpoint",
		"--tm-home", homeDir,
		"--checkpoint-path", checkpointFile,
		"--passphrase-file", nodeWalletPhraseFile,
	}

	b, err := utils.ExecuteBinary(binary, args, nil)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func VegaProtocolUpgradeProposal(binary, homeDir, releaseTag, height, nodeWalletPhraseFile string) ([]byte, error) {
	args := []string{
		"protocol_upgrade_proposal",
		"--home", homeDir,
		"--vega-release-tag", releaseTag,
		"--height", height,
		"--passphrase-file", nodeWalletPhraseFile,
	}

	b, err := utils.ExecuteBinary(binary, args, nil)
	if err != nil {
		return nil, err
	}

	return b, nil
}
