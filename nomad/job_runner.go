package nomad

import (
	"context"
	"fmt"
	"sync"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"golang.org/x/sync/errgroup"

	"github.com/hashicorp/nomad/api"
)

type JobRunner struct {
	client *Client
}

func NewJobRunner(c *Client) *JobRunner {
	return &JobRunner{
		client: c,
	}
}

func (r *JobRunner) runDockerJob(ctx context.Context, conf config.DockerConfig) (*api.Job, error) {
	j := nomadDockerJob(conf)

	if err := r.client.RunAndWait(ctx, *j); err != nil {
		return nil, fmt.Errorf("failed to run nomad docker job: %w", err)
	}

	return j, nil
}

func (r *JobRunner) RunNodeSets(ctx context.Context, nodeSets []types.NodeSet) ([]api.Job, error) {
	jobs, err := nomadNodeSetJobs(nodeSets)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nomad jobs for node sets: %w", err)
	}

	eg := new(errgroup.Group)
	for _, j := range jobs {
		j := j

		eg.Go(func() error {
			return r.client.RunAndWait(ctx, j)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to wait for node sets: %w", err)
	}

	return jobs, nil
}

func (r *JobRunner) runWallet(ctx context.Context, conf *config.WalletConfig, wallet *types.Wallet) (*api.Job, error) {
	j := nomadWalletJob(conf, wallet)

	if err := r.client.RunAndWait(ctx, *j); err != nil {
		return nil, fmt.Errorf("failed to run the wallet job: %w", err)
	}

	return j, nil
}

func (r *JobRunner) runFaucet(ctx context.Context, binary string, conf *config.FaucetConfig, fc *types.Faucet) (*api.Job, error) {
	j := nomadFaucetJob(binary, conf, fc)

	if err := r.client.RunAndWait(ctx, *j); err != nil {
		return nil, fmt.Errorf("failed to wait for faucet job: %w", err)
	}

	return j, nil
}

func (r *JobRunner) StartNetwork(gCtx context.Context, conf *config.Config, generatedSvcs *types.GeneratedServices) (*types.NetworkJobs, error) {
	g, ctx := errgroup.WithContext(gCtx)
	result := &types.NetworkJobs{
		NodesSetsJobIDs: map[string]bool{},
		ExtraJobIDs:     map[string]bool{},
	}
	var lock sync.Mutex

	for _, dc := range conf.Network.PreStart.Docker {
		// capture in the loop by copy
		dc := dc
		g.Go(func() error {
			job, err := r.runDockerJob(ctx, dc)
			if err != nil {
				return fmt.Errorf("failed to run pre start job %s: %w", dc.Name, err)
			}

			lock.Lock()

			result.ExtraJobIDs[*job.ID] = true
			lock.Unlock()

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to wait for pre-start jobs: %w", err)
	}

	// create new error group to be able call wait funcion again
	g, ctx = errgroup.WithContext(gCtx)
	if generatedSvcs.Faucet != nil {
		g.Go(func() error {
			job, err := r.runFaucet(ctx, *conf.VegaBinary, conf.Network.Faucet, generatedSvcs.Faucet)
			if err != nil {
				return fmt.Errorf("failed to run faucet: %w", err)
			}

			lock.Lock()
			result.FaucetJobID = *job.ID
			lock.Unlock()

			return nil
		})
	}

	if generatedSvcs.Wallet != nil {
		g.Go(func() error {
			job, err := r.runWallet(ctx, conf.Network.Wallet, generatedSvcs.Wallet)
			if err != nil {
				return fmt.Errorf("failed to run wallet: %w", err)
			}

			lock.Lock()
			result.WalletJobID = *job.ID
			lock.Unlock()

			return nil
		})
	}

	g.Go(func() error {
		jobs, err := r.RunNodeSets(ctx, generatedSvcs.NodeSets.ToSlice())
		if err != nil {
			return fmt.Errorf("failed to run node sets: %w", err)
		}

		lock.Lock()
		for _, job := range jobs {
			result.NodesSetsJobIDs[*job.ID] = true
		}
		lock.Unlock()

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to start vega network: %w", err)
	}
	return result, nil
}

func (r *JobRunner) stopAllJobs(ctx context.Context) error {
	jobs, _, err := r.client.API.Jobs().List(nil)
	if err != nil {
		return err
	}

	var eg errgroup.Group
	for _, j := range jobs {
		j := j
		eg.Go(func() error {
			_, err := r.client.Stop(ctx, j.ID, true)
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	// just to try - we are not interested in error
	_ = r.client.API.System().GarbageCollect()

	return nil
}

func (r *JobRunner) StopNetwork(ctx context.Context, jobs *types.NetworkJobs, nodesOnly bool) error {
	// no jobs, no network started
	if jobs == nil {
		if !nodesOnly {
			return r.stopAllJobs(ctx)
		}

		return nil
	}

	allJobs := []string{}
	if !nodesOnly {
		allJobs = append(jobs.ExtraJobIDs.ToSlice(), jobs.WalletJobID, jobs.FaucetJobID)
	}
	allJobs = append(allJobs, jobs.NodesSetsJobIDs.ToSlice()...)
	g, ctx := errgroup.WithContext(ctx)
	for _, jobName := range allJobs {
		if jobName == "" {
			continue
		}
		// Explicitly copy name
		jobName := jobName

		g.Go(func() error {
			if _, err := r.client.Stop(ctx, jobName, true); err != nil {
				return fmt.Errorf("cannot stop nomad job \"%s\": %w", jobName, err)
			}
			return nil
		})

	}

	if err := g.Wait(); err != nil {
		return err
	}

	// just to try - we are not interested in error
	_ = r.client.API.System().GarbageCollect()

	return nil
}

func (r *JobRunner) StopJobs(ctx context.Context, jobIDs []string) error {
	// no jobs to stop
	if len(jobIDs) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, jobName := range jobIDs {
		if jobName == "" {
			continue
		}
		// Explicitly copy name
		jobName := jobName

		g.Go(func() error {
			if _, err := r.client.Stop(ctx, jobName, true); err != nil {
				return fmt.Errorf("cannot stop nomad job \"%s\": %w", jobName, err)
			}
			return nil
		})

	}

	return g.Wait()
}
