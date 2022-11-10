package genesis

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type SmartContract struct {
	Ethereum string `json:"Ethereum"`
	Vega     string `json:"Vega"`
}

type TemplateContext struct {
	Addresses map[string]SmartContract
	ChainID   string
	NetworkID string
	// GenValidators []tmtypes.GenesisValidator // TODO add this to the template context
}

func NewTemplateContext(chainID, networkID string, addressesJSON []byte) (*TemplateContext, error) {
	addrs := map[string]SmartContract{}

	if err := json.Unmarshal(addressesJSON, &addrs); err != nil {
		return nil, fmt.Errorf("could not parse json smart contract addresses: %s", addressesJSON)
	}

	return &TemplateContext{
		ChainID:   chainID,
		NetworkID: networkID,
		Addresses: addrs,
	}, nil
}

func (gc TemplateContext) GetEthContractAddr(contract string) string {
	sc, ok := gc.Addresses[contract]
	if !ok {
		log.Fatalf("could not find Ethereum smart contract %q", contract)
	}

	if sc.Ethereum == "" {
		log.Fatalf("could not find Ethereum smart contract %q", contract)
	}

	return sc.Ethereum
}

func (gc TemplateContext) GetVegaContractID(contract string) string {
	sc, ok := gc.Addresses[contract]
	if !ok {
		log.Fatalf("could not find Vega smart contract %q", contract)
	}

	if sc.Vega == "" {
		log.Fatalf("could not find Vega smart contract %q", contract)
	}

	return strings.Replace(sc.Vega, "0x", "", 1)
}
