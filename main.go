package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/runner"
	"code.vegaprotocol.io/vegacapsule/runner/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
)

func generate(state *state.NetworkState) (*types.GeneratedServices, error) {
	if state.GeneratedServices != nil {
		log.Println("Network already generated. Generate skipped")
		return state.GeneratedServices, nil
	}

	if _, err := os.Stat(state.Config.OutputDir); os.IsExist(err) {
		return nil, fmt.Errorf("output directory %q already exist", state.Config.OutputDir)
	}

	log.Println("generating network")

	gen, err := generator.New(state.Config)
	if err != nil {
		return nil, err
	}

	generatedSvcs, err := gen.Generate()
	if err != nil {
		return nil, err
	}

	if err := state.Config.Persist(); err != nil {
		return nil, fmt.Errorf("failed to persist config in output directory %s", err)
	}

	log.Println("generating network success")

	return generatedSvcs, nil
}

func start(state *state.NetworkState) error {
	log.Println("starting network")

	generatedSvs, err := generate(state)
	if err != nil {
		return fmt.Errorf("failed to generate config for network: %w", err)
	}
	state.GeneratedServices = generatedSvs

	nomadRunner, err := nomad.New(nil)
	if err != nil {
		return err
	}

	runner := runner.New(nomadRunner)

	ctx := context.Background()

	for _, dc := range state.Config.Network.PreStart.Docker {
		if err := runner.RunDockerJob(ctx, dc); err != nil {
			return fmt.Errorf("failed to run pre start job %s: %w", dc.Name, err)
		}
	}
	state.PreTasks = jobNames(state.Config.Network.PreStart)

	if err := runner.StartRawNetwork(ctx, state.Config, state.GeneratedServices); err != nil {
		return fmt.Errorf("failed to start nomad network: %s", err)
	}

	log.Println("starting network success")
	return nil
}

func stop(state *state.NetworkState) {
	log.Println("stopping network")
	nomadRunner, err := nomad.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	runner := runner.New(nomadRunner)

	if err := runner.StopRawNetwork(state.GeneratedServices); err != nil {
		log.Fatalf("failed to stop nomad network: %s", err)
	}

	if err := runner.StopJobs(state.PreTasks); err != nil {
		log.Fatalf("failed to stop per-tasks: %s", err)
	}

	log.Println("stopping network success")
}

func cleanup(outputDir string) {
	log.Println("network cleaning up")

	if err := os.RemoveAll(outputDir); err != nil {
		log.Fatalf("failed cleanup network: %s", err)
	}

	log.Println("network cleaning up success")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'start'|'stop'|'destroy' subcommands")
		os.Exit(1)
	}

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	configFilePath := startCmd.String("config-path", "", "enable")

	if err := startCmd.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	conf, err := config.ParseConfigFile(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	networkState, err := state.LoadNetworkState(conf.OutputDir)
	if err != nil {
		log.Fatalf("cannot load network state: %s", err)
	}
	if networkState.Config == nil {
		networkState.Config = conf
	}

	arg := os.Args[1]
	switch arg {
	case "start":
		if err := start(networkState); err != nil {
			log.Fatal(err)
		}
		if err := networkState.Perist(); err != nil {
			log.Fatalf("Cannot save network state")
		}
	case "stop":
		stop(networkState)
	case "generate":
		if _, err := generate(networkState); err != nil {
			log.Fatal(err)
		}
		if err := networkState.Perist(); err != nil {
			log.Fatalf("Cannot save network state")
		}
	case "destroy":
		stop(networkState)
		cleanup(conf.OutputDir)
	default:
		log.Printf("unknown subcommand %s: expected 'start'|'stop'|'destroy' subcommands", arg)
		os.Exit(1)
	}

}

func jobNames(jobs *config.PrestartConfig) []string {
	if jobs == nil {
		return []string{}
	}

	jobNames := make([]string, len(jobs.Docker))

	for jobIdx, jobDetails := range jobs.Docker {
		jobNames[jobIdx] = jobDetails.Name
	}

	return jobNames
}
