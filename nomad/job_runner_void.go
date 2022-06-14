package nomad

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/hashicorp/nomad/api"
)

type VoidJobRunner struct{}

func NewVoidJobRunner() *VoidJobRunner {
	return &VoidJobRunner{}
}

func (r *VoidJobRunner) printCaveat() {
	fmt.Println("Caveat: using void nomad job runner")
}

func (r *VoidJobRunner) RunRawNomadJobs(ctx context.Context, rawJobs []string) ([]*api.Job, error) {
	r.printCaveat()
	return nil, nil
}

func (r *VoidJobRunner) StopNetwork(ctx context.Context, jobs *types.NetworkJobs, nodesOnly bool) error {
	r.printCaveat()
	return nil
}
