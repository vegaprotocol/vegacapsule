package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator"
	"code.vegaprotocol.io/vegacapsule/runner"
	"code.vegaprotocol.io/vegacapsule/runner/nomad"
	"code.vegaprotocol.io/vegacapsule/types"
)

func generate(conf *config.Config) (*types.GeneratedServices, error) {
	if _, err := os.Stat(conf.OutputDir); os.IsExist(err) {
		return nil, fmt.Errorf("output directory %q already exist", conf.OutputDir)
	}

	log.Println("generating network")

	gen, err := generator.New(conf)
	if err != nil {
		return nil, err
	}

	generatedSvcs, err := gen.Generate()
	if err != nil {
		return nil, err
	}

	if err := conf.Persist(); err != nil {
		return nil, fmt.Errorf("failed to persist config in output directory %s", conf.OutputDir)
	}

	log.Println("generating network success")

	return generatedSvcs, nil
}

func start(conf *config.Config) error {
	log.Println("starting network")
	generatedSvcs, err := generate(conf)
	if err != nil {
		return fmt.Errorf("failed to generate config for network: %w", err)
	}

	nomadRunner, err := nomad.New(nil)
	if err != nil {
		return err
	}

	runner := runner.New(nomadRunner)

	for _, dc := range conf.Network.PreStart.Docker {
		if err := runner.RunDockerJob(dc); err != nil {
			return fmt.Errorf("failed to run pre start job %s: %w", dc.Name, err)
		}
	}

	if err := runner.StartRawNetwork(conf, generatedSvcs); err != nil {
		return fmt.Errorf("failed to start nomad network: %s", err)
	}

	log.Println("starting network success")
	return nil
}

func stop() {
	log.Println("stopping network")
	nomadRunner, err := nomad.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	runner := runner.New(nomadRunner)

	if err := runner.StopRawNetwork(); err != nil {
		log.Fatalf("failed to start nomad network: %s", err)
	}
	log.Println("stopping network success")
}

func destroy(outputDir string) {
	log.Println("destroying network")
	stop()

	if err := os.RemoveAll(outputDir); err != nil {
		log.Fatalf("failed to destroy network: %s", err)
	}
	log.Println("destroying network success")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'start'|'stop'|'destroy' subcommands")
		os.Exit(1)
	}

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	configFilePathS := startCmd.String("config-path", "", "enable")

	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	configFilePath := generateCmd.String("config-path", "", "enable")

	arg := os.Args[1]
	switch arg {
	case "start":
		if err := startCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		conf, err := config.ParseConfigFile(*configFilePathS)
		if err != nil {
			log.Fatal(err)
		}

		if err := start(conf); err != nil {
			log.Fatal(err)
		}

	case "stop":
		stop()
	case "generate":
		if err := generateCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		conf, err := config.ParseConfigFile(*configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := generate(conf); err != nil {
			log.Fatal(err)
		}

	// case "destroy":
	// destroy()
	default:
		log.Printf("unknown subcommand %s: expected 'start'|'stop'|'destroy' subcommands", arg)
		os.Exit(1)
	}

}
