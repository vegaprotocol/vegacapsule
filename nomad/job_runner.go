package nomad

import (
	"context"
	"fmt"
	"log"
	"sync"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
	"golang.org/x/sync/errgroup"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec2"
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
	ports := []api.Port{}
	portLabels := []string{}
	if conf.StaticPort != nil {
		ports = append(ports, api.Port{
			Label: fmt.Sprintf("%s-port", conf.Name),
			To:    conf.StaticPort.To,
			Value: conf.StaticPort.Value,
		})
		portLabels = append(portLabels, fmt.Sprintf("%s-port", conf.Name))
	}

	j := &api.Job{
		ID:          &conf.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				Networks: []*api.NetworkResource{
					{
						ReservedPorts: ports,
					},
				},
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
							"image":          conf.Image,
							"command":        conf.Command,
							"args":           conf.Args,
							"ports":          portLabels,
							"auth_soft_fail": conf.AuthSoftFail,
						},
						Env: conf.Env,
						Resources: &api.Resources{
							CPU:      utils.IntPoint(500),
							MemoryMB: utils.IntPoint(768),
						},
					},
				},
			},
		},
	}

	if err := r.client.RunAndWait(ctx, j); err != nil {
		return nil, fmt.Errorf("failed to run nomad docker job: %w", err)
	}

	return j, nil
}

func (r *JobRunner) defaultNodeSetJob(ns types.NodeSet) *api.Job {
	tasks := make([]*api.Task, 0, 2)
	tasks = append(tasks,
		&api.Task{
			Name:   ns.Vega.Name,
			Driver: "raw_exec",
			Config: map[string]interface{}{
				"command": ns.Vega.BinaryPath,
				"args": []string{
					"node",
					"--home", ns.Vega.HomeDir,
					"--tendermint-home", ns.Tendermint.HomeDir,
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
			Name:   ns.DataNode.Name,
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

	return &api.Job{
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
	}
}

func (r *JobRunner) RunRawNomadJobs(ctx context.Context, rawJobs []string) ([]types.RawJobWithNomadJob, error) {
	var mut sync.Mutex
	jobs := make([]types.RawJobWithNomadJob, 0, len(rawJobs))

	eg := new(errgroup.Group)
	for _, rj := range rawJobs {
		rj := rj

		eg.Go(func() error {
			job, err := jobspec2.ParseWithConfig(&jobspec2.ParseConfig{
				Path:    "input.hcl",
				Body:    []byte(rj),
				ArgVars: []string{},
				AllowFS: true,
			})
			if err != nil {
				return fmt.Errorf("failed to parse Nomad job: %w", err)
			}

			if err := r.client.RunAndWait(ctx, job); err != nil {
				return err
			}

			mut.Lock()
			jobs = append(jobs, types.RawJobWithNomadJob{
				RawJob:   rj,
				NomadJob: job,
			})
			mut.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to wait for Nomad jobs: %w", err)
	}

	return jobs, nil

}

func (r *JobRunner) RunNodeSets(ctx context.Context, nodeSets []types.NodeSet) ([]types.NetworkJobState, error) {
	jobs := make([]*api.Job, 0, len(nodeSets))
	jobsStates := []types.NetworkJobState{}

	for _, ns := range nodeSets {
		if ns.RemoteCommandRunner != nil {
			job, err := jobspec2.ParseWithConfig(&jobspec2.ParseConfig{
				Path:    "command_runner.hcl",
				Body:    []byte(ns.RemoteCommandRunner.NomadJobRaw),
				ArgVars: []string{},
				AllowFS: true,
			})

			if err != nil {
				return nil, fmt.Errorf("failed to parse command runner template for node set %s: %w", ns.Name, err)
			}

			jobs = append(jobs, job)
			jobsStates = append(jobsStates, types.NetworkJobState{
				Name:    *job.ID,
				Running: true, // TODO: Is this assumption correct?
				Kind:    types.JobCommandRunner,
			})
		}

		var jobName string
		if ns.NomadJobRaw == nil {
			log.Printf("adding node set %q with default Nomad job definition", ns.Name)

			job := r.defaultNodeSetJob(ns)
			jobs = append(jobs, job)
			jobName = *job.ID
		} else {
			log.Printf("adding node set %q with custom Nomad job definition", ns.Name)
			var err error
			job, err := jobspec2.ParseWithConfig(&jobspec2.ParseConfig{
				Path:    "node_set.hcl",
				Body:    []byte(*ns.NomadJobRaw),
				ArgVars: []string{},
				AllowFS: true,
			})

			if err != nil {
				return nil, err
			}

			jobs = append(jobs, job)
			jobName = *job.ID
		}

		jobsStates = append(jobsStates, types.NetworkJobState{
			Name:    jobName,
			Running: true, // TODO: Is this assumption correct?
			Kind:    types.JobNodeSet,
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

	return jobsStates, nil
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
								"--automatic-consent",
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

	if err := r.client.RunAndWait(ctx, j); err != nil {
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

	if err := r.client.RunAndWait(ctx, j); err != nil {
		return nil, fmt.Errorf("failed to wait for faucet job: %w", err)
	}

	return j, nil
}

func (r *JobRunner) runDockerJobs(ctx context.Context, dockerConfigs []config.DockerConfig) ([]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	jobIDs := make([]string, 0, len(dockerConfigs))
	var jobIDsLock sync.Mutex

	for _, dc := range dockerConfigs {
		// capture in the loop by copy
		dc := dc
		g.Go(func() error {
			job, err := r.runDockerJob(ctx, dc)
			if err != nil {
				return fmt.Errorf("failed to run pre start job %s: %w", dc.Name, err)
			}

			jobIDsLock.Lock()
			jobIDs = append(jobIDs, *job.ID)
			jobIDsLock.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to wait for docker jobs: %w", err)
	}

	return jobIDs, nil
}

func (r *JobRunner) StartNetwork(gCtx context.Context, conf *config.Config, generatedSvcs *types.GeneratedServices) (types.JobStateMap, error) {
	g, ctx := errgroup.WithContext(gCtx)

	result := types.JobStateMap{}
	var lock sync.Mutex

	result.CreateAndAppendJobs(generatedSvcs.PreGenerateJobsIDs(), types.JobPreStart)

	if conf.Network.PreStart != nil {
		ExtraJobs, err := r.runDockerJobs(ctx, conf.Network.PreStart.Docker)
		if err != nil {
			return nil, fmt.Errorf("failed to run pre start jobs: %w", err)
		}

		result.CreateAndAppendJobs(ExtraJobs, types.JobPreStart)
	}

	// create new error group to be able call wait funcion again
	if generatedSvcs.Faucet != nil {
		g.Go(func() error {
			job, err := r.runFaucet(ctx, *conf.VegaBinary, conf.Network.Faucet, generatedSvcs.Faucet)
			if err != nil {
				return fmt.Errorf("failed to run faucet: %w", err)
			}

			lock.Lock()
			result.Append(types.NetworkJobState{
				Name:    *job.ID,
				Running: true,
				Kind:    types.JobFaucet,
			})
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
			result.Append(types.NetworkJobState{
				Name:    *job.ID,
				Kind:    types.JobWallet,
				Running: true,
			})
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
			result.Append(job)
		}
		lock.Unlock()

		return nil
	})

	if err := g.Wait(); err != nil {
		// return result even if job failed to save current network state
		return result, fmt.Errorf("failed to start vega network: %w", err)
	}

	if conf.Network.PostStart != nil {
		extraJobsIDs, err := r.runDockerJobs(gCtx, conf.Network.PostStart.Docker)
		if err != nil {
			// return result even if job failed to save current network state
			return result, fmt.Errorf("failed to run post start jobs: %w", err)
		}

		result.CreateAndAppendJobs(extraJobsIDs, types.JobPostStart)
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

func (r *JobRunner) StopNetwork(ctx context.Context, jobs []types.NetworkJobState) error {
	// no jobs, no network started
	if len(jobs) == 0 {
		return r.stopAllJobs(ctx)

	}

	g, ctx := errgroup.WithContext(ctx)
	for _, jobState := range jobs {
		if jobState.Name == "" {
			continue
		}
		// Explicitly copy name
		jobName := jobState.Name

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
