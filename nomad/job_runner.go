package nomad

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"
	"sync"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/logscollector"
	"code.vegaprotocol.io/vegacapsule/types"
	"golang.org/x/sync/errgroup"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec2"
)

type JobRunner struct {
	Client        *Client
	capsuleBinary string
	logsOutputDir string
}

func NewJobRunner(c *Client, capsuleBinaryPath, logsOutputDir string) (*JobRunner, error) {
	return &JobRunner{
		Client:        c,
		capsuleBinary: capsuleBinaryPath,
		logsOutputDir: logsOutputDir,
	}, nil
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

			if err := r.runAndWait(ctx, job, nil); err != nil {
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

func (r *JobRunner) runAndWait(
	ctx context.Context,
	job *api.Job,
	probes *types.ProbesConfig,
) error {
	err := r.Client.RunAndWait(ctx, job, probes)
	if err == nil {
		return nil
	}

	if (IsJobTimeoutErr(err)) && hasLogsCollectorTask(job) {
		fmt.Printf("\nLogs from failed %s job:\n", *job.ID)

		if err := logscollector.TailLastLogs(path.Join(r.logsOutputDir, *job.ID)); err != nil {
			return fmt.Errorf("failed to print logs from failed job: %w", err)
		}
	}

	return err
}

type jobWithPreProbe struct {
	Job    *api.Job
	Probes *types.ProbesConfig
}

// RunNodeSets returns list of started jobs.
// When one of the jobs fails during startup it returns jobs that has alredy started before that.
func (r *JobRunner) RunNodeSets(ctx context.Context, nodeSets []types.NodeSet) ([]*api.Job, error) {
	jobs := make([]jobWithPreProbe, 0, len(nodeSets))

	for _, ns := range nodeSets {
		if ns.NomadJobRaw == nil {
			log.Printf("adding node set %q with default Nomad job definition", ns.Name)

			jobs = append(jobs, jobWithPreProbe{Job: r.defaultNodeSetJob(ns), Probes: ns.PreStartProbe})
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

		jobs = append(jobs, jobWithPreProbe{Job: job, Probes: ns.PreStartProbe})
	}

	var mut sync.Mutex
	startedJobs := make([]*api.Job, 0, len(nodeSets))

	eg := new(errgroup.Group)
	for _, j := range jobs {
		j := j

		eg.Go(func() error {
			if err := r.runAndWait(ctx, j.Job, j.Probes); err != nil {
				if _, err := r.stopJobsByIDs(ctx, []string{*j.Job.ID}); err != nil {
					log.Printf("Failed to stop failed job %s, this job might need to be stopped manually. Reason %s", *j.Job.ID, err)
				}

				return err
			}

			mut.Lock()
			startedJobs = append(startedJobs, j.Job)
			mut.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return startedJobs, fmt.Errorf("failed to wait for node sets: %w", err)
	}

	return startedJobs, nil
}

func (r *JobRunner) runDockerJobs(ctx context.Context, dockerConfigs []config.DockerConfig) ([]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	jobIDs := make([]string, 0, len(dockerConfigs))

	var jobIDsLock sync.Mutex

	for _, dc := range dockerConfigs {
		// capture in the loop by copy
		dc := dc
		g.Go(func() error {
			// Skip for already running jobs
			if r.Client.JobRunning(ctx, dc.Name) {
				return nil
			}

			job := r.defaultDockerJob(ctx, dc)

			if err := r.runAndWait(ctx, job, nil); err != nil {
				return fmt.Errorf("failed to run pre start job %q: %w", *job.ID, err)
			}

			jobIDsLock.Lock()
			jobIDs = append(jobIDs, *job.ID)
			jobIDsLock.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return jobIDs, nil
}

func (r *JobRunner) runExecJobs(ctx context.Context, execConfigs []config.ExecConfig) ([]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	jobIDs := make([]string, 0, len(execConfigs))

	var jobIDsLock sync.Mutex

	for _, ec := range execConfigs {
		// capture in the loop by copy
		ec := ec
		g.Go(func() error {
			// Skip for already running jobs
			if r.Client.JobRunning(ctx, ec.Name) {
				return nil
			}

			job := r.defaultExecJob(ctx, ec)

			if err := r.runAndWait(ctx, job, nil); err != nil {
				return fmt.Errorf("failed to run pre start job %q: %w", *job.ID, err)
			}

			jobIDsLock.Lock()
			jobIDs = append(jobIDs, *job.ID)
			jobIDsLock.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return jobIDs, nil
}

func (r *JobRunner) StartNetwork(
	ctx context.Context,
	conf *config.Config,
	generatedSvcs *types.GeneratedServices,
	stopAllJobsOnFailure bool,
) (*types.NetworkJobs, error) {
	netJobs, err := r.startNetwork(ctx, conf, generatedSvcs)
	if err != nil {
		if stopAllJobsOnFailure {
			if _, err := r.stopAllJobs(ctx); err != nil {
				log.Printf("Failed to stop all registered jobs - please clean up Nomad manually: %s", err)
			}
			return nil, err
		}
		log.Println("Part of the network could not start, but it has been required to not stop existing jobs on failure, so we continue as normal...")
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
		extraJobIDs, err := r.runExecJobs(ctx, conf.Network.PreStart.Exec)
		if err != nil {
			return result, fmt.Errorf("failed to run pre start exec jobs: %w", err)
		}
		result.AddExtraJobIDs(extraJobIDs)

		extraJobIDs, err = r.runDockerJobs(ctx, conf.Network.PreStart.Docker)
		if err != nil {
			return result, fmt.Errorf("failed to run pre start docker jobs: %w", err)
		}
		result.AddExtraJobIDs(extraJobIDs)
	}

	// create new error group to be able to call the `wait` function again
	if generatedSvcs.Faucet != nil {
		g.Go(func() error {
			// Skip for already running faucet
			if r.Client.JobRunning(ctx, generatedSvcs.Faucet.Name) {
				return nil
			}

			job := r.defaultFaucetJob(conf.Network.Faucet, generatedSvcs.Faucet)

			if err := r.runAndWait(ctx, job, nil); err != nil {
				return fmt.Errorf("failed to run the faucet job %q: %w", *job.ID, err)
			}

			lock.Lock()
			result.FaucetJobID = *job.ID
			lock.Unlock()

			return nil
		})
	}

	if generatedSvcs.Wallet != nil {
		g.Go(func() error {
			// Skip for already running wallet
			if r.Client.JobRunning(ctx, generatedSvcs.Wallet.Name) {
				return nil
			}

			job := r.defaultWalletJob(generatedSvcs.Wallet)

			if err := r.runAndWait(ctx, job, nil); err != nil {
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

		lock.Lock()
		for _, job := range jobs {
			result.NodesSetsJobIDs[*job.ID] = true
		}
		lock.Unlock()

		if err != nil {
			return fmt.Errorf("failed to run node sets: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return result, fmt.Errorf("failed to start vega network: %w", err)
	}

	if conf.Network.PostStart != nil {
		extraJobIDs, err := r.runDockerJobs(gCtx, conf.Network.PostStart.Docker)
		if err != nil {
			return result, fmt.Errorf("failed to run post start jobs: %w", err)
		}

		result.AddExtraJobIDs(extraJobIDs)
	}

	return result, nil
}

func (r *JobRunner) stopAllJobs(ctx context.Context) ([]string, error) {
	allJobs, _, err := r.Client.API.Jobs().List(nil)
	if err != nil {
		return nil, err
	}

	allJobIDs := []string{}
	for _, job := range allJobs {
		if job.ID == "" {
			continue
		}
		allJobIDs = append(allJobIDs, job.ID)
	}

	return r.stopJobsByIDs(ctx, allJobIDs)
}

func (r *JobRunner) StopNetwork(ctx context.Context, jobs *types.NetworkJobs, nodesOnly bool) ([]string, error) {
	// no jobs, no network started
	if jobs == nil {
		if !nodesOnly {
			return r.stopAllJobs(ctx)
		}

		return nil, nil
	}

	allJobIDs := []string{}
	if !nodesOnly {
		allJobIDs = append(jobs.ExtraJobIDs.ToSlice(), jobs.WalletJobID, jobs.FaucetJobID)
	}
	allJobIDs = append(allJobIDs, jobs.NodesSetsJobIDs.ToSlice()...)

	return r.stopJobsByIDs(ctx, allJobIDs)
}

func (r *JobRunner) StopJobs(ctx context.Context, jobIDs []string) ([]string, error) {
	return r.stopJobsByIDs(ctx, jobIDs)
}

// ListExposedPortsPerJob returns exposed ports per node
func (r *JobRunner) ListExposedPortsPerJob(ctx context.Context, jobID string) ([]int64, error) {
	job, err := r.Client.Info(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return GetJobPorts(job), nil
}

// ListExposedPorts returns exposed ports across all nodes
func (r *JobRunner) ListExposedPorts(ctx context.Context) (map[string][]int64, error) {
	jobs, err := r.Client.List(ctx)
	if err != nil {
		return nil, err
	}

	portsPerJob := map[string][]int64{}

	for _, j := range jobs {
		if j.Status != Running {
			continue
		}

		ports, err := r.ListExposedPortsPerJob(ctx, j.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list ports for job %q: %w", j.ID, err)
		}

		portsPerJob[j.ID] = ports
	}

	return portsPerJob, nil
}

func (r *JobRunner) stopJobsByIDs(ctx context.Context, allJobIDs []string) ([]string, error) {
	// Apparently, we can have blank job IDs, so skipping them.
	cleanedUpJobIDs := []string{}
	for _, jobID := range allJobIDs {
		if jobID == "" {
			continue
		}
		cleanedUpJobIDs = append(cleanedUpJobIDs, jobID)
	}

	if len(cleanedUpJobIDs) == 0 {
		log.Println("No job to be stopped.")
		return nil, nil
	}

	log.Printf("Trying to stop jobs: %s\n", strings.Join(cleanedUpJobIDs, ", "))

	g, ctx := errgroup.WithContext(ctx)
	for _, jobID := range cleanedUpJobIDs {
		cpyJobID := jobID
		g.Go(func() error {
			if err := r.Client.Stop(ctx, cpyJobID, true); err != nil {
				return fmt.Errorf("cannot stop job %q: %w", cpyJobID, err)
			}

			log.Printf("Job %q stopped\n", cpyJobID)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("could not stop all jobs: %w", err)
	}

	log.Println("Jobs have been stopped.")

	// just to try - we are not interested in error
	_ = r.Client.API.System().GarbageCollect()

	return cleanedUpJobIDs, nil
}
