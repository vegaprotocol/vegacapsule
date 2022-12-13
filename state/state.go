package state

import (
	"fmt"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

type StateType string

type NetworkState struct {
	Config            *config.Config
	GeneratedServices *types.GeneratedServices
	RunningJobs       *types.NetworkJobs
	VegaChainID       string
}

func (ns *NetworkState) Empty() bool {
	return ns == nil || ns.Config == nil || len(ns.GeneratedServices.NodeSets) == 0
}

func (ns *NetworkState) Running() bool {
	return !ns.Empty() && ns.RunningJobs != nil && len(ns.RunningJobs.NodesSetsJobIDs) != 0
}

func (ns NetworkState) Persist() error {
	networkBytes, err := encodeState(ns)
	if err != nil {
		return fmt.Errorf("failed to persist network state: %w", err)
	}

	if err := os.WriteFile(stateFilePath(*ns.Config.OutputDir), networkBytes, 0o644); err != nil {
		return fmt.Errorf("failed to persist network state: %w", err)
	}

	return nil
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

	networkBytes, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read network state: %w", err)
	}

	netState, err := decodeState(networkBytes)
	if err != nil {
		return nil, err
	}

	netState.Config.OutputDir = &networkDir

	return netState, nil
}

func stateFilePath(networkDir string) string {
	return filepath.Join(networkDir, "network.dat")
}
