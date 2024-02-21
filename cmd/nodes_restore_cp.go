package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	vspaths "code.vegaprotocol.io/vega/paths"
	"code.vegaprotocol.io/vegacapsule/commands"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var checkpointFile string

func incrementChainID(chainID string) (string, error) {
	s := strings.Split(chainID, "-")
	i, err := strconv.ParseInt(s[len(s)-1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("could not increment chain-id %s: %w", chainID, err)
	}
	return strings.Join(s[:len(s)-1], "-") + fmt.Sprintf("-%03d", i+1), nil
}

func updateGenesisChainID(ns types.NodeSet, chainID string) error {
	f, err := os.Open(ns.Tendermint.GenesisFilePath)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		return err
	}

	var genesis map[string]interface{}
	if err := json.Unmarshal(b, &genesis); err != nil {
		return fmt.Errorf("failed unmarshal tendermint config: %w", err)
	}

	genesis["chain_id"] = chainID
	b, err = json.MarshalIndent(genesis, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(ns.Tendermint.GenesisFilePath, b, 0o644)
}

func updateDataNodeChainID(ns types.NodeSet, chainID string) error {
	f, err := os.Open(ns.DataNode.ConfigFilePath)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		return err
	}

	var cfg map[string]interface{}
	if _, err := toml.Decode(string(b), &cfg); err != nil {
		return fmt.Errorf("failed decode datanode config: %w", err)
	}

	cfg["ChainID"] = chainID
	if err := vspaths.WriteStructuredFile(ns.DataNode.ConfigFilePath, cfg); err != nil {
		return fmt.Errorf("failed to write configuration file at %s: %w", ns.DataNode.ConfigFilePath, err)
	}

	return nil
}

var nodesRestoreCheckpointCmd = &cobra.Command{
	Use:   "restore-checkpoint",
	Short: "Restore all Vega nodes state from checkpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("nodes restore-checkpoint")
		}

		if checkpointFile == "" {
			return fmt.Errorf("parameter checkpoint-file can not be empty")
		}

		chainID, err := incrementChainID(netState.VegaChainID)
		if err != nil {
			return err
		}

		netState.VegaChainID = chainID
		if err := netState.Persist(); err != nil {
			return err
		}
		for _, ns := range netState.GeneratedServices.NodeSets {
			r, err := commands.VegaRestoreCheckpoint(
				*netState.Config.VegaBinary,
				ns.Tendermint.HomeDir,
				checkpointFile,
				ns.Vega.NodeWalletPassFilePath,
			)
			if err != nil {
				return fmt.Errorf("failed to restore node %q from checkpoint: %w", ns.Name, err)
			}

			if err := updateGenesisChainID(ns, chainID); err != nil {
				return fmt.Errorf("unable to add new chain id to genesis file for %s: %w", ns.Name, err)
			}

			if ns.DataNode != nil {
				if err := updateDataNodeChainID(ns, chainID); err != nil {
					return fmt.Errorf("unable to add new chain id to datanode config for %s: %w", ns.Name, err)
				}
			}

			fmt.Printf("applied transaction for node set %q: %s", ns.Name, r)
		}
		return nil
	},
}

func init() {
	nodesRestoreCheckpointCmd.PersistentFlags().StringVar(&checkpointFile,
		"checkpoint-file",
		"",
		"Path to the checkpoint file",
	)
	nodesRestoreCheckpointCmd.MarkFlagRequired("checkpoint-file")
}
