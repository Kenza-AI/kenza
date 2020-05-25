package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const sshLongDescription = `
The ssh command executes a command on the Kenza manager machine.

If you wish to SSH into the machine and obtain a shell, use "docker-machine ssh your_machine_name" instead.
`

const sshExamples = `
To run a command (e.g. ls) using the Kenza name used when the "kenza provision" command was run:
kenza ssh your_machine_name ls
`

var sshCmd = &cobra.Command{
	Use:     "ssh machine_name command_to_execute",
	Short:   "Executes a command on a Kenza manager machine",
	Long:    sshLongDescription,
	Example: sshExamples,
	Args:    cobra.ExactArgs(2), // name of machine to SSH into and command to be executed

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf(`Attempting to run command "%s" on machine "%s"`, args[1], args[0]))
		execute(sshCommand(args[0], args[1]))
	},
}

func sshCommand(machineName string, command string) []string {
	return []string{fmt.Sprintf("docker-machine ssh %s %s", machineName, command)}
}
