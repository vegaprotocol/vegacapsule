package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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

func (n *NomadRunner) Stop(job *api.Job, purge bool) (bool, error) {
	jobs := n.NomadClient.Jobs()

	jId, _, err := jobs.Deregister(*job.ID, purge, &api.WriteOptions{})
	if err != nil {
		log.Printf("error stopping the job: %+v", err)
		return false, err
	}

	log.Printf("Stopped Job: %+v - %+v", *job.Name, jId)
	return true, nil
}

func ganacheCheck(timeout time.Duration) error {
	for start := time.Now(); time.Since(start) < timeout; {
		time.Sleep(1 * time.Second)
		postBody, _ := json.Marshal(map[string]string{
			"method": "web3_clientVersion",
		})
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post("http://127.0.0.1:8545/", "application/json", responseBody)
		if err != nil {
			log.Println("ganache not yet ready", err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		if strings.Contains(string(body), "EthereumJS") {
			log.Println("ganache is ready")
			return nil
		}
		continue
	}

	return fmt.Errorf("ganache container has timed out")
}
