package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
)

func registerJobs(client *api.Client, dir string) error {

	files, err := os.ReadDir(dir)
	path, _ := filepath.Abs(dir)

	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, file := range files {
		p := filepath.Join(path, file.Name())

		j, err := jobspec.ParseFile(p)
		if err != nil {
			fmt.Println(err)
			return err
		}
		client.Jobs().Register(j, &api.WriteOptions{})
	}
	return nil
}

func deregisterJobs(client *api.Client) error {
	jobs, _, err := client.Jobs().List(nil)
	if err != nil {
		log.Fatalf("[ERR] nomad: failed listing jobs: %v", err)
		return err
	}
	fmt.Println(jobs)
	for _, job := range jobs {
		if _, _, err := client.Jobs().Deregister(job.ID, true, nil); err != nil {
			log.Fatalf("[ERR] nomad: failed deregistering job: %v", err)
			return err
		}
	}
	return nil
}
