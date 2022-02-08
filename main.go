package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func generate(config *Config) ([]nodeSet, error) {
	if _, err := os.Stat(config.OutputDir); os.IsExist(err) {
		return nil, fmt.Errorf("output directory %q already exist", config.OutputDir)
	}

	log.Println("generating network")

	gen, err := NewGenerator(config)
	if err != nil {
		return nil, err
	}

	nodeSets, err := gen.Generate()
	if err != nil {
		return nil, err
	}

	if err := config.Persist(); err != nil {
		return nil, fmt.Errorf("failed to persist config in output directory %s", config.OutputDir)
	}

	log.Println("generating network success")

	return nodeSets, nil
}

func start(config *Config) error {
	log.Println("starting network")
	nodeSets, err := generate(config)
	if err != nil {
		return fmt.Errorf("failed to generate config for network: %w", err)
	}

	nomadRunner, err := NewNomadRunner(nil)
	if err != nil {
		return err
	}

	runner := NewRunner(nomadRunner)

	if err := runner.StartRawNetwork(config, nodeSets); err != nil {
		return fmt.Errorf("failed to start nomad network: %s", err)
	}

	log.Println("starting network success")
	return nil
}

func stop() {
	log.Println("stopping network")
	nomadRunner, err := NewNomadRunner(nil)
	if err != nil {
		log.Fatal(err)
	}

	runner := NewRunner(nomadRunner)

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

		config, err := ParseConfigFile(*configFilePathS)
		if err != nil {
			log.Fatal(err)
		}

		if err := start(config); err != nil {
			log.Fatal(err)
		}

	case "stop":
		stop()
	case "generate":
		if err := generateCmd.Parse(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

		config, err := ParseConfigFile(*configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := generate(config); err != nil {
			log.Fatal(err)
		}

	// case "destroy":
	// destroy()
	default:
		log.Printf("unknown subcommand %s: expected 'start'|'stop'|'destroy' subcommands", arg)
		os.Exit(1)
	}

}

func strPoint(s string) *string {
	return &s
}

func intPoint(i int) *int {
	return &i
}
