package cli

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execute(cmd []string) error {
	commands := append([]string{"-c"}, cmd...)
	output, err := exec.Command("sh", commands...).CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))
	return err
}

func executeNoPrint(cmd []string) (out string, err error) {
	commands := append([]string{"-c"}, cmd...)
	output, err := exec.Command("sh", commands...).CombinedOutput()
	return string(output), err
}

func executeWithErrOut(cmd []string, stderr *bytes.Buffer, stdout *bytes.Buffer) error {
	commands := append([]string{"-c"}, cmd...)
	command := exec.Command("sh", commands...)

	command.Stdout = stdout
	command.Stderr = stderr
	return command.Run()
}
