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
					Attempts: utils.ToPoint(0),
					Mode:     utils.ToPoint("fail"),
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
							CPU:      utils.ToPoint(500),
							MemoryMB: utils.ToPoint(768),
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

func (r *JobRunner) RunNodeSets(ctx context.Context, nodeSets []types.NodeSet) ([]*api.Job, error) {
	jobs := make([]*api.Job, 0, len(nodeSets))

	for _, ns := range nodeSets {
		if ns.NomadJobRaw == nil {
			log.Printf("adding node set %q with default Nomad job definition", ns.Name)

			jobs = append(jobs, r.defaultNodeSetJob(ns))
			continue
		}

		job, err := jobspec2.ParseWithConfig(&jobspec2.ParseConfig{
			Path:    "input.hcl",
			Body:    []byte(*ns.NomadJobRaw),
			ArgVars: []string{},
			AllowFS: true,
		})

		if err != nil {
			return nil, err
		}

		log.Printf("adding node set %q with custom Nomad job definition", ns.Name)

		jobs = append(jobs, job)
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

func (r *JobRunner) StartNetwork(
	ctx context.Context,
	conf *config.Config,
	generatedSvcs *types.GeneratedServices,
) (*types.NetworkJobs, error) {
	netJobs, err := r.startNetwork(ctx, conf, generatedSvcs)
	if err != nil {
		if err := r.stopAllJobs(ctx); err != nil {
			log.Printf("Failed to stop all registered jobs - please clean up Nomad manually: %s", err)
		}

		return nil, err
	}

	return netJobs, nil
}

func (r *JobRunner) startNetwork(
	gCtx context.Context,
	conf *config.Config,
	generatedSvcs *types.GeneratedServices,
) (*types.NetworkJobs, error) {
	g, ctx := errgroup.WithContext(gCtx)

	result := &types.NetworkJobs{
		NodesSetsJobIDs: map[string]bool{},
		ExtraJobIDs:     map[string]bool{},
	}
	var lock sync.Mutex

	result.AddExtraJobIDs(generatedSvcs.PreGenerateJobsIDs())

	if conf.Network.PreStart != nil {
		extraJobIDs, err := r.runDockerJobs(ctx, conf.Network.PreStart.Docker)
		if err != nil {
			return nil, fmt.Errorf("failed to run pre start jobs: %w", err)
		}

		result.AddExtraJobIDs(extraJobIDs)
	}

	// create new error group to be able call wait funcion again
	if generatedSvcs.Faucet != nil {
		g.Go(func() error {
			job := r.defaultFaucetJob(*conf.VegaBinary, conf.Network.Faucet, generatedSvcs.Faucet)
			if err := r.client.RunAndWait(ctx, job); err != nil {
				return fmt.Errorf("failed to run the fauce job %q: %w", *job.ID, err)
			}

			lock.Lock()
			result.FaucetJobID = *job.ID
			lock.Unlock()

			return nil
		})
	}

	if generatedSvcs.Wallet != nil {
		g.Go(func() error {
			job := r.defaultWalletJob(conf.Network.Wallet, generatedSvcs.Wallet)
			if err := r.client.RunAndWait(ctx, job); err != nil {
				return fmt.Errorf("failed to run the wallet job %q: %w", *job.ID, err)
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

	if conf.Network.PostStart != nil {
		extraJobIDs, err := r.runDockerJobs(gCtx, conf.Network.PostStart.Docker)
		if err != nil {
			return nil, fmt.Errorf("failed to run post start jobs: %w", err)
		}

		result.AddExtraJobIDs(extraJobIDs)
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
			return r.client.Stop(ctx, j.ID, true)
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
			if err := r.client.Stop(ctx, jobName, true); err != nil {
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
			if err := r.client.Stop(ctx, jobName, true); err != nil {
				return fmt.Errorf("cannot stop nomad job \"%s\": %w", jobName, err)
			}
			return nil
		})

	}

	return g.Wait()
}
