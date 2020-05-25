package sagify

import (
	"errors"
	"strconv"
	"time"

	"github.com/kenza-ai/kenza/worker/job"
)

// Sagify build handler
type Sagify struct {
	commands       []string
	next           job.Handler
	config         *configuration
	buildConfig    BuildConfiguration
	job            *job.Request
	sageMakerJobID string
}

type configuration struct {
	Region               string `json:"aws_region"`
	RequirementsFilePath string `json:"requirements_dir"`
}

// BuildConfiguration - Sagify build configuration i.e. the commands in the build file
type BuildConfiguration struct {
	Train                Train
	BatchTransform       BatchTransform       `yaml:"batch_transform"`
	HyperparameterTuning HyperparameterTuning `yaml:"hyperparameter_optimization"`
}

// Train command configuration
type Train struct {
	Timeout         string
	Metrics         string
	Ec2Type         string `yaml:"ec2_type"`
	InputDir        string `yaml:"input_s3_dir"`
	OutputDir       string `yaml:"output_s3_dir"`
	VolumeSize      string `yaml:"volume_size"`
	BaseJobName     string `yaml:"base_job_name"`
	HyperparamsFile string `yaml:"hyperparameters_file"`
	Deploy          Deploy
	Schedules       *Schedules `yaml:"schedule"`
}

// Deploy command configuration
type Deploy struct {
	Endpoint      string
	Ec2Type       string `yaml:"ec2_type"`
	Instances     string `yaml:"instances_count"`
	ModelLocation string `yaml:"s3_model_location"`
}

// HyperparameterTuning command configuration
type HyperparameterTuning struct {
	Timeout              string
	Deploy               Deploy
	Ec2Type              string     `yaml:"ec2_type"`
	MaxJobs              string     `yaml:"max_jobs"`
	VolumeSize           string     `yaml:"volume_size"`
	InputDir             string     `yaml:"input_s3_dir"`
	OutputDir            string     `yaml:"output_s3_dir"`
	BaseJobName          string     `yaml:"base_job_name"`
	MaxParallelJobs      string     `yaml:"max_parallel_jobs"`
	HyperparamRangesFile string     `yaml:"hyperparameter_ranges_file"`
	Schedules            *Schedules `yaml:"schedule"`
}

// BatchTransform command configuration
type BatchTransform struct {
	Ec2Type       string     `yaml:"ec2_type"`
	Instances     string     `yaml:"instances_count"`
	ModelLocation string     `yaml:"model_s3_location"`
	InputDir      string     `yaml:"features_s3_location"`
	OutputDir     string     `yaml:"predictions_s3_location"`
	Schedules     *Schedules `yaml:"schedule"`
}

var (
	// errMissingDeploymentEndpointName is returned when the mandatory endpoint name is missing from the config file.
	errMissingDeploymentEndpointName = errors.New("Deploy entry is missing mandatory endpoint key and/or value")
)

// NewSagify - Sagify initializer
func NewSagify(buildConfig BuildConfiguration) *Sagify {
	return &Sagify{config: &configuration{}, buildConfig: buildConfig}
}

// Handle - Sagify's Handler implementation
//
// Parses build file and runs Sagify commands.
func (h *Sagify) Handle(r *job.Request) {
	h.job = r

	// IMPORTANT: set a unique SageMaker ID. As of this writing, SageMaker does not offer a job deletion API.
	// Setting the ID to a pseudo-constant value e.g. "kenza-{projectID}-{jobID}" will cause issues among
	// Kenza installations targeting the same AWS account. https://forums.aws.amazon.com/thread.jspa?threadID=268989
	projectID, jobID := toString(h.job.JobQueued.ProjectID), toString(h.job.JobQueued.JobID)
	h.sageMakerJobID = "kenza-" + projectID + "-" + jobID + "-" + toString(time.Now().UTC().UnixNano()) // avoid SageMaker ID conflicts
	r.SageMakerJobID = h.sageMakerJobID

	h.parseConfiguration()

	// Install project dependencies and extract sagify version / info and apply
	// AWS region. These happen regardless of job type (scheduled or normal) and
	// always happen first so we can determine if the requested Sagify version is supported.
	h.addPipInstallCommand()
	h.addSagifyInfoCommand()
	r.Region = h.config.Region

	if err := h.addTrainOrTuningCommand(); err != nil {
		r.Fail(err)
		return
	}
	h.addBatchTransformCommand()
	if err := h.job.Notify(); err != nil {
		e(err.Error())
	}

	h.notifySchedules()

	if err := run(h.commands, r.WorkDir, map[string]string{},
		h.buildConfig.HyperparameterTuning.winningModelResolver); err != nil {
		r.Fail(err)
		return
	}

	r.Status = "completed"
	if err := r.Notify(); err != nil {
		e(err.Error())
	}
}

// SetNext sets the next Handler
func (h *Sagify) SetNext(next job.Handler) {
	h.next = next
}

func toString(id int64) string {
	return strconv.FormatInt(id, 10)
}

func (h *Sagify) notifySchedules() {
	if err := h.job.NotifySchedules(); err != nil {
		e(err.Error())
	}
}
