package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var netDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy existing network will stop network and removes all it's files",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("destroy")
		}

		if err := netStop(context.Background(), netState); err != nil {
			return fmt.Errorf("failed to stop network: %w", err)
		}

		return netCleanup(homePath)
	},
}

func netCleanup(outputDir string) error {
	log.Println("network cleaning up")

	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("failed cleanup network: %w", err)
	}

	log.Println("network cleaning up success")

	return nil
}
