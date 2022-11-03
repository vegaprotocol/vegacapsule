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

	dataNodeCmdName = "datanode"
	walletCmdName   = "wallet"

	portStatusListen = "LISTEN"
	nomadProcessName = "nomad"
	visorProcessName = "visor"
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

		if !(strings.Contains(parentName, nomadProcessName) || strings.Contains(parentName, visorProcessName)) {
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

			if len(cmds) == 0 {
				continue
			}

			switch cmds[1] {
			case faucetProcessName:
				outName = faucetProcessName
			case dataNodeCmdName:
				outName = dataNodeProcessName
			case walletCmdName:
				outName = walletCmdName
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
