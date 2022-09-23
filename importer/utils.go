package importer

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

type tendermintKey struct {
	PrivateKey string
	PublicKey  string
	NodeID     string
}

func decodeBase64TendermintPrivateKey(privateKey string) (*tendermintKey, error) {
	var privKey ed25519.PrivKey

	privKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key %s, %w", privateKey, err)
	}
	publicKey := base64.StdEncoding.EncodeToString(privKey.PubKey().Bytes())
	nodeId := privKey.PubKey().Address().String()

	return &tendermintKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		NodeID:     nodeId,
	}, nil
}

func createTempFile(content string) (string, error) {
	file, err := ioutil.TempFile("", "vegacapsule_import")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	if _, err := file.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write content to temp file: %w", err)
	}

	return file.Name(), nil
}
