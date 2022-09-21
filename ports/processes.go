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
	delveName        = "dlv"
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

		outName := currentName

		switch outName {
		case delveName:
			pc, outNameCh, err := findChild(p)
			if err != nil {
				continue
			}

			if err = addPort(p, outNameCh, out); err != nil {
				continue
			}
			p, outName = pc, outNameCh
			fallthrough
		case vegaProcessName:
			cmds, err := p.CmdlineSlice()
			if err != nil {
				continue
			}

			// Check if command is faucet
			if len(cmds) != 0 && cmds[1] == faucetProcessName {
				outName = faucetProcessName
			}
		default:
			if _, ok := networkProcessesNames[currentName]; !ok {
				continue
			}
		}

		if err = addPort(p, outName, out); err != nil {
			continue
		}
	}

	return out, nil
}

func addPort(p *process.Process, outName string, out map[int64]string) error {
	cs, err := p.Connections()
	if err != nil {
		log.Printf("failed to get listen connections for process %s: %s", outName, err)
		return err
	}

	for _, c := range cs {
		if c.Status == portStatusListen {
			out[int64(c.Laddr.Port)] = outName
		}
	}
	return nil
}

// recursive find child that matches process name in networkProcessesNames
func findChild(p *process.Process) (*process.Process, string, error) {
	children, err := p.Children()
	if err != nil {
		return nil, "", err
	}

	for _, child := range children {
		childName, err := child.Name()
		if err != nil {
			continue
		}

		if _, ok := networkProcessesNames[childName]; ok {
			return child, childName, nil
		}

		return findChild(child)
	}

	return nil, "", fmt.Errorf("process not recognised")
}
