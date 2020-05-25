package sagify

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
)

const supportedSagifyVersions string = ">0.0.0 <1.0.0"

// Hyperaparameter tuning jobs produce a "winner" job/model. There is no way to know the winner in advance,
// so we keep track of the "winning" job (produced as an output of hyperparameter tuning commands)
// so it can be fed to the deploy command (otherwise there's no way to know which model to deploy).
var winningTuningJob string

// Provides the model location for a winning tuning job (needed to deploy the winning model).
type hyperparameterTuningWinningModelResolver = func(winningJobID string) string

func run(commands []string, workDir string, variables map[string]string,
	hyperparameterTuningWinningModelResolver func(winningJobID string) string) error {
	os.Chdir(workDir)

	awsRole := variables["IAM_ROLE"]
	awsExternalID := variables["EXTERNAL_ID"]

	for _, cmd := range commands {

		if canUseAwsRoleAndExternalID(cmd) {
			if awsRole != "" {
				cmd = commandByAppendingArgument("--iam-role-arn", awsRole, cmd)
			}
			if awsExternalID != "" {
				cmd = commandByAppendingArgument("--external-id", awsExternalID, cmd)
			}
		}

		if isDeployCommand(cmd) && winningTuningJob != "" {
			cmd += fmt.Sprintf(" -m %s", hyperparameterTuningWinningModelResolver(winningTuningJob))
		}

		output, err := execute(cmd)

		if isHyperparameterTuningCommand(cmd) {
			winningTuningJob = getWinningHyperparameterTuningTrainingJob(output)
		}

		if isSagifyInfoCommand(cmd) {
			sagifyVersion := getSagifyVersionFromPipShowInfo(output)
			if err := isSupportedSagifyVersion(sagifyVersion, supportedSagifyVersions); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func canUseAwsRoleAndExternalID(cmd string) bool {
	return isPushCommand(cmd) || isTrainCommand(cmd) || isDeployCommand(cmd) || isHyperparameterTuningCommand(cmd) || isBatchTransformCommand(cmd)
}

func isSagifyInfoCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "pip3 show sagify")
}

func isPushCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "sagify -v push")
}

func isTrainCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "sagify -v cloud train")
}

func isDeployCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "sagify -v cloud deploy")
}

func isHyperparameterTuningCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "sagify -v cloud hyperparameter-optimization")
}

func isBatchTransformCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "sagify -v cloud batch-transform")
}

func commandByAppendingArgument(arg string, value string, cmd string) string {
	return strings.Join([]string{cmd, arg, value}, " ")
}

func isSupportedSagifyVersion(version, supportedVersionRange string) error {
	v, err := semver.Parse(version)
	if err != nil {
		return err
	}

	expectedRange, err := semver.ParseRange(supportedVersionRange)
	if err != nil {
		return err
	}

	if !expectedRange(v) {
		return fmt.Errorf("Sagify version %s outside of supported range %s", v, supportedVersionRange)
	}
	return nil
}

func getSagifyVersionFromPipShowInfo(info string) string {
	scanner := bufio.NewScanner(strings.NewReader(info))
	for scanner.Scan() {
		infoLine := scanner.Text()
		pipShowVersionPrefix := "Version: "
		if strings.HasPrefix(infoLine, pipShowVersionPrefix) {
			return strings.TrimPrefix(infoLine, pipShowVersionPrefix)
		}
	}
	return ""
}

func getWinningHyperparameterTuningTrainingJob(tuningOutput string) string {
	scanner := bufio.NewScanner(strings.NewReader(tuningOutput))
	for scanner.Scan() {
		outputLine := scanner.Text()
		winningJobPrefix := "Best job name: "
		if strings.HasPrefix(outputLine, winningJobPrefix) {
			return strings.TrimPrefix(outputLine, winningJobPrefix)
		}
	}
	return ""
}
