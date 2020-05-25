package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [SERVICE...]",
	Short: "Stop Kenza services",
	Long:  "Stop Kenza services. Does not result in volume / data loss.",

	ValidArgs: kenzaServices,
	Args:      cobra.OnlyValidArgs,

	Run: func(cmd *cobra.Command, args []string) {
		if err := execute(stopCommand(args)); err == nil {
			fmt.Println("Stopped successfully")
		}
	},
}

func stopCommand(args []string) []string {
	ackMessage := "Stopping"
	if len(args) == 0 {
		fmt.Println(ackMessage + " Kenza")
		return []string{"docker stack rm kenza"}
	}

	// https://docs.docker.com/engine/reference/commandline/service_scale/#scale-multiple-services
	command := "docker service scale"
	for _, arg := range args {
		command = command + " " + "kenza_" + arg + "=0"
		ackMessage = ackMessage + " " + arg
	}
	fmt.Println(command)
	fmt.Println(ackMessage)
	return []string{command}
}
