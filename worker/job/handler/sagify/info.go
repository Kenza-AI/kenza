package sagify

import "fmt"

func (h *Sagify) addSagifyInfoCommand() {
	h.commands = append(h.commands, fmt.Sprint("pip3 show sagify"))
}
