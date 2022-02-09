package nomad

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/nomad/api"
)

const (
	DeploymentStatusCanceled = "canceled"
	DeploymentStatusSuccess  = "successful"
	AllocationStateDead      = "dead"
)
const Running = "running"

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

// TODO this needs some love. Sometimes the function returns and error even though the job is running.....
func (n *NomadRunner) waitForDeployment(jobID string, deployID string) error {
	log.Printf("waiting for job: %q", jobID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	topics := map[api.Topic][]string{
		api.TopicDeployment: {deployID},
		api.TopicAllocation: {deployID},
	}

	eventCh, err := n.NomadClient.EventStream().Stream(ctx, topics, 0, &api.QueryOptions{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-eventCh:
			if event.Err != nil {
				log.Println("error from event stream", "error", err)
				break
			}
			if event.IsHeartbeat() {
				continue
			}

			for _, e := range event.Events {
				switch e.Topic {
				case api.TopicDeployment:
					dep, _ := e.Deployment()
					log.Printf("deployment (%s) update for job: %q, status: %q, another: %s", dep.ID, dep.JobID, dep.Status, dep.StatusDescription)

					switch dep.Status {
					case DeploymentStatusCanceled:
						return fmt.Errorf("failed to run %s job", jobID)
					case DeploymentStatusSuccess:
						return nil
					}
				case api.TopicAllocation:
					alloc, _ := e.Allocation()

					for _, ts := range alloc.TaskStates {
						for _, tse := range ts.Events {
							log.Printf("update for jobID %q, state: %q, type: %q, reason %q", jobID, ts.State, tse.Type, tse.DisplayMessage)
						}

						if ts.State == AllocationStateDead {
							return fmt.Errorf("failed to run %q job", jobID)
						}
					}

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

func (n *NomadRunner) RunAndWait(job *api.Job) error {
	jobs := n.NomadClient.Jobs()

	resp, _, err := jobs.Register(job, nil)
	if err != nil {
		return fmt.Errorf("error running jobs: %w", err)
	}

	eval, _, err := n.NomadClient.Evaluations().Info(resp.EvalID, nil)
	if err != nil {
		log.Fatalf("error registering jobs: %+v", err)
	}

	if err := n.waitForDeployment(eval.JobID, eval.DeploymentID); err != nil {
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
