package state

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type StateType string

type NetworkState struct {
	PreTasks []string

	Config            *config.Config
	GeneratedServices *types.GeneratedServices
}

func (ns NetworkState) Perist() error {
	networkBytes, err := encodeState(ns)
	if err != nil {
		return fmt.Errorf("cannot persist network state: %w", err)
	}

	return ioutil.WriteFile(stateFilePath(ns.Config.OutputDir), networkBytes, 0644)
}

func LoadNetworkState(networkDir string) (*NetworkState, error) {
	statePath := stateFilePath(networkDir)
	configExists, err := utils.FileExists(statePath)
	if err != nil {
		return nil, fmt.Errorf("cannot check network state: %w", err)
	}

	if !configExists {
		return &NetworkState{}, nil
	}

	networkBytes, err := ioutil.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read network state: %w", err)
	}

	return decodeState(networkBytes)
}

func stateFilePath(networkDir string) string {
	return filepath.Join(networkDir, "network.dat")
}
