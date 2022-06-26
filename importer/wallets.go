package importer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
)

type isolatedVegaWallet struct {
	VegaHomePath           string
	RecoveryPhraseFilePath string
	VegaWalletPassFilePath string
	IsolatedWalletName     string
}

func createIsolatedVegaWallet(vegaBinary string, data isolatedVegaWallet, force bool) (string, error) {
	vegaWalletFilePath := filepath.Join(data.VegaHomePath, "data", "wallets", data.IsolatedWalletName)

	if force {
		if err := os.RemoveAll(vegaWalletFilePath); err != nil {
			return "", fmt.Errorf("failed to remove existing vega wallet: %w", err)
		}
	}

	args := []string{
		"wallet", "import",
		"--home", data.VegaHomePath,
		"--no-version-check",
		"--output", "json",
		"--recovery-phrase-file", data.RecoveryPhraseFilePath,
		"--passphrase-file", data.VegaWalletPassFilePath,
		"--wallet", data.IsolatedWalletName,
	}

	out, err := utils.ExecuteBinary(vegaBinary, args, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create isolated vega wallet: %w", err)
	}
	_ = out // TODO: Fix it
	return vegaWalletFilePath, nil
}

type importNodeWalletInput struct {
	VegaHomePath       string
	TendermintHomePath string
	PassphraseFilePath string

	EthKeystoreFilePath     string
	EthKeystorePassFilePath string

	VegaWalletFilePath     string
	VegaWalletPassFilePath string
}

func importVegaNodeWallet(vegaBinary string, data importNodeWalletInput) error {
	tmImportArgs := []string{
		"nodewallet", "import", "--force",
		"--home", data.VegaHomePath,
		"--chain", types.NodeWalletChainTypeTendermint,
		"--passphrase-file", data.PassphraseFilePath,
		"--output", "json",
		"--tendermint-home", data.TendermintHomePath,
	}
	if _, err := utils.ExecuteBinary(vegaBinary, tmImportArgs, nil); err != nil {
		return fmt.Errorf("failed to import tendermint to vega nodewallet: %w", err)
	}

	ethImportArgs := []string{
		"nodewallet", "import", "--force",
		"--home", data.VegaHomePath,
		"--chain", types.NodeWalletChainTypeEthereum,
		"--passphrase-file", data.PassphraseFilePath,
		"--output", "json",
		"--wallet-passphrase-file", data.EthKeystorePassFilePath,
		"--wallet-path", data.EthKeystoreFilePath,
	}
	if _, err := utils.ExecuteBinary(vegaBinary, ethImportArgs, nil); err != nil {
		return fmt.Errorf("failed to import ethereum to vega nodewallet: %w", err)
	}

	vegaImpotyArgs := []string{
		"nodewallet", "import", "--force",
		"--home", data.VegaHomePath,
		"--chain", types.NodeWalletChainTypeVega,
		"--passphrase-file", data.PassphraseFilePath,
		"--output", "json",
		"--wallet-passphrase-file", data.VegaWalletPassFilePath,
		"--wallet-path", data.VegaWalletFilePath,
	}
	if _, err := utils.ExecuteBinary(vegaBinary, vegaImpotyArgs, nil); err != nil {
		return fmt.Errorf("failed to import ethereum to vega nodewallet: %w", err)
	}

	return nil
}

func verifyTendermintNode(vegaBinary, tendermintHomePath, expectedNodeID string) error {
	out, err := utils.ExecuteBinary(vegaBinary, []string{"tm", "show-node-id", "--home", tendermintHomePath}, nil)
	if err != nil {
		return fmt.Errorf("failed to get the tendermint node id: %w", err)
	}

	tendermintNodeID := strings.ToLower(strings.Trim(string(out), " \t\n"))

	if tendermintNodeID != strings.ToLower(expectedNodeID) {
		return fmt.Errorf("tendermint node is invalid: expected \"%s\", got \"%s\"", expectedNodeID, tendermintNodeID)
	}

	return nil
}

func importTendermintNodeKey(tendermintHomePath, privateKey string) error {
	nodeKeyFilePath := filepath.Join(tendermintHomePath, "config", "node_key.json")
	rootNode, err := dasel.NewFromFile(nodeKeyFilePath, "json")
	if err != nil {
		return fmt.Errorf("failed to load private validator file: %w", err)
	}
	if err := rootNode.Put(".priv_key.value", privateKey); err != nil {
		return fmt.Errorf("failed to update address in she priv_validator_key.json file: %w", err)
	}

	if err := rootNode.WriteToFile(nodeKeyFilePath, "json", []storage.ReadWriteOption{}); err != nil {
		return fmt.Errorf("failed to write the node_key.json file: %w", err)
	}

	return nil
}

type tendermintPrivateValidatorData struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

func importTendermintPrivateValidator(tendermintHomePath string, data tendermintPrivateValidatorData) error {
	privateValidatorKeyFilePath := filepath.Join(tendermintHomePath, "config", "priv_validator_key.json")

	rootNode, err := dasel.NewFromFile(privateValidatorKeyFilePath, "json")
	if err != nil {
		return fmt.Errorf("failed to load private validator file: %w", err)
	}
	if err := rootNode.Put(".address", data.Address); err != nil {
		return fmt.Errorf("failed to update address in the priv_validator_key.json file: %w", err)
	}
	if err := rootNode.Put(".pub_key.value", data.PublicKey); err != nil {
		return fmt.Errorf("failed to update public key in the priv_validator_key.json file: %w", err)
	}
	if err := rootNode.Put(".priv_key.value", data.PrivateKey); err != nil {
		return fmt.Errorf("failed to update private key in the priv_validator_key.json file: %w", err)
	}
	if err := rootNode.WriteToFile(privateValidatorKeyFilePath, "json", []storage.ReadWriteOption{}); err != nil {
		return fmt.Errorf("failed to write the priv_validator_key.json file: %w", err)
	}

	return nil
}

type importedEthereumKeystoreInfo struct {
	keystoreFilePath string
	ethereumAddress  string
}

func importPrivateKeyIntoKeystore(privateKeyHex, passwordFilePath, keystorePath string) (*importedEthereumKeystoreInfo, error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	ecdsaKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore file: %w", err)
	}

	passwordBytes, err := ioutil.ReadFile(passwordFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read password file: %w", err)
	}
	password := string(passwordBytes)

	ks := keystore.NewKeyStore(keystorePath, keystore.StandardScryptN, keystore.StandardScryptP)
	if ks == nil {
		return nil, fmt.Errorf("failed to create new keystore")
	}

	ethAccount, err := ks.ImportECDSA(ecdsaKey, password)
	if err != nil {
		if !strings.Contains(err.Error(), "account already exists") {
			return nil, fmt.Errorf("failed to import given private key into created key store: %w", err)
		}

		accountFileName, err := filterEthAccountFileName(keystorePath, ethAccount.Address.Hex())
		if err != nil {
			return nil, fmt.Errorf("the ethereum private key has been imporeted in the past but keystore file cannot be found: %w", err)
		}
		return &importedEthereumKeystoreInfo{
			ethereumAddress:  ethAccount.Address.Hex(),
			keystoreFilePath: filepath.Join(keystorePath, accountFileName),
		}, nil
	}

	return &importedEthereumKeystoreInfo{
		ethereumAddress:  ethAccount.Address.Hex(),
		keystoreFilePath: filepath.Join(keystorePath, ethAccount.URL.Path),
	}, nil
}

func filterEthAccountFileName(keystorePath, publicKey string) (string, error) {
	publicKey = strings.ToLower(strings.TrimPrefix(publicKey, "0x"))

	files, err := ioutil.ReadDir(keystorePath)
	if err != nil {
		return "", fmt.Errorf("failed to read keystore directory: %w", err)
	}

	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name()), fmt.Sprintf("--%s", publicKey)) {
			return file.Name(), nil
		}
	}
	return "", fmt.Errorf("account file for given public key not found")
}
