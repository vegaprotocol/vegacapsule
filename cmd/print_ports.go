package cmd

import (
	"context"
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/ports"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/spf13/cobra"
)

var printPortsCmd = &cobra.Command{
	Use: "print-ports",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("nodes start")
		}

		nomadClient, err := nomad.NewClient(nil)
		if err != nil {
			return fmt.Errorf("failed to create nomad client: %w", err)
		}

		nomadRunner, err := nomad.NewJobRunner(nomadClient, "", "")
		if err != nil {
			return fmt.Errorf("failed to create job runner: %w", err)
		}

		printNetworkPorts(cmd.Context(), nomadRunner, networkState.GeneratedServices)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(printPortsCmd)
}

func printNetworkPorts(ctx context.Context, nomadRunner *nomad.JobRunner, genServices *types.GeneratedServices) error {
	log.Println("collecting exposed nodes addresses")

	jobs, err := nomadRunner.Client.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list running network jobs: %w", err)
	}

	portsPerProcessName, err := ports.ScanNetworkProcessesPorts()
	if err != nil {
		return fmt.Errorf("failed to scan on open ports on os: %w", err)
	}

	for _, j := range jobs {
		gss := genServices.GetByName(j.ID)
		if len(gss) == 0 {
			nomadExposedPorts, err := nomadRunner.ListExposedPorts(ctx, j.ID)
			if err != nil {
				log.Printf("failed to ") // TODO add the log
				continue
			}

			for _, openPort := range nomadExposedPorts {
				printJobPort(openPort, "", j.ID, "")
			}
			continue
		}

		for _, gs := range gss {
			configuredPorts, err := ports.ExtractPortsFromConfig(gs.ConfigFilePath)
			if err != nil {
				log.Printf("failed to extract ports from config file %q for job %s: %s", gs.ConfigFilePath, j.ID, err)
				continue
			}

			for openPort, processName := range portsPerProcessName {
				if portName, ok := configuredPorts[openPort]; ok {
					printJobPort(openPort, portName, j.ID, processName)
				}
			}
		}
	}

	return nil
}

func printJobPort(
	port int64,
	portName string,
	jobName string,
	taskName string,
) {
	if taskName != "" {
		fmt.Printf("\n%q %s:\n", jobName, taskName)
	} else {
		fmt.Printf("\n%q:\n", jobName)
	}

	if portName != "" {
		fmt.Printf("- %s: localhost:%d\n", portName, port)
	} else {
		fmt.Printf("- localhost:%d\n", port)
	}

}
