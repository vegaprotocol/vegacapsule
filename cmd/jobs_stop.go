package cmd

import (
	"context"
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var jobName string

var jobsStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a specific job",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("jobs stop")
		}

		updatedNetworkState, err := stopJob(context.Background(), *networkState, jobName)
		if err != nil {
			return fmt.Errorf("failed stop job: %w", err)
		}

		return updatedNetworkState.Persist()
	},
}

func init() {
	jobsStopCmd.PersistentFlags().StringVar(&jobName,
		"name",
		"",
		"Name of the job to stop",
	)
	jobsStopCmd.MarkPersistentFlagRequired("name")
}

func stopJob(ctx context.Context, state state.NetworkState, name string) (*state.NetworkState, error) {
	log.Printf("stopping %s job", name)

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, *state.Config.VegaCapsuleBinary, state.Config.LogsDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create job runner: %w", err)
	}

	toRemove := []string{name}
	stoppedJobs, err := nomadRunner.StopJobs(ctx, toRemove)
	if err != nil {
		return nil, fmt.Errorf("failed to stop nomad job %q: %w", name, err)
	}

	if state.RunningJobs != nil {
		state.RunningJobs.RemoveRunningJobsIDs(stoppedJobs)
	}

	log.Printf("stopping %s job success", name)
	return &state, nil
}
