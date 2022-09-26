package ethereum

import (
	"fmt"
	"strings"

	"code.vegaprotocol.io/vegacapsule/utils"
)

func getEthereumWalletArgs(signer Signer) []string {
	if signer.ClefRPCAddress != "" {
		return []string{
			"--passphrase-file", signer.WalletPassFilePath,
		}
	}

	return []string{
		"--private-key", signer.KeyPair.PrivateKey,
	}
}

func setThresholdSignature(vegaBinary string, newThreshold int, nonce uint64, submitter string, signers SignersList) (string, error) {
	result := "0x"

	for _, signer := range signers {
		args := []string{
			"set_threshold",
			"--home", signer.HomeAddress,
			"--new-threshold", fmt.Sprintf("%d", newThreshold),
			"--submitter", submitter,
			"--nonce", fmt.Sprintf("%d", nonce),
		}

		args = append(args, getEthereumWalletArgs(signer)...)

		signature, err := callVegaBridgeERC20(vegaBinary, args)
		if err != nil {
			return "", fmt.Errorf("failed to compute set_threshold signature for validator: %s: %w", signer.KeyPair.Address, err)
		}

		result += strings.Trim(
			strings.TrimPrefix(string(signature), "0x"),
			"\n",
		)
	}

	return result, nil
}

func addSignerSignature(vegaBinary string, newSigner string, nonce uint64, submitter string, signers SignersList) (string, error) {
	result := "0x"

	for _, signer := range signers {
		args := []string{
			"add_signer",
			"--home", signer.HomeAddress,
			"--new-signer", newSigner,
			"--submitter", submitter,
			"--nonce", fmt.Sprintf("%d", nonce),
		}

		args = append(args, getEthereumWalletArgs(signer)...)

		signature, err := callVegaBridgeERC20(vegaBinary, args)
		if err != nil {
			return "", fmt.Errorf("failed to compute set_threshold signature for validator: %s: %w", signer.KeyPair.Address, err)
		}

		result += strings.Trim(
			strings.TrimPrefix(string(signature), "0x"),
			"\n",
		)
	}

	return result, nil
}

func removeSignerSignature(vegaBinary string, oldSigner string, nonce uint64, submitter string, signers SignersList) (string, error) {
	result := "0x"

	for _, signer := range signers {
		args := []string{
			"remove_signer",
			"--home", signer.HomeAddress,
			"--old-signer", oldSigner,
			"--submitter", submitter,
			"--nonce", fmt.Sprintf("%d", nonce),
		}

		args = append(args, getEthereumWalletArgs(signer)...)

		signature, err := callVegaBridgeERC20(vegaBinary, args)
		if err != nil {
			return "", fmt.Errorf("failed to compute set_threshold signature for validator: %s: %w", signer.KeyPair.Address, err)
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
