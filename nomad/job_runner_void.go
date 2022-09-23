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

func (r *VoidJobRunner) RunRawNomadJobs(ctx context.Context, rawJobs []string) ([]types.RawJobWithNomadJob, error) {
	r.printCaveat()
	return nil, nil
}

func (r *VoidJobRunner) StopNetwork(ctx context.Context, jobs *types.NetworkJobs, nodesOnly bool) error {
	r.printCaveat()
	return nil
}

func (r *VoidJobRunner) GetJobPorts(job *api.Job) []int64 {
	r.printCaveat()
	return nil
}
