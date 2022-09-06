package nomad

import "github.com/hashicorp/nomad/api"

func GetJobPorts(job *api.Job) []int64 {
	var ports []int64

	for _, tg := range job.TaskGroups {
		for _, net := range tg.Networks {
			for _, p := range append(net.DynamicPorts, net.ReservedPorts...) {
				ports = append(ports, int64(p.Value))
			}
		}
	}

	return ports
}
