package sagify

import "fmt"

func (h *Sagify) addBuildCommand() {
	h.commands = append(h.commands, fmt.Sprint("sagify -v build"))
}
