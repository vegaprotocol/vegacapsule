package ports

import (
	"context"
	"fmt"
	"log"
	"strings"

	"code.vegaprotocol.io/vegacapsule/types"
)

type JobWithTask struct {
	Name     string
	TaskName string
}

type PortWithName struct {
	Name string
	Port int64
}

type genServiceGetter interface {
	GetByName(name string) []types.GeneratedService
}

func OpenPortsPerJob(
	ctx context.Context,
	nomadExposedPortsPerJob map[string][]int64,
	genServices genServiceGetter,
) (map[JobWithTask][]PortWithName, error) {
	portsPerProcessName, err := ScanNetworkProcessesPorts()
	if err != nil {
		return nil, fmt.Errorf("failed to scan on open ports on os: %w", err)
	}

	openPortsPerJob := map[JobWithTask][]PortWithName{}

	for jobID, nomadExposedPorts := range nomadExposedPortsPerJob {
		gss := genServices.GetByName(jobID)

		// Add ports expose on Nomad job specifically.
		for _, openPort := range nomadExposedPorts {
			key := JobWithTask{Name: jobID}
			openPortsPerJob[key] = append(openPortsPerJob[key], PortWithName{Port: openPort})
		}

		// Add ports from processes running raw on OS
		for _, gs := range gss {
			configuredPorts, err := ExtractPortsFromConfig(gs.ConfigFilePath)
			if err != nil {
				log.Printf("failed to extract ports from config file %q for job %s: %s", gs.ConfigFilePath, jobID, err)
				continue
			}

			for openPort, processName := range portsPerProcessName {
				if portName, ok := configuredPorts[openPort]; ok && strings.Contains(gs.Name, processName) {
					key := JobWithTask{Name: jobID, TaskName: processName}
					openPortsPerJob[key] = append(openPortsPerJob[key], PortWithName{
						Port: openPort,
						Name: portName,
					})
				}
			}
		}
	}

	return openPortsPerJob, nil
}
