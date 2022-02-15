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
		return nil, fmt.Errorf("failed to persist config in output directory %s: %w", state.Config.OutputDir, err)
	}

	log.Println("generating network success")

	return generatedSvcs, nil
}

func start(ctx context.Context, state *state.NetworkState) error {
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

	res, err := runner.StartNetwork(ctx, state.Config, state.GeneratedServices)
	if err != nil {

		return fmt.Errorf("failed to start nomad network: %s", err)
	}
	state.RunningJobs = res

	log.Println("starting network success")
	return nil
}

func stop(ctx context.Context, state *state.NetworkState) {
	log.Println("stopping network")
	nomadRunner, err := nomad.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	runner := runner.New(nomadRunner)

	if err := runner.StopNetwork(ctx, state.RunningJobs); err != nil {
		log.Fatalf("failed to stop nomad network: %s", err)
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

	ctx := context.Background()
	arg := os.Args[1]
	switch arg {
	case "start":
		if err := start(ctx, networkState); err != nil {
			log.Fatal(err)
		}
		if err := networkState.Perist(); err != nil {
			log.Fatalf("Cannot save network state")
		}
	case "stop":
		stop(ctx, networkState)
	case "generate":
		if _, err := generate(networkState); err != nil {
			log.Fatal(err)
		}
		if err := networkState.Perist(); err != nil {
			log.Fatalf("Cannot save network state")
		}
	case "destroy":
		stop(ctx, networkState)
		cleanup(conf.OutputDir)
	default:
		log.Printf("unknown subcommand %s: expected 'start'|'stop'|'destroy' subcommands", arg)
		os.Exit(1)
	}

}
