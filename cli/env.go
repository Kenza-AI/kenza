package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const envLongDescription = `
The env command prints the currently "active" Kenza installation (and its underlying Docker Machine equivalent).
Use it with "eval" to "activate" a specific Kenza installation.
`

const envExamples = `
To "enable" an existing Kenza installation named "kenza-aws":
eval $(kenza env kenza-aws)
`

var envCmd = &cobra.Command{
	Use:     "env [KENZA_INSTALLATION_NAME]",
	Short:   "",
	Long:    envLongDescription,
	Example: envExamples,
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		command := fmt.Sprintf("docker-machine env %s", args[0])
		output, err := executeNoPrint([]string{command})
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		printKenzaMachineEnvironment(output)
	},
}

// Sanitized `docker-machine env machine_name` output; skips lines not describing env vars
func printKenzaMachineEnvironment(dockerMachineEnvCommandOutput string) {
	scanner := bufio.NewScanner(strings.NewReader(string(dockerMachineEnvCommandOutput)))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "export") {
			continue
		}
		fmt.Println(line)
	}
}
