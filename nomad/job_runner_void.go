package nomad

import (
	"context"
	"fmt"

	"github.com/hashicorp/nomad/api"
)

type VoidJobRunner struct{}

func NewVoidJobRunner() *VoidJobRunner {
	return &VoidJobRunner{}
}

func (r *VoidJobRunner) RunRawNomadJobs(ctx context.Context, rawJobs []string) ([]*api.Job, error) {
	fmt.Println("Caveat: using void nomad job runner")
	return nil, nil
}
