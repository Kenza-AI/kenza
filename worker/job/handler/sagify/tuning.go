package sagify

import (
	"fmt"
)

func (h *Sagify) addTuningCommand() error {
	if tuning := h.buildConfig.HyperparameterTuning; (tuning != HyperparameterTuning{}) {
		h.job.Type = "tuning"
		h.job.Schedules = tuning.Schedules

		h.commands = append(h.commands,
			fmt.Sprintf(`sagify -v cloud hyperparameter-optimization -i %s -o %s -h %s -e %s -v %s -s %s -n %s --job-name %s -m %s -p %s --wait`,
				tuning.InputDir, tuning.OutputDir, tuning.HyperparamRangesFile,
				tuning.Ec2Type, tuning.VolumeSize, tuning.Timeout, tuning.BaseJobName, h.sageMakerJobID, tuning.MaxJobs, tuning.MaxParallelJobs))
		if err := h.addDeployCommand(); err != nil {
			return err
		}
	}
	return nil
}

func (h *HyperparameterTuning) winningModelResolver(winningJobID string) string {
	return h.OutputDir + "/" + winningJobID + "/output/model.tar.gz"
}
