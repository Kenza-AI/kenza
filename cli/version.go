package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version - source control release tag
	Version = "dev"
	// Commit ID
	Commit = "none"
	// Date binary was built
	Date = "unknown"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Args:  cobra.NoArgs,
	Short: "View Kenza configuration information",

	Run: func(cmd *cobra.Command, args []string) {
		runInfoCommand()
	},
}

func runInfoCommand() {
	fmt.Println("Kenza info")
	fmt.Println()

	fmt.Printf("Version: v%s\n", Version)
	fmt.Printf("Built:   %s\n", Date)
	fmt.Printf("Commit:  %s\n", Commit)

	fmt.Println()
}
