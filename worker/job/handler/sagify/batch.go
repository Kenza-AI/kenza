package sagify

import (
	"fmt"
)

func (h *Sagify) addBatchTransformCommand() {
	if batch := h.buildConfig.BatchTransform; (batch != BatchTransform{}) {
		h.job.Type = "batchtransform"
		h.job.Schedules = batch.Schedules

		featuresLocation, predictionsLocation := batch.InputDir, batch.OutputDir

		h.commands = append(h.commands,
			fmt.Sprintf(`sagify -v cloud batch-transform -m %s -i %s -o %s -n %s -e %s --job-name %s --wait`,
				batch.ModelLocation, featuresLocation, predictionsLocation, batch.Instances, batch.Ec2Type, h.sageMakerJobID))
	}
}
