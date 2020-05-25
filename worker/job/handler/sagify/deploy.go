package sagify

import (
	"fmt"
)

func (h *Sagify) addDeployCommand() error {
	// Deployments are only added after we have established the job type (e.g. training)
	// meaning it's safe to assume the corresponding deployment config has been set by this point.
	deploymentConfig := h.buildConfig.Train.Deploy
	if h.job.Type == "tuning" {
		deploymentConfig = h.buildConfig.HyperparameterTuning.Deploy
	}

	if deploy := deploymentConfig; (deploy != Deploy{}) {
		if deploy.Endpoint == "" {
			return errMissingDeploymentEndpointName
		}

		h.job.Endpoint = deploy.Endpoint

		if err := h.job.Notify(); err != nil {
			e(err.Error())
		}

		deployCommand := fmt.Sprintf(`sagify -v cloud deploy -n %s -e %s --endpoint-name %s`,
			deploy.Instances, deploy.Ec2Type, deploy.Endpoint)
		if h.job.Type == "training" {
			deployCommand += fmt.Sprintf(" -m %s", h.buildConfig.Train.modelLocation(h.sageMakerJobID))
		}

		h.commands = append(h.commands, deployCommand)
	}
	return nil
}
