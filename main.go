package main

import (
	"context"
<<<<<<< HEAD
	"encoding/json"
=======
	"errors"
>>>>>>> create default config && add simple add-validator func skeleton
	"flag"
	"fmt"
	"log"
	"os"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/nomad"
	nmrunner "code.vegaprotocol.io/vegacapsule/nomad/runner"
	"code.vegaprotocol.io/vegacapsule/state"
	"code.vegaprotocol.io/vegacapsule/types"
	"code.vegaprotocol.io/vegacapsule/utils"
)

func generate(state state.NetworkState, force bool) (*state.NetworkState, error) {
	if force {
		if err := os.RemoveAll(*state.Config.OutputDir); err != nil {
			return nil, fmt.Errorf("failed to remove output folder with --force flag: %w", err)
		}
	} else if state.GeneratedServices != nil {
		return nil, fmt.Errorf("failed to generate network: network is already generated")
	}

	if netDirExists, _ := utils.FileExists(*state.Config.OutputDir); netDirExists {
		return nil, fmt.Errorf("output directory %q already exist", *state.Config.OutputDir)
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
		return nil, fmt.Errorf("failed to persist config in output directory %s: %w", *state.Config.OutputDir, err)
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

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner := nomad.NewJobRunner(nomadClient)

	res, err := nomadRunner.StartNetwork(ctx, state.Config, state.GeneratedServices)
	if err != nil {
		return nil, fmt.Errorf("failed to start nomad network: %s", err)
	}
	state.RunningJobs = res

	log.Println("starting network success")
	return &state, nil
}

func bootstrap(ctx context.Context, state state.NetworkState, force bool) (*state.NetworkState, error) {
	log.Println("bootstraping network")

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

func stop(ctx context.Context, state *state.NetworkState) error {
	log.Println("stopping network")
	if state.Empty() {
		log.Fatalf("cannot stop network: network is not bootstrapped")
	}

	nomadClient, err := nomad.NewClient(nil)
	if err != nil {
		return fmt.Errorf("failed to create nomad client: %w", err)
	}

	nomadRunner := nomad.NewJobRunner(nomadClient)

	if err := nomadRunner.StopNetwork(ctx, state.RunningJobs); err != nil {
		return fmt.Errorf("failed to stop nomad network: %w", err)
	}

	log.Println("stopping network success")
	return nil
}

func cleanup(outputDir string) {
	log.Println("network cleaning up")

	if err := os.RemoveAll(outputDir); err != nil {
		log.Fatalf("failed cleanup network: %s", err)
	}

	log.Println("network cleaning up success")
}

func addValidator(state state.NetworkState) error {
	gen, err := generator.New(state.Config)
	if err != nil {
		return err
	}

	for _, ns := range state.GeneratedServices.NodeSets {
		if ns.Mode == types.NodeModeValidator {
			return gen.AddValidator(len(state.GeneratedServices.NodeSets), ns)
		}
	}

	return errors.New("no validator found")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'generate'|'bootstrap'|'start'|'stop'|'destroy' subcommands")
		os.Exit(1)
	}

	homePath, err := config.DefaultNetworkHome()
	if err != nil {
		log.Fatalf("Failed to get default network home: %s", err.Error())
	}

	subcommand := os.Args[1]
	command := flag.NewFlagSet("main", flag.ExitOnError)
	configFilePath := command.String("config-path", "", "enable")
	force := command.Bool("force", false, "enable")
	networkHomePath := command.String("home-path", homePath, "enable")
	// validatoSetName := command.String("name", homePath, "enable")

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

		networkState, err := state.LoadNetworkState(*conf.OutputDir)
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
			log.Fatalf("Failed to %s network: %s", subcommand, err)
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
			if err := stop(ctx, networkState); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := stop(ctx, networkState); err != nil {
				log.Fatal(err)
			}
			cleanup(*networkHomePath)
		}
	case "nomad":
		if err := nmrunner.StartAgent(*configFilePath); err != nil {
			log.Fatal(err)
		}
	case "list-validators":
		networkState, err := state.LoadNetworkState(*networkHomePath)
		if err != nil {
			log.Fatalf("Failed list validators: %s", err)
		}

		if networkState.Empty() {
			log.Fatalf("Failed list validators: network not bootstrapped. Use the 'bootstrap' subcommand or provide different network home with the `-home-path` flag")
		}

		validators := networkState.ListValidators()

		validatorsJson, err := json.MarshalIndent(validators, "", "\t")
		if err != nil {
			log.Fatalf("failed to marshal validators info: %s", err.Error())
		}

		fmt.Println(string(validatorsJson))
	case "add-validator-set":
		networkState, err := state.LoadNetworkState(*networkHomePath)
		if err != nil {
			log.Fatalf("Failed to %s network: %s", subcommand, err)
		}

		if networkState.Empty() {
			log.Fatalf("Failed to %s network: network not bootstrapped. Use the 'bootstrap' subcommand or provide different network home with the `-home-path` flag", subcommand)
		}

		if err := addValidator(*networkState); err != nil {
			log.Fatal(err)
		}
	default:
		log.Printf("unknown subcommand %s: expected 'generate'|'bootstrap'|'start'|'stop'|'destroy'|'nomad'|'add-validator-set' subcommands", subcommand)
		os.Exit(1)
	}
}
