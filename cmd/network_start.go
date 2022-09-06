package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"code.vegaprotocol.io/vegacapsule/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

var (
	doNotStopOnFailure bool
)

var netStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts existing network",
	RunE: func(cmd *cobra.Command, args []string) error {
		netState, err := state.LoadNetworkState(homePath)
		if err != nil {
			return err
		}

		if netState.Empty() {
			return networkNotBootstrappedErr("start")
		}

		updatedNetState, err := netStart(context.Background(), *netState)
		if err != nil {
			return fmt.Errorf("failed to start network: %w", err)
		}

		return updatedNetState.Persist()
	},
}

func init() {
	netStartCmd.PersistentFlags().BoolVar(&doNotStopOnFailure,
		"do-not-stop-on-failure",
		false,
		"Do not stop partially running network when failed to start",
	)
}

func netStart(ctx context.Context, state state.NetworkState) (*state.NetworkState, error) {
	log.Println("starting network")

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner, err := nomad.NewJobRunner(nomadClient, *state.Config.VegaCapsuleBinary, state.Config.LogsDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create job runner: %w", err)
	}

	res, err := nomadRunner.StartNetwork(ctx, state.Config, state.GeneratedServices, !doNotStopOnFailure)
	if err != nil {
		return nil, fmt.Errorf("failed to start nomad network: %s", err)
	}
	state.RunningJobs = res

	log.Println("starting network success")

	printPorts(state)

	return &state, nil
}

func printPorts(state state.NetworkState) {
	log.Println("collecting exposed nodes addresses")

	time.Sleep(time.Second * 5)
	openPorts := getRunningPorts()

	if fOpenPorts, ok := openPorts["faucet"]; ok && state.GeneratedServices.Faucet != nil {
		configuredPorts, err := ExtractPortsFromTOML(state.GeneratedServices.Faucet.ConfigFilePath)
		if err == nil {
			for _, port := range fOpenPorts {
				if len(fOpenPorts) != 0 {
					fmt.Printf("%q:\n", state.GeneratedServices.Faucet.Name)
				}
				if portName, ok := configuredPorts[port]; ok {
					fmt.Printf("- %s localhost:%d\n", portName, port)
				}
			}
		}
	}

	if wOpenPorts, ok := openPorts["wallet"]; ok && state.GeneratedServices.Wallet != nil {
		configuredPorts, err := ExtractPortsFromTOML(state.GeneratedServices.Wallet.ConfigFilePath)
		if err == nil {
			for _, port := range wOpenPorts {
				if len(wOpenPorts) != 0 {
					fmt.Printf("\n%q:\n", state.GeneratedServices.Wallet.Name)
				}
				if portName, ok := configuredPorts[port]; ok {
					fmt.Printf("- %s: localhost:%d\n", portName, port)
				}
			}
		}
	}

	for _, ns := range state.GeneratedServices.NodeSets {
		if vOpenPorts, ok := openPorts["vega"]; ok {
			configuredPorts, err := ExtractPortsFromTOML(ns.Vega.ConfigFilePath)
			if err == nil {
				if len(vOpenPorts) != 0 {
					fmt.Printf("\n%q Vega:\n", ns.Name)
				}
				for _, port := range vOpenPorts {
					if portName, ok := configuredPorts[port]; ok {
						fmt.Printf("- %s: localhost:%d\n", portName, port)
					}
				}
			}
		}

		if dnOpenPorts, ok := openPorts["data-node"]; ok && ns.DataNode != nil {
			configuredPorts, err := ExtractPortsFromTOML(ns.DataNode.ConfigFilePath)
			if err == nil {
				if len(dnOpenPorts) != 0 {
					fmt.Printf("\n%q Data Node:\n", ns.Name)
				}
				for _, port := range dnOpenPorts {
					if portName, ok := configuredPorts[port]; ok {
						fmt.Printf("- %s: localhost:%d\n", portName, port)
					}
				}
			}

		}
	}
}

func getRunningPorts() map[string][]int64 {
	ps, err := process.Processes()
	if err != nil {
		panic(err)
	}

	out := map[string][]int64{}

	for _, p := range ps {
		parent, err := p.Parent()
		if err != nil {
			continue
		}

		parentName, err := parent.Name()
		if err != nil {
			continue
		}

		if parentName != "nomad_1.3.1" {
			continue
		}

		n, err := p.Name()
		if err != nil {
			continue
		}

		switch n {
		case "vega", "data-node", "vegawallet":
			cs, err := p.Connections()
			if err != nil {
				fmt.Println(err)
				continue
			}

			cmds, err := p.CmdlineSlice()
			if err != nil {
				continue
			}

			outName := n
			if n == "vega" {
				if cmds[1] == "faucet" {
					outName = "faucet"
				}
			}

			if _, ok := out[outName]; !ok {
				out[outName] = []int64{}
			}

			for _, c := range cs {
				if c.Status == "LISTEN" {
					out[outName] = append(out[outName], int64(c.Laddr.Port))
				}
			}
		}
	}

	return out
}
