package cmd

import (
	"context"
	"fmt"
	"log"
	"sort"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/ports"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"
)

var netPrintPortsCmd = &cobra.Command{
	Use:   "addresses",
	Short: "Print all exposed addresses and ports per running network job to stdout",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("network addresses")
		}

		if !networkState.Running() {
			return networkNotRunningErr("network addresses")
		}

		nomadClient, err := nomad.NewClient(nil)
		if err != nil {
			return fmt.Errorf("failed to create nomad client: %w", err)
		}

		nomadRunner, err := nomad.NewJobRunner(nomadClient, "", "")
		if err != nil {
			return fmt.Errorf("failed to create job runner: %w", err)
		}

		return printNetworkAddresses(cmd.Context(), nomadRunner, networkState.GeneratedServices)
	},
}

func printNetworkAddresses(ctx context.Context, nomadRunner *nomad.JobRunner, genServices *types.GeneratedServices) error {
	log.Println("printing exposed network addresses")

	nomadExposedPorts, err := nomadRunner.ListExposedPorts(ctx)
	if err != nil {
		return err
	}

	allOpenPorts, err := ports.OpenPortsPerJob(ctx, nomadExposedPorts, genServices)
	if err != nil {
		return err
	}

	keys := []ports.JobWithTask{}
	for v := range allOpenPorts {
		keys = append(keys, v)
	}

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Name == keys[j].Name {
			return keys[i].TaskName < keys[j].TaskName
		}
		return keys[i].Name < keys[j].Name
	})

	for _, job := range keys {
		ports := allOpenPorts[job]

		if job.TaskName != "" {
			fmt.Printf("Job %q - %q\n", job.Name, job.TaskName)
		} else {
			fmt.Printf("Job %q\n", job.Name)
		}

		for _, port := range ports {
			if port.Name != "" {
				fmt.Printf("  - %s: localhost:%d\n", port.Name, port.Port)
			} else {
				fmt.Printf("  - localhost:%d\n", port.Port)
			}
		}
	}

	return nil
}
