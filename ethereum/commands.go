package ethereum

import (
	"fmt"
	"strings"

	"code.vegaprotocol.io/vegacapsule/utils"
)

func setThresholdSignature(vegaBinary string, vegaHome string, newThreshold int, nonce uint64, submitter string, signers []string) (string, error) {
	result := "0x"

	for _, privKey := range signers {
		signature, err := callVegaBridgeERC20(vegaBinary, []string{
			"set_threshold",
			"--home", vegaHome,
			"--new-threshold", fmt.Sprintf("%d", newThreshold),
			"--submitter", submitter,
			"--nonce", fmt.Sprintf("%d", nonce),
			"--private-key", privKey,
		})

		if err != nil {
			return "", fmt.Errorf("failed to compute set_threshold signature for validator: %s: %w", privKey, err)
		}

		result += strings.Trim(
			strings.TrimPrefix(string(signature), "0x"),
			"\n",
		)
	}

	return result, nil
}

func addSignerSignature(vegaBinary string, vegaHome string, newSigner string, nonce uint64, submitter string, signers []string) (string, error) {
	result := "0x"

	for _, privKey := range signers {
		signature, err := callVegaBridgeERC20(vegaBinary, []string{
			"add_signer",
			"--home", vegaHome,
			"--new-signer", newSigner,
			"--submitter", submitter,
			"--nonce", fmt.Sprintf("%d", nonce),
			"--private-key", privKey,
		})

		if err != nil {
			return "", fmt.Errorf("failed to compute set_threshold signature for validator: %s: %w", privKey, err)
		}

		result += strings.Trim(
			strings.TrimPrefix(string(signature), "0x"),
			"\n",
		)
	}

	return result, nil
}

func removeSignerSignature(vegaBinary string, vegaHome string, oldSigner string, nonce uint64, submitter string, signers []string) (string, error) {
	result := "0x"

	for _, privKey := range signers {
		signature, err := callVegaBridgeERC20(vegaBinary, []string{
			"remove_signer",
			"--home", vegaHome,
			"--old-signer", oldSigner,
			"--submitter", submitter,
			"--nonce", fmt.Sprintf("%d", nonce),
			"--private-key", privKey,
		})

		if err != nil {
			return "", fmt.Errorf("failed to compute set_threshold signature for validator: %s: %w", privKey, err)
		}

		result += strings.Trim(
			strings.TrimPrefix(string(signature), "0x"),
			"\n",
		)
	}

	return result, nil
}

func callVegaBridgeERC20(vegaBinary string, params []string) ([]byte, error) {
	return utils.ExecuteBinary(vegaBinary, append([]string{
		"bridge",
		"erc20",
	}, params...), nil)
}
