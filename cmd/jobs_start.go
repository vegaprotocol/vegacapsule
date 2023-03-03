package cmd

import (
	"context"
	"fmt"
	"log"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/spf13/cobra"
)

var jobsStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a specific job",
	RunE: func(cmd *cobra.Command, args []string) error {
		networkState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return fmt.Errorf("failed load network state: %w", err)
		}

		if networkState.Empty() {
			return networkNotBootstrappedErr("nodes start")
		}

		// try if job name refers to node set
		nodeSet, err := networkState.GeneratedServices.GetNodeSet(jobName)
		if nodeSet != nil && err == nil {
			nomadJobID, err := nodesStartNode(
				context.Background(),
				nodeSet,
				networkState.Config,
				vegaBinary,
			)
			if err != nil {
				return fmt.Errorf("failed start job: %w", err)
			}

			networkState.RunningJobs.NodesSetsJobIDs[nomadJobID] = true

			return networkState.Persist()
		}

		networkState, err = startJob(
			context.Background(),
			*networkState,
			networkState.Config,
			jobName,
		)
		if err != nil {
			return fmt.Errorf("failed start job: %w", err)
		}

		return networkState.Persist()
	},
}

func init() {
	jobsStartCmd.PersistentFlags().StringVar(&jobName,
		"name",
		"",
		"Name of the job to stop",
	)
	jobsStartCmd.MarkPersistentFlagRequired("name")
}

func startJob(ctx context.Context, state state.NetworkState, conf *config.Config, name string) (*state.NetworkState, error) {
	log.Printf("starting %s node set", name)

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, *conf.VegaCapsuleBinary, conf.LogsDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create job runner: %w", err)
	}

	if wallet := state.GeneratedServices.Wallet; wallet != nil && wallet.Name == name {
		log.Printf("starting wallet %s", name)
		jobID, err := nomadRunner.StartWallet(ctx, conf.Network.Wallet, wallet)
		if err != nil {
			return nil, err
		}

		state.RunningJobs.WalletJobID = jobID
	}

	if faucet := state.GeneratedServices.Faucet; faucet != nil && faucet.Name == name {
		log.Printf("starting faucet %s", name)
		jobID, err := nomadRunner.StartFaucet(ctx, conf.Network.Faucet, faucet)
		if err != nil {
			return nil, err
		}

		state.RunningJobs.FaucetJobID = jobID
	}

	for _, e := range conf.Network.PreStart.Exec {
		if e.Name != name {
			continue
		}

		log.Printf("starting job %s", name)
		jobID, err := nomadRunner.RunExecJob(ctx, e)
		if err != nil {
			return nil, err
		}

		state.RunningJobs.AddExtraJobIDs([]string{jobID})
	}

	for _, d := range conf.Network.PreStart.Docker {
		if d.Name != name {
			continue
		}

		log.Printf("starting job %s", name)
		jobID, err := nomadRunner.RunDockerJob(ctx, d)
		if err != nil {
			return nil, err
		}

		state.RunningJobs.AddExtraJobIDs([]string{jobID})
	}

	return &state, nil
}
