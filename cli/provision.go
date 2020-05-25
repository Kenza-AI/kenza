package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const provisionLongDescription = `
The provision command prepares an environment based on a driver (e.g. amazonec2) for Kenza to run in.
`

const provisionExamples = `
To provision resources for a Kenza installation named 'kenza-aws' that uses the 'ml-role' IAM role on AWS EC2:
kenza provision --driver amazonec2 --amazonec2-iam-instance-profile ml-role kenza-aws
`

var provisionCmd = &cobra.Command{
	Use:     "provision [flags]",
	Short:   "Provisions resources Kenza uses to operate on the cloud",
	Long:    provisionLongDescription,
	Example: provisionExamples,

	// pass all args and flags as-is to "docker-machine create"
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Kenza resource provisioning...")
		fmt.Println("Preparing EC2 environment (manager EC2 instance, Security Group etc)")
		execProvisionCommand(args)
	},
}

func execProvisionCommand(args []string) {
	// Request Security Group to open ports 8080 (Kenza API) and 80 (Kenza wep app)
	argsString := " --amazonec2-open-port 8080 --amazonec2-open-port 80"

	// Then append the rest of args from the cli as-is
	for _, arg := range args {
		argsString += fmt.Sprintf(" %s", arg)
	}

	command := fmt.Sprint("docker-machine create" + argsString)
	if err := execute([]string{command}); err != nil {
		os.Exit(2)
	}
}
