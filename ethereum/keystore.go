package ethereum

import (
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeyPair struct {
	Address    string
	PrivateKey string
}

func DescribeKeyPair(keyFilePath, password string) (*KeyPair, error) {
	keys, err := ioutil.ReadFile(keyFilePath)
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
