package sagify

import "fmt"

func (h *Sagify) addPushCommand() {
	h.commands = append(h.commands, fmt.Sprint("sagify -v push"))
}
