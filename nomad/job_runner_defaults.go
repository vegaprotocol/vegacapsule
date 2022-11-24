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

var (
	defaultLogConfig = &api.LogConfig{
		MaxFileSizeMB: utils.ToPoint(50), // 500 Mb
	}
	defaultResourcesConfig = &api.Resources{
		CPU:      utils.ToPoint(500),
		MemoryMB: utils.ToPoint(512),
		DiskMB:   utils.ToPoint(550),
	}
	defaultRestartPolicy = &api.RestartPolicy{
		Attempts: utils.ToPoint(0),
		Interval: utils.ToPoint(time.Second * 5),
		Delay:    utils.ToPoint(time.Second * 1),
		Mode:     utils.ToPoint("fail"),
	}
	defaultReschedulePolicy = &api.ReschedulePolicy{
		Attempts:  utils.ToPoint(0),
		Unlimited: utils.ToPoint(false),
	}
	defaultKillTimeout = utils.ToPoint(time.Duration(time.Second * 20))
)

func mergeResourcesWithDefault(customRes *config.Resources) *api.Resources {
	result := *defaultResourcesConfig

	if customRes == nil {
		return &result
	}

	if customRes.CPU != nil {
		result.CPU = customRes.CPU
	}

	if customRes.Cores != nil {
		result.Cores = customRes.Cores
	}

	if customRes.DiskMB != nil {
		result.DiskMB = customRes.DiskMB
	}

	if customRes.MemoryMB != nil {
		result.MemoryMB = customRes.MemoryMB
	}

	if customRes.MemoryMaxMB != nil {
		result.MemoryMaxMB = customRes.MemoryMaxMB
	}

	return &result
}

func hasLogsCollectorTask(job *api.Job) bool {
	for _, tg := range job.TaskGroups {
		for _, task := range tg.Tasks {
			if task.Name == types.NomadLogsCollectorTaskName {
				return true
			}
		}
	}

	return false
}

func (r *JobRunner) defaultLogCollectorTask(jobName string) *api.Task {
	return &api.Task{
		Name:   types.NomadLogsCollectorTaskName,
		Driver: "raw_exec",
		Config: map[string]interface{}{
			"command": r.capsuleBinary,
			"args": []string{
				"nomad", "logscollector",
				"--out-dir", path.Join(r.logsOutputDir, jobName),
			},
		},
		LogConfig:   defaultLogConfig,
		Resources:   defaultResourcesConfig,
		KillTimeout: defaultKillTimeout,
	}
}

func (r *JobRunner) defaultNodeSetTasks(ns types.NodeSet) []*api.Task {
	if ns.Visor != nil {
		return []*api.Task{
			{
				Name:   ns.Visor.Name,
				Driver: "raw_exec",
				Config: map[string]interface{}{
					"command": ns.Visor.BinaryPath,
					"args": []string{
						"run",
						"--home", ns.Visor.HomeDir,
					},
				},
				LogConfig: defaultLogConfig,
				Resources: &api.Resources{
					CPU:      utils.ToPoint(1000),
					MemoryMB: utils.ToPoint(512),
				},
				KillTimeout: defaultKillTimeout,
			},
			r.defaultLogCollectorTask(ns.Name),
		}
	}

	tasks := make([]*api.Task, 0, 2)

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
			LogConfig:   defaultLogConfig,
			Resources:   defaultResourcesConfig,
			KillTimeout: defaultKillTimeout,
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
					config.DataNodeSubCmd,
					"node",
					"--home", ns.DataNode.HomeDir,
				},
			},
			LogConfig:   defaultLogConfig,
			Resources:   defaultResourcesConfig,
			KillTimeout: defaultKillTimeout,
		})
	}

	return tasks
}

func (r *JobRunner) defaultNodeSetJob(ns types.NodeSet) *api.Job {
	return &api.Job{
		ID:          utils.ToPoint(ns.Name),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				EphemeralDisk: &api.EphemeralDisk{
					SizeMB: utils.ToPoint(550),
				},
				Name:             utils.ToPoint("vega"),
				RestartPolicy:    defaultRestartPolicy,
				ReschedulePolicy: defaultReschedulePolicy,
				Tasks:            r.defaultNodeSetTasks(ns),
			},
		},
	}
}

func (r *JobRunner) defaultWalletJob(wallet *types.Wallet) *api.Job {
	return &api.Job{
		ID:          &wallet.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				EphemeralDisk: &api.EphemeralDisk{
					SizeMB: utils.ToPoint(550),
				},
				Name:             utils.ToPoint("vega"),
				RestartPolicy:    defaultRestartPolicy,
				ReschedulePolicy: defaultReschedulePolicy,
				Tasks: []*api.Task{
					{
						Name:   "wallet-1",
						Driver: "raw_exec",
						Leader: true,
						Config: map[string]interface{}{
							"command": wallet.BinaryPath,
							"args": []string{
								config.WalletSubCmd,
								"service",
								"run",
								"--network", wallet.Network,
								"--automatic-consent",
								"--no-version-check",
								"--output", "json",
								"--home", wallet.HomeDir,
							},
						},
						LogConfig:   defaultLogConfig,
						Resources:   defaultResourcesConfig,
						KillTimeout: defaultKillTimeout,
					},
					r.defaultLogCollectorTask(wallet.Name),
				},
			},
		},
	}
}

func (r *JobRunner) defaultFaucetJob(conf *config.FaucetConfig, fc *types.Faucet) *api.Job {
	return &api.Job{
		ID:          &fc.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				EphemeralDisk: &api.EphemeralDisk{
					SizeMB: utils.ToPoint(550),
				},
				Name:             &conf.Name,
				RestartPolicy:    defaultRestartPolicy,
				ReschedulePolicy: defaultReschedulePolicy,
				Tasks: []*api.Task{
					{
						Name:   conf.Name,
						Driver: "raw_exec",
						Leader: true,
						Config: map[string]interface{}{
							"command": fc.BinaryPath,
							"args": []string{
								config.FaucetSubCmd,
								"run",
								"--passphrase-file", fc.WalletPassFilePath,
								"--home", fc.HomeDir,
							},
						},
						LogConfig:   defaultLogConfig,
						Resources:   defaultResourcesConfig,
						KillTimeout: defaultKillTimeout,
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
				EphemeralDisk: &api.EphemeralDisk{
					SizeMB: utils.ToPoint(550),
				},
				Name:             &conf.Name,
				RestartPolicy:    defaultRestartPolicy,
				ReschedulePolicy: defaultReschedulePolicy,
				Networks: []*api.NetworkResource{
					{
						ReservedPorts: ports,
					},
				},
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
							"volumes":        conf.VolumeMounts,
						},
						Env:         conf.Env,
						LogConfig:   defaultLogConfig,
						Resources:   mergeResourcesWithDefault(conf.Resources),
						KillTimeout: defaultKillTimeout,
					},
					r.defaultLogCollectorTask(conf.Name),
				},
			},
		},
	}
}
