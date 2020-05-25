package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const statusExamples = `
Get the status of all services:
kenza status

Get the status of only the web and progress services:
kenza status web progress
`

var statusCmd = &cobra.Command{
	Use:     "status [SERVICE...]",
	Short:   "Get service(s) status information",
	Example: statusExamples,

	ValidArgs: kenzaServices,
	Args:      cobra.OnlyValidArgs,

	Run: func(cmd *cobra.Command, args []string) {
		execute(statusCommand(args))
	},
}

func statusCommand(args []string) []string {
	// https://docs.docker.com/engine/reference/commandline/stack_ps/#filtering
	return []string{"docker stack ps" + " --no-trunc" + serviceFilters(args) + " kenza"}
}

func serviceFilters(args []string) string {
	serviceFilters := ""
	for _, arg := range args {
		serviceFilters = serviceFilters + fmt.Sprintf(` --filter "name=kenza_%s"`, arg)
	}
	return serviceFilters
}
