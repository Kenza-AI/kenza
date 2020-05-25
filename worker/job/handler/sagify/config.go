package sagify

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func (h *Sagify) parseConfiguration() {
	configFilePath := filepath.Join(h.job.WorkDir, ".sagify.json")
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		h.job.Fail(err)
	}

	if err := json.Unmarshal(configData, h.config); err != nil {
		h.job.Fail(err)
	}
	i("sagify config %+v", *h.config)
}
