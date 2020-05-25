package cli

import (
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs SERVICE",
	Short: "View service log output",

	ValidArgs: kenzaServices,
	Args:      cobra.ExactValidArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		execute(logsCommand(args))
	},
}

func logsCommand(args []string) []string {
	command := "docker service logs"
	for _, arg := range args {
		command = command + " " + "kenza_" + arg
	}
	return []string{command}
}
