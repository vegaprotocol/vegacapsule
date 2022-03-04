package nomad

import (
	"context"
	"fmt"
	"sync"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
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
	j := &api.Job{
		ID:          &conf.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.IntPoint(0),
					Mode:     utils.StrPoint("fail"),
				},
				Name: &conf.Name,
				Tasks: []*api.Task{
					{
						Name:   conf.Name,
						Driver: "docker",
						Config: map[string]interface{}{
							"image":   conf.Image,
							"command": conf.Command,
							"args":    conf.Args,
						},
						Resources: &api.Resources{
							Networks: []*api.NetworkResource{
								{
									ReservedPorts: []api.Port{
										{
											Label: fmt.Sprintf("%s-port", conf.Name),
											Value: conf.StaticPort,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if err := r.client.RunAndWait(ctx, *j); err != nil {
		return nil, fmt.Errorf("failed to run nomad docker job: %w", err)
	}

	return j, nil
}

func (r *JobRunner) RunNodeSets(ctx context.Context, vegaBinary string, nodeSets []types.NodeSet) ([]api.Job, error) {
	jobs := make([]api.Job, 0, len(nodeSets))

	for i, ns := range nodeSets {
		tasks := make([]*api.Task, 0, 3)
		tasks = append(tasks,
			&api.Task{
				Name:   fmt.Sprintf("tendermint-%d", i),
				Driver: "raw_exec",
				Config: map[string]interface{}{
					"command": vegaBinary,
					"args": []string{
						"tm",
						"node",
						"--home", ns.Tendermint.HomeDir,
					},
				},
				Resources: &api.Resources{
					CPU:      utils.IntPoint(500),
					MemoryMB: utils.IntPoint(512),
				},
			},
			&api.Task{
				Name:   fmt.Sprintf("vega-%s-%d", ns.Mode, i),
				Driver: "raw_exec",
				Config: map[string]interface{}{
					"command": vegaBinary,
					"args": []string{
						"node",
						"--home", ns.Vega.HomeDir,
						"--nodewallet-passphrase-file", ns.Vega.NodeWalletPassFilePath,
					},
				},
				Resources: &api.Resources{
					CPU:      utils.IntPoint(500),
					MemoryMB: utils.IntPoint(512),
				},
			})

		if ns.DataNode != nil {
			tasks = append(tasks, &api.Task{
				Name:   fmt.Sprintf("data-node-%s-%d", ns.Mode, i),
				Driver: "raw_exec",
				Config: map[string]interface{}{
					"command": ns.DataNode.BinaryPath,
					"args": []string{
						"node",
						"--home", ns.DataNode.HomeDir,
					},
				},
				Resources: &api.Resources{
					CPU:      utils.IntPoint(500),
					MemoryMB: utils.IntPoint(512),
				},
			})
		}

		jobs = append(jobs, api.Job{
			ID:          utils.StrPoint(ns.Name),
			Datacenters: []string{"dc1"},
			TaskGroups: []*api.TaskGroup{
				{
					RestartPolicy: &api.RestartPolicy{
						Attempts: utils.IntPoint(0),
						Mode:     utils.StrPoint("fail"),
					},
					Name:  utils.StrPoint("vega"),
					Tasks: tasks,
				},
			},
		})
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
	j := &api.Job{
		ID:          &wallet.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.IntPoint(0),
					Mode:     utils.StrPoint("fail"),
				},
				Name: utils.StrPoint("vega"),
				Tasks: []*api.Task{
					{
						Name:   "wallet-1",
						Driver: "raw_exec",
						Config: map[string]interface{}{
							"command": conf.Binary,
							"args": []string{
								"service",
								"run",
								"--network", wallet.Network,
								"--no-version-check",
								"--output", "json",
								"--home", wallet.HomeDir,
							},
						},
						Resources: &api.Resources{
							CPU:      utils.IntPoint(500),
							MemoryMB: utils.IntPoint(512),
						},
					},
				},
			},
		},
	}

	if err := r.client.RunAndWait(ctx, *j); err != nil {
		return nil, fmt.Errorf("failed to run the wallet job: %w", err)
	}

	return j, nil
}

func (r *JobRunner) runFaucet(ctx context.Context, binary string, conf *config.FaucetConfig, fc *types.Faucet) (*api.Job, error) {
	j := &api.Job{
		ID:          &fc.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.IntPoint(0),
					Mode:     utils.StrPoint("fail"),
				},
				Name: &conf.Name,
				Tasks: []*api.Task{
					{
						Name:   conf.Name,
						Driver: "raw_exec",
						Config: map[string]interface{}{
							"command": binary,
							"args": []string{
								"faucet",
								"run",
								"--passphrase-file", fc.WalletPassFilePath,
								"--home", fc.HomeDir,
							},
						},
						Resources: &api.Resources{
							CPU:      utils.IntPoint(500),
							MemoryMB: utils.IntPoint(512),
						},
					},
				},
			},
		},
	}

	if err := r.client.RunAndWait(ctx, *j); err != nil {
		return nil, fmt.Errorf("failed to wait for faucet job: %w", err)
	}

	return j, nil
}

func (r *JobRunner) StartNetwork(gCtx context.Context, conf *config.Config, generatedSvcs *types.GeneratedServices) (*types.NetworkJobs, error) {
	g, ctx := errgroup.WithContext(gCtx)
	result := &types.NetworkJobs{}
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
			result.ExtraJobIDs = append(result.ExtraJobIDs, *job.ID)
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
		jobs, err := r.RunNodeSets(ctx, *conf.VegaBinary, generatedSvcs.NodeSets)
		if err != nil {
			return fmt.Errorf("failed to run node sets: %w", err)
		}

		lock.Lock()
		for _, job := range jobs {
			result.NodesSetsJobIDs = append(result.NodesSetsJobIDs, *job.ID)
		}
		lock.Unlock()

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to start vega network: %w", err)
	}
	return result, nil
}

func (r *JobRunner) StopNetwork(ctx context.Context, jobs *types.NetworkJobs) error {
	// no jobs, no network started
	if jobs == nil {
		return nil
	}

	allJobs := append(jobs.ExtraJobIDs, jobs.WalletJobID, jobs.FaucetJobID)
	allJobs = append(allJobs, jobs.NodesSetsJobIDs...)
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

	return g.Wait()
}
