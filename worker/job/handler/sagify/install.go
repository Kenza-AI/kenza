package sagify

import (
	"fmt"
	"path/filepath"
)

func (h *Sagify) addPipInstallCommand() {
	h.commands = append(h.commands, fmt.Sprintf("pip3 install -r %s", filepath.Join(h.job.WorkDir, h.config.RequirementsFilePath)))
}
