package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var netDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Stop the network and removes all of its files",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if err := netStop(context.Background(), netState); err != nil {
			if nomad.IsConnectionErr(err) {
				log.Println("Couldn't connect to nomad, skipping the network shutdown...")
			} else {
				return fmt.Errorf("failed to stop the network: %w", err)
			}
		}

		return netCleanup(homePath)
	},
}

func netCleanup(outputDir string) error {
	log.Println("Cleaning up the network...")

	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("failed to cleanup the network: %w", err)
	}

	log.Println("Network has been successfully cleaned up.")

	return nil
}
