package nomad

import (
	"context"
	"fmt"
	"path"
	"time"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"

	"github.com/hashicorp/nomad/api"
)

var defaultResourcesConfig = &api.Resources{
	CPU:      utils.ToPoint(500),
	MemoryMB: utils.ToPoint(512),
	DiskMB:   utils.ToPoint(550),
}

var defaultLogConfig = &api.LogConfig{
	MaxFileSizeMB: utils.ToPoint(500), // 500 Mb
}

func (r *JobRunner) defaultLogCollectorTask(jobName string) *api.Task {
	return &api.Task{
		Name:   "logger",
		Driver: "raw_exec",
		Config: map[string]interface{}{
			"command": r.capsuleBinary,
			"args": []string{
				"nomad", "logscollector",
				"--out-dir", path.Join(r.logsOutputDir, jobName),
			},
		},
		LogConfig: defaultLogConfig,
		Resources: defaultResourcesConfig,
	}
}

func (r *JobRunner) defaultNodeSetJob(ns types.NodeSet) *api.Job {
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
			LogConfig: defaultLogConfig,
			Resources: defaultResourcesConfig,
		},
		r.defaultLogCollectorTask(ns.Name),
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
			LogConfig: defaultLogConfig,
			Resources: defaultResourcesConfig,
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
						Leader: true,
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
						LogConfig: defaultLogConfig,
						Resources: defaultResourcesConfig,
					},
					r.defaultLogCollectorTask(wallet.Name),
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
						Leader: true,
						Config: map[string]interface{}{
							"command": binary,
							"args": []string{
								"faucet",
								"run",
								"--passphrase-file", fc.WalletPassFilePath,
								"--home", fc.HomeDir,
							},
						},
						LogConfig: defaultLogConfig,
						Resources: defaultResourcesConfig,
					},
					r.defaultLogCollectorTask(fc.Name),
				},
			},
		},
	}
}

func (r *JobRunner) defaultDockerJob(ctx context.Context, conf config.DockerConfig) *api.Job {
	ports := []api.Port{}
	portLabels := []string{}
	if conf.StaticPort != nil {
		ports = append(ports, api.Port{
			Label: fmt.Sprintf("%s-port", conf.Name),
			To:    conf.StaticPort.To,
			Value: conf.StaticPort.Value,
		})
		portLabels = append(portLabels, fmt.Sprintf("%s-port", conf.Name))
	}

	return &api.Job{
		ID:          &conf.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				Networks: []*api.NetworkResource{
					{
						ReservedPorts: ports,
					},
				},
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.ToPoint(0),
					Mode:     utils.ToPoint("fail"),
				},
				Name: &conf.Name,
				Tasks: []*api.Task{
					{
						Name:   conf.Name,
						Driver: "docker",
						Leader: true,
						Config: map[string]interface{}{
							"image":          conf.Image,
							"command":        conf.Command,
							"args":           conf.Args,
							"ports":          portLabels,
							"auth_soft_fail": conf.AuthSoftFail,
						},
						Env:       conf.Env,
						LogConfig: defaultLogConfig,
						Resources: &api.Resources{
							CPU:      utils.ToPoint(500),
							MemoryMB: utils.ToPoint(768),
						},
					},
					r.defaultLogCollectorTask(conf.Name),
				},
			},
		},
	}
}
