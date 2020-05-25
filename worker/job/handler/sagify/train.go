package sagify

import (
	"fmt"
	"unicode"
)

func (h *Sagify) addTrainOrTuningCommand() error {
	h.addBuildCommand()
	h.addPushCommand()

	// Train and Tuning are mutually exclusive. Train takes precedence.
	if train := h.buildConfig.Train; (train != Train{}) {
		h.job.Type = "training"
		h.job.Schedules = train.Schedules

		trainCommand := fmt.Sprintf(`sagify -v cloud train -i %s -o %s -e %s -v %s -s %s -n %s --job-name %s`,
			train.InputDir, train.OutputDir, train.Ec2Type, train.VolumeSize,
			train.Timeout, train.BaseJobName, h.sageMakerJobID)

		if train.Metrics != "" {
			trainCommand += fmt.Sprintf(` --metric-names %s`, removeWhiteSpace(train.Metrics))
		}

		if train.HyperparamsFile != "" {
			trainCommand += fmt.Sprintf(` --hyperparams-file %s`, train.HyperparamsFile)
		}

		h.commands = append(h.commands, trainCommand)

		if err := h.addDeployCommand(); err != nil {
			return err
		}
	} else {
		if err := h.addTuningCommand(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Train) modelLocation(jobID string) string {
	return t.OutputDir + "/" + jobID + "/output/model.tar.gz"
}

func removeWhiteSpace(s string) string {
	trimmed := []rune{}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			trimmed = append(trimmed, r)
		}
	}
	return string(trimmed)
}
