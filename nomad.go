package main

import (
	"log"

	"github.com/hashicorp/nomad/api"
)

const Running = "running"

type NomadRunner struct {
	NomadClient *api.Client
}

func NewNomadRunner(config *api.Config) (*NomadRunner, error) {
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

	return &NomadRunner{NomadClient: cl}, nil
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

	resp, _, err := jobs.Register(job, nil)
	if err != nil {
		log.Fatalf("error registering jobs: %+v", err)
	}
	log.Printf("Success Reponse: %+v", resp)

	return false, nil
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
