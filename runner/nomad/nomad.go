package nomad

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/nomad/api"
)

const (
	DeploymentStatusRunning  = "running"
	DeploymentStatusCanceled = "canceled"
	DeploymentStatusSuccess  = "successful"
	AllocationStateDead      = "dead"
	Running                  = "running"
)

type NomadRunner struct {
	NomadClient *api.Client
}

func New(config *api.Config) (*NomadRunner, error) {
	nomadConfig := api.DefaultConfig()
	if config != nil {
		nomadConfig.Address = config.Address
		nomadConfig.TLSConfig = config.TLSConfig
	}

	cl, err := api.NewClient(nomadConfig)
	if err != nil {
		log.Fatalf("Error creating client %+v", err)
		return nil, err
	}

	r := &NomadRunner{NomadClient: cl}

	return r, nil
}

// TODO maybe improve the logging?
func (n *NomadRunner) waitForDeployment(jobID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(time.Second * 4)
			deployments, _, err := n.NomadClient.Jobs().Deployments(jobID, true, &api.QueryOptions{})
			if err != nil {
				return err
			}

			for _, dep := range deployments {
				log.Printf("deployment (%s) update for job: %q, status: %q, another: %s", dep.ID, dep.JobID, dep.Status, dep.StatusDescription)

				switch dep.Status {
				case DeploymentStatusCanceled:
					return fmt.Errorf("failed to run %s job", jobID)
				case DeploymentStatusSuccess:
					return nil
				}
			}
		}
	}
}

func (n *NomadRunner) Run(job *api.Job) (bool, error) {
	jobs := n.NomadClient.Jobs()

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

func (n *NomadRunner) RunAndWait(job api.Job) error {
	jobs := n.NomadClient.Jobs()

	_, _, err := jobs.Register(&job, nil)
	if err != nil {
		return fmt.Errorf("error running jobs: %w", err)
	}

	if err := n.waitForDeployment(*job.ID); err != nil {
		return err
	}

	return nil
}

func (n *NomadRunner) Stop(jobID string, purge bool) (bool, error) {
	jobs := n.NomadClient.Jobs()

	jId, _, err := jobs.Deregister(jobID, purge, &api.WriteOptions{})
	if err != nil {
		log.Printf("error stopping the job: %+v", err)
		return false, err
	}

	log.Printf("Stopped Job: %+v - %+v", jobID, jId)
	return true, nil
}
