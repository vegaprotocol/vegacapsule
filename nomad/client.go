package nomad

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/hashicorp/nomad/api"
	"golang.org/x/sync/errgroup"
)

type ConnectionError struct {
	Err error
}

func (ce *ConnectionError) Error() string {
	return fmt.Sprintf("failed to connect to nomad: %s", ce.Err.Error())
}

func newConnectionErr(err error) *ConnectionError {
	return &ConnectionError{
		Err: err,
	}
}

func IsConnectionErr(err error) bool {
	var cerr *ConnectionError
	return errors.As(err, &cerr)
}

const (
	DeploymentStatusRunning  = "running"
	DeploymentStatusCanceled = "cancelled"
	DeploymentStatusSuccess  = "successful"
	AllocationStateDead      = "dead"
	Running                  = "running"
)

type Client struct {
	API *api.Client
}

func NewClient(config *api.Config) (*Client, error) {
	nomadConfig := api.DefaultConfig()
	if config != nil {
		nomadConfig.Address = config.Address
		nomadConfig.TLSConfig = config.TLSConfig
	}

	api, err := api.NewClient(nomadConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	// Ping Nomad
	if _, err := api.Operator().Metrics(nil); err != nil {
		return nil, newConnectionErr(err)
	}

	return &Client{API: api}, nil
}

// TODO maybe improve the logging?
func (n *Client) waitForDeployment(ctx context.Context, jobID string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("failed to run %s job: starting deadline has been exceeded", jobID)
			}
			return ctx.Err()
		default:
			time.Sleep(time.Second * 4)

			job, _, err := n.API.Jobs().Info(jobID, &api.QueryOptions{})
			if err != nil {
				return err
			}

			deployments, _, err := n.API.Jobs().Deployments(jobID, true, &api.QueryOptions{})
			if err != nil {
				return err
			}

			for _, dep := range deployments {
				log.Printf("Update for job: %q, jobStatus: %s, deploymentStatus: %q, another: %s", dep.JobID, *job.Status, dep.Status, dep.StatusDescription)

				switch dep.Status {
				case DeploymentStatusCanceled:
					return fmt.Errorf("failed to run %s job %s", jobID, dep.StatusDescription)
				case DeploymentStatusSuccess:
					return nil
				}
			}
		}
	}
}

func (n *Client) Run(job *api.Job) (bool, error) {
	jobs := n.API.Jobs()

	info, _, err := jobs.Info(*job.ID, &api.QueryOptions{})
	if err != nil {
		//NOTE: Handle 404 status code
		log.Printf("Error getting job info: %+v", err)
	} else if *info.Status == Running {
		return true, nil
	}

	_, _, err = jobs.Register(job, nil)
	if err != nil {
		log.Fatalf("error registering jobs: %+v", err)
	}

	return false, nil
}

func (n *Client) logJob(ctx context.Context, job *api.Job) error {

	// /client/allocation/:alloc_id/gc
	jobs := n.API.Jobs()
	allocsApi := n.API.AllocFS()

	var err error
	var allocs []*api.AllocationListStub

	for len(allocs) == 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			allocs, _, err = jobs.Allocations(*job.ID, true, &api.QueryOptions{})
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Logs for", *job.ID, allocs)

	eg, egCtx := errgroup.WithContext(ctx)

	for _, as := range allocs {
		file, err := os.Create(path.Join("/Users/karel/.vegacapsule/testnet", as.Name))
		if err != nil {
			return err
		}

		fmt.Println("alloc: ", as.NodeID, as.Name, file.Name())
		cAlloc := &api.Allocation{ID: as.ID, NodeID: as.NodeID}
		for taskName := range as.TaskStates {

			// n.API.AllocFS().Logs()

			for _, logType := range logTypes {
				taskName := taskName
				logType := logType
				file := file
				eg.Go(func() error {
					defer file.Close()

					cancelCh := make(chan struct{})
					framesCh, errsCh := allocsApi.Logs(cAlloc, true, taskName, logType, "start", 0, cancelCh, &api.QueryOptions{})
					fmt.Println("starting ", taskName, logType)
					for {
						select {
						case <-egCtx.Done():
							return egCtx.Err()
						case frame := <-framesCh:
							if frame.IsHeartbeat() {
								break
							}

							fmt.Printf("---- writing to file: %s", file.Name())

							if _, err := file.Write(frame.Data); err != nil {
								return err
							}
						case err := <-errsCh:
							return err
						}
					}
				})

			}
		}
	}

	return eg.Wait()
}

func (n *Client) RunAndWait(ctx context.Context, job *api.Job) error {
	jobs := n.API.Jobs()

	_, _, err := jobs.Register(job, nil)
	if err != nil {
		return fmt.Errorf("error running jobs: %w", err)
	}

	// go func() {
	// 	log.Printf("Starting log watcher for job %s", *job.ID)

	// 	if err := n.logJob(ctx, job); err != nil {
	// 		log.Printf("Failed to log job %s: %s", *job.ID, err)
	// 	}
	// }()

	if err := n.waitForDeployment(ctx, *job.ID); err != nil {
		return err
	}

	return nil
}

// Stop stops a specific job
func (n *Client) Stop(ctx context.Context, jobID string, purge bool) error {
	jobs := n.API.Jobs()

	writeOpts := new(api.WriteOptions).WithContext(ctx)
	jId, _, err := jobs.Deregister(jobID, purge, writeOpts)
	if err != nil {
		log.Printf("error stopping the job: %+v", err)
		return err
	}

	log.Printf("Stopped Job: %+v - %+v", jobID, jId)
	return nil
}
