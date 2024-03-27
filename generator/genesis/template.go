package genesis

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"code.vegaprotocol.io/vegacapsule/config"
)

type TemplateContext struct {
	PrimaryBridge   EthereumBridge
	SecondaryBridge EthereumBridge
}

func NewTemplateContext(cfg config.NetworkConfig) (*TemplateContext, error) {
	primaryAddrs := map[string]SmartContract{}
	if err := json.Unmarshal([]byte(*cfg.SmartContractsAddresses), &primaryAddrs); err != nil {
		return nil, fmt.Errorf("could not parse json primary smart contract addresses: %w", err)
	}

	secondaryAddrs := map[string]SmartContract{}
	if err := json.Unmarshal([]byte(*cfg.SecondarySmartContractsAddresses), &secondaryAddrs); err != nil {
		return nil, fmt.Errorf("could not parse json secondary smart contract addresses: %w", err)
	}

	return &TemplateContext{
		PrimaryBridge: EthereumBridge{
			Addresses: primaryAddrs,
			NetworkID: cfg.Ethereum.NetworkID,
			ChainID:   cfg.Ethereum.ChainID,
		},
		SecondaryBridge: EthereumBridge{
			Addresses: secondaryAddrs,
			NetworkID: cfg.SecondaryEthereum.NetworkID,
			ChainID:   cfg.SecondaryEthereum.ChainID,
		},
	}, nil
}

/*
description: |

	Template context also includes functions:
	- `.GetEthContractAddr "contract_name"` - returns contract address based on name.
	- `.GetVegaContractID "contract_name"` - returns contract vega ID based on name.
*/
type EthereumBridge struct {
	// description: Ethereum smart contract addresses created by Vega. These can represent bridges or ERC20 tokens.
	Addresses map[string]SmartContract
	// description: Ethereum network ID.
	NetworkID string
	// description: Ethereum chain ID.
	ChainID string
	// GenValidators []tmtypes.GenesisValidator // TODO add this to the template context
}

type SmartContract struct {
	// description: Ethereum address.
	Ethereum string `json:"Ethereum"`
	// description:  Vega contract ID.
	Vega string `json:"Vega"`
}

func (gc EthereumBridge) GetEthContractAddr(contract string) string {
	sc, ok := gc.Addresses[contract]
	if !ok {
		log.Fatalf("could not find Ethereum smart contract %q", contract)
	}

	if sc.Ethereum == "" {
		log.Fatalf("could not find Ethereum smart contract %q", contract)
	}

	return sc.Ethereum
}

func (gc EthereumBridge) GetVegaContractID(contract string) string {
	sc, ok := gc.Addresses[contract]
	if !ok {
		log.Fatalf("could not find Vega smart contract %q", contract)
	}

	if sc.Vega == "" {
		log.Fatalf("could not find Vega smart contract %q", contract)
	}

	return strings.Replace(sc.Vega, "0x", "", 1)
}
