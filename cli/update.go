package cli

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

var verbose bool

func init() {
	parseUpdateCmdFlags()
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates Kenza to the latest version",
	Args:  cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		doSelfUpdate(verbose)
	},
}

func doSelfUpdate(verbose bool) {
	fmt.Println("Current Kenza version:", Version)

	if verbose {
		selfupdate.EnableLog()
	}

	fmt.Println("Checking for newer Kenza versions...")
	v := semver.MustParse(strings.TrimPrefix(Version, "v"))
	latest, err := selfupdate.UpdateSelf(v, "Kenza-AI/kenza")
	if err != nil {
		fmt.Println("Binary update failed:", err)
		return
	}

	if latest.Version.Equals(v) {
		fmt.Println("Current binary is the latest version", Version)
	} else {
		fmt.Println("Successfully updated to version", latest.Version)
		fmt.Println("Release note:\n", latest.ReleaseNotes)
	}
}

func parseUpdateCmdFlags() {
	updateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output for the update command")
}
