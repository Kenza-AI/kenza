package cli

import (
	"github.com/spf13/cobra"
)

// AvailServices avaiable as args in kenza subcommands
var kenzaServices = []string{"api", "db", "web", "progress", "pubsub", "worker", "scheduler"}

func init() {
	rootCmd.AddCommand(sshCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(scaleCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(provisionCmd)
}

const rootCommandLongDescription = `
A Machine Learning focused CI/CD pipeline system for the cloud and container age.

Complete documentation available at https://docs.kenza.ai
`

var rootCmd = &cobra.Command{
	Use:   "kenza",
	Short: "Kenza is a Machine Learning CI/CD pipeline system.",
	Long:  rootCommandLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute the "kenza" command
func Execute() error {
	return rootCmd.Execute()
}
