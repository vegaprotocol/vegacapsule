package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/runner"
	"code.vegaprotocol.io/vegacapsule/runner/nomad"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/utils"
)

func generate(state state.NetworkState, force bool) (*state.NetworkState, error) {
	if force {
		if err := os.RemoveAll(state.Config.OutputDir); err != nil {
			return nil, fmt.Errorf("failed to remove output folder with --force flag: %w", err)
		}
	} else if state.GeneratedServices != nil {
		return nil, fmt.Errorf("failed to generate network: network is already generated")
	}

	if netDirExists, _ := utils.FileExists(state.Config.OutputDir); netDirExists {
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

	state.GeneratedServices = generatedSvcs
	return &state, nil
}

func start(ctx context.Context, state state.NetworkState) (*state.NetworkState, error) {
	log.Println("starting network")
	if state.Empty() {
		return nil, fmt.Errorf("failed to start network: network is not bootstrapped")
	}

	nomadRunner, err := nomad.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start the network: %w", err)
	}

	runner := runner.New(nomadRunner)

	res, err := runner.StartNetwork(ctx, state.Config, state.GeneratedServices)
	if err != nil {
		return nil, fmt.Errorf("failed to start nomad network: %s", err)
	}
	state.RunningJobs = res

	log.Println("starting network success")
	return &state, nil
}

func bootstrap(ctx context.Context, state state.NetworkState, force bool) (*state.NetworkState, error) {
	log.Println("starting network")

	updatedState, err := generate(state, force)
	if err != nil {
		return nil, fmt.Errorf("failed to generate config for network: %w", err)
	}

	updatedState, err = start(ctx, *updatedState)
	if err != nil {
		return nil, fmt.Errorf("failed to start network: %w", err)
	}

	return updatedState, nil
}

func stop(ctx context.Context, state *state.NetworkState) {
	log.Println("stopping network")
	if state.Empty() {
		log.Fatalf("cannot stop network: network is not bootstrapped")
	}

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
		fmt.Println("expected 'generate'|'bootstrap'|'start'|'stop'|'destroy' subcommands")
		os.Exit(1)
	}

	subcommand := os.Args[1]
	command := flag.NewFlagSet("main", flag.ExitOnError)
	configFilePath := command.String("config-path", "", "enable")
	force := command.Bool("force", false, "enable")
	networkHomePath := command.String("home-path", defaultNetworkHome(), "enable")

	if err := command.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	switch subcommand {
	case "bootstrap", "generate":
		if *configFilePath == "" {
			log.Fatalf("Missing config file path. Use the `-config-path` flag")
		}

		conf, err := config.ParseConfigFile(*configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		networkState, err := state.LoadNetworkState(conf.OutputDir)
		if err != nil {
			log.Fatalf("cannot load network state: %s", err)
		}
		networkState.Config = conf

		if subcommand == "bootstrap" {
			networkState, err = bootstrap(ctx, *networkState, *force)
		} else {
			networkState, err = generate(*networkState, *force)
		}

		if err != nil {
			log.Fatal(err)
		}
		if err := networkState.Perist(); err != nil {
			log.Fatalf("Cannot save network state")
		}
	case "start", "stop", "destroy":
		if *configFilePath != "" {
			log.Printf("Flag `-config-path` is ignored for %s subcommand. Use the `-home-path` flag.", subcommand)
		}
		log.Printf("Using network network home: %s", *networkHomePath)

		networkState, err := state.LoadNetworkState(*networkHomePath)
		if err != nil {
			log.Fatalf("failed to %s network: %s", subcommand, err)
		}

		if networkState.Empty() {
			log.Fatalf("Failed to %s network: network not bootstrapped. Use the 'bootstrap' subcommand or provide different network home with the `-home-path` flag", subcommand)
		}

		if subcommand == "start" {
			networkState, err := start(ctx, *networkState)
			if err != nil {
				log.Fatalf("failed to start network: %s", err)
			}
			if err := networkState.Perist(); err != nil {
				log.Fatalf("cannot persist network state: %s", err)
			}
		} else if subcommand == "stop" {
			stop(ctx, networkState)
		} else {
			stop(ctx, networkState)
			cleanup(*networkHomePath)
		}
	default:
		log.Printf("unknown subcommand %s: expected 'generate'|'bootstrap'|'start'|'stop'|'destroy' subcommands", subcommand)
		os.Exit(1)
	}
}

// TODO: This will be replaced during CLI refactor
func defaultNetworkHome() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s/%s", user.HomeDir, "vega/vegacapsule/testnet")
}
