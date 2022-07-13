package importer

import (
	"fmt"
	"log"

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

	log.Println("generating and saving genesis")
	if err := gen.GenerateAndSave(netState.GeneratedServices.GetValidators(), netState.GeneratedServices.GetNonValidators(), tendermintGen.GenesisValidators()); err != nil {
		return fmt.Errorf("failed generating and saving genesis: %w", err)
	}
	log.Println("done")

	return nil
}
