package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/runner/nomad"
	"code.vegaprotocol.io/vegacapsule/types"

	"github.com/hashicorp/nomad/api"
)

var (
	dockerGanacheImage    = "trufflesuite/ganache-cli:v6.12.2"
	dockerVegaImage       = "ghcr.io/vegaprotocol/vega/vega:$vega_version"
	dockerVegawalletImage = "vegaprotocol/vegawallet:$vegawallet_version"
	dockerDatanodeImage   = "ghcr.io/vegaprotocol/data-node/data-node:$datanode_version"
	dockerVegatoolsImage  = "vegaprotocol/vegatools:$vegatools_version"
	dockerClefImage       = "ghcr.io/vegaprotocol/devops-infra/clef:latest"
)

type Runner struct {
	nomad              *nomad.NomadRunner
	ganacheJobName     string
	vegaNetworkJobName string
}

func New(n *nomad.NomadRunner) *Runner {
	return &Runner{
		nomad:              n,
		ganacheJobName:     fmt.Sprintf("test-vega-ganache-%d", 1),
		vegaNetworkJobName: fmt.Sprintf("test-vega-network-%d", 1),
	}
}

func ganacheCheck(url string, timeout time.Duration) error {
	for start := time.Now(); time.Since(start) < timeout; {
		time.Sleep(2 * time.Second)
		postBody, _ := json.Marshal(map[string]string{
			"method": "web3_clientVersion",
		})
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(url, "application/json", responseBody)
		if err != nil {
			log.Println("ganache not yet ready", err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println("body:", string(body))
		if err != nil {
			log.Println(err)
			continue
		}

		if strings.Contains(string(body), "EthereumJS") {
			log.Println("ganache is ready")
			return nil
		}
		continue
	}

	return fmt.Errorf("ganache container has timed out")
}

func (r *Runner) startGanache() error {
	j := &api.Job{
		ID:          strPoint(r.ganacheJobName),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: intPoint(0),
					Mode:     strPoint("fail"),
				},
				Name: strPoint("ganache"),
				Tasks: []*api.Task{
					{
						Name:   "ganache-1",
						Driver: "docker",
						Config: map[string]interface{}{
							"image":   "ghcr.io/vegaprotocol/devops-infra/ganache:latest",
							"command": "ganache-cli",
							"args": []string{
								"--blockTime", "1",
								"--chainId", "1440",
								"--networkId", "1441",
								"-h", "0.0.0.0",
								"-p", "8545",
								"-m", "cherry manage trip absorb logic half number test shed logic purpose rifle",
								"--db", "/app/ganache-db",
							},
						},
						Resources: &api.Resources{
							Networks: []*api.NetworkResource{
								{
									ReservedPorts: []api.Port{
										{
											Label: "ganache-port",
											Value: 8545,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if err := r.nomad.RunAndWait(j); err != nil {
		return fmt.Errorf("failed to run nomad ganache job: %w", err)
	}

	// if err := ganacheCheck(url, time.Minute*3); err != nil {
	// 	return fmt.Errorf("failed to start ganache container: %w", err)
	// }

	return nil
}

func (r *Runner) RunGanache() error {
	return r.startGanache()
}

func (r *Runner) StartRawNetwork(conf *config.Config, nodeSets []types.NodeSet) error {
	// TODO this should come from config - do the pre-start network configuration
	if err := r.startGanache(); err != nil {
		return err
	}

	tasks := make([]*api.Task, 0, len(nodeSets))

	for i, ns := range nodeSets {
		tasks = append(tasks, &api.Task{
			Name:   fmt.Sprintf("tendermint-%d", i),
			Driver: "raw_exec",
			Config: map[string]interface{}{
				"command": conf.VegaBinary,
				"args": []string{
					"tm",
					"node",
					"--home", ns.Tendermint.HomeDir,
				},
			},
			Resources: &api.Resources{
				CPU:      intPoint(500),
				MemoryMB: intPoint(256),
			},
		},
			&api.Task{
				Name:   fmt.Sprintf("vega-%s-%d", ns.Mode, i),
				Driver: "raw_exec",
				Config: map[string]interface{}{
					"command": conf.VegaBinary,
					"args": []string{
						"node",
						"--home", ns.Vega.HomeDir,
						"--nodewallet-passphrase-file", ns.Vega.NodeWalletPassFilePath,
					},
				},
				Resources: &api.Resources{
					CPU:      intPoint(500),
					MemoryMB: intPoint(256),
				},
			})
	}

	j := &api.Job{
		ID:          strPoint(r.vegaNetworkJobName),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				Name:  strPoint("vega"),
				Tasks: tasks,
			},
		},
	}

	if err := r.nomad.RunAndWait(j); err != nil {
		return fmt.Errorf("failed to run nomad network: %w", err)
	}

	return nil
}

func (r *Runner) StopRawNetwork() error {
	var gErr error
	if _, err := r.nomad.Stop(r.ganacheJobName, true); err != nil {
		gErr = fmt.Errorf("failed to stop %q job: %w", r.ganacheJobName, err)
	}
	if _, err := r.nomad.Stop(r.vegaNetworkJobName, true); err != nil {
		gErr = fmt.Errorf("failed to stop %q job: %w: %s", r.vegaNetworkJobName, err, gErr)
	}

	return gErr
}

func strPoint(s string) *string {
	return &s
}

func intPoint(i int) *int {
	return &i
}
