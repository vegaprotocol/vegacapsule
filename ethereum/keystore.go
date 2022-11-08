package ethereum

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeyPair struct {
	Address    string
	PrivateKey string
}

type Signer struct {
	HomeAddress        string
	WalletPassFilePath string
	ClefRPCAddress     string
	KeyPair            KeyPair
}

type SignersList []Signer

func (l SignersList) EthPrivateKeys() []string {
	result := make([]string, len(l))

	for idx, validator := range l {
		result[idx] = validator.KeyPair.PrivateKey
	}

	return result
}

func DescribeKeyPair(keyFilePath, password string) (*KeyPair, error) {
	keys, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore file: %w", err)
	}

	key, err := keystore.DecryptKey(keys, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt ethereum keystore: %w", err)
	}

	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)

	return &KeyPair{
		Address:    key.Address.Hex(),
		PrivateKey: hexutil.Encode(privateKeyBytes)[2:],
	}, nil
}

type KeyStore struct {
	FilePath         string
	PasswordFilePath string
	KeyPair          KeyPair
}

func ImportPrivateKeyIntoKeystore(privateKeyHex, passwordFilePath string, keystorePath string) (*KeyStore, error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	ecdsaKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore file: %w", err)
	}

	passwordBytes, err := os.ReadFile(passwordFilePath)
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
		return &KeyStore{
			FilePath:         filepath.Join(keystorePath, accountFileName),
			PasswordFilePath: passwordFilePath,
			KeyPair: KeyPair{
				PrivateKey: privateKeyHex,
				Address:    ethAccount.Address.Hex(),
			},
		}, nil
	}

	return &KeyStore{
		FilePath:         ethAccount.URL.Path,
		PasswordFilePath: passwordFilePath,
		KeyPair: KeyPair{
			PrivateKey: privateKeyHex,
			Address:    ethAccount.Address.Hex(),
		},
	}, nil
}

func filterEthAccountFileName(keystorePath, publicKey string) (string, error) {
	publicKey = strings.ToLower(strings.TrimPrefix(publicKey, "0x"))

	files, err := os.ReadDir(keystorePath)
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
