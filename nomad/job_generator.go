package nomad

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
)

func GenerateNomadNetworkJobs(conf *config.Config, generatedSvcs *types.GeneratedServices) ([]api.Job, error) {
	jobs := []api.Job{}
	for _, dc := range conf.Network.PreStart.Docker {
		jobs = append(jobs, *nomadDockerJob(dc))
	}

	if generatedSvcs.Faucet != nil {
		jobs = append(jobs, *nomadFaucetJob(*conf.VegaBinary, conf.Network.Faucet, generatedSvcs.Faucet))
	}

	if generatedSvcs.Wallet != nil {
		jobs = append(jobs, *nomadWalletJob(conf.Network.Wallet, generatedSvcs.Wallet))
	}

	nodeSetsJobs, err := nomadNodeSetJobs(generatedSvcs.NodeSets.ToSlice())
	if err != nil {
		return nil, fmt.Errorf("failed to generate node set jobs: %w", err)
	}
	jobs = append(jobs, nodeSetsJobs...)

	return jobs, nil
}

func PersistJobsTemplates(outputDir string, jobs []api.Job) error {
	if fileExists, _ := utils.FileExists(outputDir); !fileExists {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create missing output directory: %w", err)
		}
	}

	for _, job := range jobs {
		jobFile := filepath.Join(outputDir, fmt.Sprintf("%s.hcl", *job.ID))
		tplFile, err := utils.CreateFile(jobFile)
		if err != nil {
			return fmt.Errorf("failed to create template file: %w", err)
		}
		defer tplFile.Close()

		data, err := json.Marshal(job)
		if err != nil {
			return fmt.Errorf("failed to marshal nomad job: %w", err)
		}
		if _, err := tplFile.Write(data); err != nil {
			return fmt.Errorf("failed to write template file: %w", err)
		}
	}

	return nil
}

func defaultNodeSetJob(ns types.NodeSet) api.Job {
	tasks := make([]*api.Task, 0, 3)
	tasks = append(tasks,
		&api.Task{
			Name:   ns.Tendermint.Name,
			Driver: "raw_exec",
			Config: map[string]interface{}{
				"command": ns.Tendermint.BinaryPath,
				"args": []string{
					"tm",
					"node",
					"--home", ns.Tendermint.HomeDir,
				},
			},
			Resources: &api.Resources{
				CPU:      utils.IntPoint(500),
				MemoryMB: utils.IntPoint(512),
			},
		},
		&api.Task{
			Name:   ns.Vega.Name,
			Driver: "raw_exec",
			Config: map[string]interface{}{
				"command": ns.Vega.BinaryPath,
				"args": []string{
					"node",
					"--home", ns.Vega.HomeDir,
					"--nodewallet-passphrase-file", ns.Vega.NodeWalletPassFilePath,
				},
			},
			Resources: &api.Resources{
				CPU:      utils.IntPoint(500),
				MemoryMB: utils.IntPoint(512),
			},
		})

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
				CPU:      utils.IntPoint(500),
				MemoryMB: utils.IntPoint(512),
			},
		})
	}

	return api.Job{
		ID:          utils.StrPoint(ns.Name),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.IntPoint(0),
					Mode:     utils.StrPoint("fail"),
				},
				Name:  utils.StrPoint("vega"),
				Tasks: tasks,
			},
		},
	}
}

func nomadDockerJob(conf config.DockerConfig) *api.Job {
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
					Attempts: utils.IntPoint(0),
					Mode:     utils.StrPoint("fail"),
				},
				Name: &conf.Name,
				Tasks: []*api.Task{
					{
						Name:   conf.Name,
						Driver: "docker",
						Config: map[string]interface{}{
							"image":          conf.Image,
							"command":        conf.Command,
							"args":           conf.Args,
							"ports":          portLabels,
							"auth_soft_fail": conf.AuthSoftFail,
						},
						Env: conf.Env,
						Resources: &api.Resources{
							CPU:      utils.IntPoint(500),
							MemoryMB: utils.IntPoint(768),
						},
					},
				},
			},
		},
	}
}

func nomadNodeSetJobs(nodeSets []types.NodeSet) ([]api.Job, error) {
	jobs := make([]api.Job, 0, len(nodeSets))

	for _, ns := range nodeSets {
		if ns.NomadJobRaw == nil {
			log.Printf("adding node set %q with default Nomad job definition", ns.Name)

			jobs = append(jobs, defaultNodeSetJob(ns))
			continue
		}

		job, err := jobspec.Parse(strings.NewReader(*ns.NomadJobRaw))
		if err != nil {
			return nil, err
		}

		log.Printf("adding node set %q with custom Nomad job definition", ns.Name)

		jobs = append(jobs, *job)
	}

	return jobs, nil
}

func nomadWalletJob(conf *config.WalletConfig, wallet *types.Wallet) *api.Job {
	return &api.Job{
		ID:          &wallet.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.IntPoint(0),
					Mode:     utils.StrPoint("fail"),
				},
				Name: utils.StrPoint("vega"),
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
							CPU:      utils.IntPoint(500),
							MemoryMB: utils.IntPoint(512),
						},
					},
				},
			},
		},
	}
}

func nomadFaucetJob(binary string, conf *config.FaucetConfig, fc *types.Faucet) *api.Job {
	return &api.Job{
		ID:          &fc.Name,
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			{
				RestartPolicy: &api.RestartPolicy{
					Attempts: utils.IntPoint(0),
					Mode:     utils.StrPoint("fail"),
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
							CPU:      utils.IntPoint(500),
							MemoryMB: utils.IntPoint(512),
						},
					},
				},
			},
		},
	}
}
