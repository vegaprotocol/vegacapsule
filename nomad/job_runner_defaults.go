package nomad

import (
	"time"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	"github.com/hashicorp/nomad/api"
)

func (r *JobRunner) defaultNodeSetJob(ns types.NodeSet, capsuleBinary string) *api.Job {
	tasks := make([]*api.Task, 0, 3)
	tasks = append(tasks,
		&api.Task{
			Leader: true,
			Name:   ns.Vega.Name,
			Driver: "raw_exec",
			Config: map[string]interface{}{
				"command": ns.Vega.BinaryPath,
				"args": []string{
					"node",
					"--home", ns.Vega.HomeDir,
					"--tendermint-home", ns.Tendermint.HomeDir,
					"--nodewallet-passphrase-file", ns.Vega.NodeWalletPassFilePath,
				},
			},
			Resources: &api.Resources{
				CPU:      utils.ToPoint(500),
				MemoryMB: utils.ToPoint(512),
			},
		},
		&api.Task{
			Name:   "logger",
			Driver: "raw_exec",
			Config: map[string]interface{}{
				"command": capsuleBinary,
				"args": []string{
					"nomad", "log",
					"--job-id", ns.Name,
				},
			},
			Resources: &api.Resources{
				CPU:      utils.ToPoint(500),
				MemoryMB: utils.ToPoint(512),
			},
		},
	)

	if ns.DataNode != nil {
		tasks = append(tasks, &api.Task{
			Name:   ns.DataNode.Name,
			Driver: "raw_exec",
			Config: map[string]interface{}{
				"command": ns.DataNode.BinaryPath,
				"args": []string{
					"node",
					"--home", ns.DataNode.HomeDir,
				},
			},
			Resources: &api.Resources{
				CPU:      utils.ToPoint(500),
				MemoryMB: utils.ToPoint(512),
			},
		})
	}

	return &api.Job{
		ID:          utils.ToPoint(ns.Name),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.ToPoint(0),
					Interval: utils.ToPoint(time.Second * 5),
					Mode:     utils.ToPoint("fail"),
				},
				Name:  utils.ToPoint("vega"),
				Tasks: tasks,
			},
		},
	}
}

func (r *JobRunner) defaultWalletJob(conf *config.WalletConfig, wallet *types.Wallet) *api.Job {
	return &api.Job{
		ID:          &wallet.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.ToPoint(0),
					Mode:     utils.ToPoint("fail"),
				},
				Name: utils.ToPoint("vega"),
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
								"--automatic-consent",
								"--no-version-check",
								"--output", "json",
								"--home", wallet.HomeDir,
							},
						},
						Resources: &api.Resources{
							CPU:      utils.ToPoint(500),
							MemoryMB: utils.ToPoint(512),
						},
					},
				},
			},
		},
	}
}

func (r *JobRunner) defaultFaucetJob(binary string, conf *config.FaucetConfig, fc *types.Faucet) *api.Job {
	return &api.Job{
		ID:          &fc.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.ToPoint(0),
					Mode:     utils.ToPoint("fail"),
				},
				Name: &conf.Name,
				Tasks: []*api.Task{
					{
						Name:   conf.Name,
						Driver: "raw_exec",
						Config: map[string]interface{}{
							"command": binary,
							"args": []string{
								"faucet",
								"run",
								"--passphrase-file", fc.WalletPassFilePath,
								"--home", fc.HomeDir,
							},
						},
						Resources: &api.Resources{
							CPU:      utils.ToPoint(500),
							MemoryMB: utils.ToPoint(512),
						},
					},
				},
			},
		},
	}
}
