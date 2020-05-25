package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const scaleLongDescription = `
The scale command scales one or multiple Kenza services. 

The command returns immediately, but the actual scaling of the service may take some time. 
To stop a service yo,u can set the scale to 0 or run 'kenza stop SERVICE_NAME'.
`

const scaleExamples = `
To be able to run 3 jobs simultaneously, scale the worker service accordingly:
kenza scale worker=3
`

var scaleCmd = &cobra.Command{
	Use:     "scale [SERVICE=NUMBER_OF_REPLICAS...]",
	Example: scaleExamples,
	Long:    scaleLongDescription,
	Short:   "Scale one or multiple Kenza services",
	Args:    cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		execute(scaleCommand(args))
	},
}

func scaleCommand(args []string) []string {
	ackMessage := "Scaling Kenza"

	// https: //docs.docker.com/engine/reference/commandline/service_scale/#scale-multiple-services
	command := "docker service scale"
	for _, arg := range args {
		command = command + " " + "kenza_" + arg
		ackMessage = ackMessage + " " + arg
	}
	fmt.Printf(ackMessage+" (%s)", command)
	return []string{command}
}
