package cmd

import "github.com/spf13/cobra"

var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Allos to start/stop a specific jobs",
	Long:  "The command allows some operations on a specific job.",
}

func init() {
	jobsCmd.AddCommand(jobsStartCmd)
	jobsCmd.AddCommand(jobsStopCmd)
}
