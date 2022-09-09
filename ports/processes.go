package ports

import (
	"fmt"
	"log"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

const (
	vegaProcessName     = "vega"
	faucetProcessName   = "faucet"
	dataNodeProcessName = "data-node"
	walletProcessName   = "vegawallet"

	portStatusListen = "LISTEN"
	nomadProcessName = "nomad"
)

var networkProcessesNames = map[string]struct{}{
	vegaProcessName:     {},
	dataNodeProcessName: {},
	walletProcessName:   {},
}

// ScanNetworkProcessesPorts returns map of ports to process name for Capsule own processes.
func ScanNetworkProcessesPorts() (map[int64]string, error) {
	ps, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	out := map[int64]string{}

	for _, p := range ps {
		parent, err := p.Parent()
		if err != nil {
			continue
		}

		parentName, err := parent.Name()
		if err != nil {
			continue
		}

		if !strings.Contains(parentName, nomadProcessName) {
			continue
		}

		currentName, err := p.Name()
		if err != nil {
			continue
		}

		if _, ok := networkProcessesNames[currentName]; !ok {
			continue
		}

		cs, err := p.Connections()
		if err != nil {
			log.Printf("failed to get listen connections for process %s: %s", currentName, err)
			continue
		}

		outName := currentName
		if outName == vegaProcessName {
			cmds, err := p.CmdlineSlice()
			if err != nil {
				continue
			}

			// Check if command is faucet
			if len(cmds) != 0 && cmds[1] == faucetProcessName {
				outName = faucetProcessName
			}
		}

		for _, c := range cs {
			if c.Status == portStatusListen {
				out[int64(c.Laddr.Port)] = outName
			}
		}
	}

	return out, nil
}
