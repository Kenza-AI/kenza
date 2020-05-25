package main

import (
	"fmt"
	"os"

	"github.com/kenza-ai/kenza/cli"
)

var (
	version string
	commit  string
	date    string
)

func init() {
	cli.Version = version
	cli.Commit = commit
	cli.Date = date
}

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
