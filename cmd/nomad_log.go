package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var nomadJobID string

var nomadLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Log a specific Nomad job",
	RunE: func(cmd *cobra.Command, args []string) error {
		exec, err := os.Executable()
		if err != nil {
			log.Println(err)
		}

		fmt.Println("Exec path: ", exec)

		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		log.Printf("Logging for job %q in path %q", nomadJobID, path)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigs

		log.Printf("Recived signal: %s", sig)

		return nil
	},
}

func init() {
	nomadLogCmd.PersistentFlags().StringVar(&nomadJobID,
		"job-id",
		"",
		"Nomad ID of the job to be logged",
	)
}
