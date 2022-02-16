package runner

import (
	"fmt"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/runner/nomad"
	"code.vegaprotocol.io/vegacapsule/types"
	"golang.org/x/sync/errgroup"

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
	vegaNetworkJobName string
}

func New(n *nomad.NomadRunner) *Runner {
	return &Runner{
		nomad:              n,
		vegaNetworkJobName: fmt.Sprintf("test-vega-network-%d", 1),
	}
}

func (r *Runner) RunDockerJob(conf config.DockerConfig) error {
	j := api.Job{
		ID:          strPoint(conf.Name),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: intPoint(0),
					Mode:     strPoint("fail"),
				},
				Name: strPoint(conf.Name),
				Tasks: []*api.Task{
					{
						Name:   conf.Name,
						Driver: "docker",
						Config: map[string]interface{}{
							"image":   conf.Image,
							"command": conf.Command,
							"args":    conf.Args,
						},
						Resources: &api.Resources{
							Networks: []*api.NetworkResource{
								{
									ReservedPorts: []api.Port{
										{
											Label: fmt.Sprintf("%s-port", conf.Name),
											Value: conf.StaticPort,
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
		return fmt.Errorf("failed to run nomad docker job: %w", err)
	}

	return nil
}

func (r *Runner) runNodeSets(conf *config.Config, nodeSets []types.NodeSet) error {
	jobs := []api.Job{}

	for i, ns := range nodeSets {
		tasks := make([]*api.Task, 0, 3)
		tasks = append(tasks,
			&api.Task{
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
					MemoryMB: intPoint(512),
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
					MemoryMB: intPoint(512),
				},
			})

		if ns.DataNode != nil {
			tasks = append(tasks, &api.Task{
				Name:   fmt.Sprintf("data-node-%s-%d", ns.Mode, i),
				Driver: "raw_exec",
				Config: map[string]interface{}{
					"command": ns.DataNode.BinaryPath,
					"args": []string{
						"node",
						"--home", ns.DataNode.HomeDir,
					},
				},
				Resources: &api.Resources{
					CPU:      intPoint(500),
					MemoryMB: intPoint(512),
				},
			})
		}

		jobs = append(jobs, api.Job{
			ID:          strPoint(fmt.Sprintf("nodeset-%s-%d", ns.Mode, i)),
			Datacenters: []string{"dc1"},
			TaskGroups: []*api.TaskGroup{
				{
					RestartPolicy: &api.RestartPolicy{
						Attempts: intPoint(0),
						Mode:     strPoint("fail"),
					},
					Name:  strPoint("vega"),
					Tasks: tasks,
				},
			},
		})
	}

	eg := new(errgroup.Group)
	for _, j := range jobs {
		j := j
		eg.Go(func() error {
			return r.nomad.RunAndWait(j)
		})
	}

	return eg.Wait()
}

func (r *Runner) runWallet(conf *config.WalletConfig, wallet *types.Wallet) error {
	j := &api.Job{
		ID:          strPoint("wallet-1"),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: intPoint(0),
					Mode:     strPoint("fail"),
				},
				Name: strPoint("vega"),
				Tasks: []*api.Task{
					{
						Name:   "wallet-1",
						Driver: "raw_exec",
						Config: map[string]interface{}{
							"command": conf.Binary,
							"args": []string{
								"service",
								"run",
								"--network", wallet.Network,
								"--no-version-check",
								"--output", "json",
								"--home", wallet.HomeDir,
							},
						},
						Resources: &api.Resources{
							CPU:      intPoint(500),
							MemoryMB: intPoint(512),
						},
					},
				},
			},
		},
	}

	return r.nomad.RunAndWait(*j)
}

func (r *Runner) StartRawNetwork(conf *config.Config, generatedSvcs *types.GeneratedServices) error {
	if generatedSvcs.Wallet != nil {
		if err := r.runWallet(conf.Network.Wallet, generatedSvcs.Wallet); err != nil {
			return fmt.Errorf("failed to run wallet: %w", err)
		}
	}

	if err := r.runNodeSets(conf, generatedSvcs.NodeSets); err != nil {
		return fmt.Errorf("failed to run node sets: %w", err)
	}

	return nil
}

func (r *Runner) StopRawNetwork() error {
	var gErr error
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
