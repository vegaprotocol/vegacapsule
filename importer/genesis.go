package importer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/generator/genesis"
	"code.vegaprotocol.io/vegacapsule/generator/tendermint"
	"code.vegaprotocol.io/vegacapsule/state"
)

func updateGenesis(netState state.NetworkState) error {
	log.Println("update genesis file for the network")
	gen, err := genesis.NewGenerator(netState.Config, *netState.Config.Network.GenesisTemplate)
	if err != nil {
		return err
	}
	var tendermintGen *tendermint.ConfigGenerator
	tendermintGen, err = tendermint.NewConfigGenerator(netState.Config, netState.GeneratedServices.NodeSets.ToSlice())
	if err != nil {
		return err
	}

	buff, err := gen.Generate(netState.GeneratedServices.GetValidators(), tendermintGen.GenesisValidators())
	if err != nil {
		return fmt.Errorf("failed to generate new genesis; %w", err)
	}
	log.Println("...update genesis file for the network finished")

	for _, nodeSet := range netState.GeneratedServices.NodeSets {
		genesisPath := filepath.Join(nodeSet.Tendermint.HomeDir, "config", "genesis.json")
		log.Printf("write new genesis to the %s file\n", genesisPath)
		if err := os.RemoveAll(genesisPath); err != nil {
			return fmt.Errorf("failed to remove old genesis.json: %w", err)
		}

		if err := ioutil.WriteFile(genesisPath, buff.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to new write genesis.json: %w", err)
		}
	}

	return nil
}
