package sagify

import (
	"os/exec"
	"strings"
)

func execute(cmd string) (output string, err error) {
	i("executing command: %v", cmd)

	command := strings.Split(cmd, " ")
	out, err := exec.Command(command[0], command[1:]...).CombinedOutput()
	output = string(out)
	i(output)

	return output, err
}
