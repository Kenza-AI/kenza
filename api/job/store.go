package job

import (
	"time"

	"github.com/kenza-ai/kenza/event"
)

// The Job model.
type Job struct {
	ID                    int64                 `json:"id"`
	Status                string                `json:"status"`
	Type                  string                `json:"type"`
	Submitter             string                `json:"submitter"`
	Project               Project               `json:"project"`
	CommitID              string                `json:"commitID"`
	Started               time.Time             `json:"started"`
	Created               string                `json:"created"`
	Updated               string                `json:"updated"`
	SageMakerID           string                `json:"sagemakerID"`
	Endpoint              string                `json:"endpoint,omitempty"`
	Region                string                `json:"region"`
	EndpointInfo          interface{}           `json:"endpoint_info"`
	SageMakerTuningJob    SageMakerTuningJob    `json:"sagemaker_tuning_job,omitempty"`
	SageMakerTrainingJob  SageMakerTrainingJob  `json:"sagemaker_training_job,omitempty"`
	SageMakerTransformJob SageMakerTransformJob `json:"sagemaker_transform_job,omitempty"`
}

// Project model from a Job's point of view.
type Project struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Repo string `json:"repo"`
}

// Metric info, maps SageMaker's `MetricData` struct
type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// SageMakerTrainingJob model
type SageMakerTrainingJob struct {
	Metrics                   []Metric           `json:"metrics"`
	HyperParameters           map[string]*string `json:"hyperparameters"`
	S3ModelLocation           string             `json:"s3_model_location"`
	TrainingTimeInSeconds     int64              `json:"training_time_seconds"`
	ResourceConfig            ResourceConfig     `json:"resource_config"`
	DescribeTrainingJobOutput interface{}
}

// SageMakerTransformJob model
type SageMakerTransformJob struct {
	S3InputLocation            string `json:"s3_input_location"`
	S3OutputLocation           string `json:"s3_output_location"`
	DescribeTransformJobOutput interface{}
}

// SageMakerTuningJob model
type SageMakerTuningJob struct {
	S3InputLocation                       string `json:"s3_input_location"`
	S3OutputLocation                      string `json:"s3_output_location"`
	DescribeHyperParameterTuningJobOutput interface{}
}

// ResourceConfig describes the resources, including ML compute instances and ML storage volumes,
// to use for model training.
type ResourceConfig struct {
	VolumeSizeInGB int64  `json:"volume_size_gb"`
	InstanceType   string `json:"instance_type"`
	InstanceCount  int64  `json:"instance_count"`
}

// Store - job store abstraction
type Store interface {
	// Delete a list of jobs.
	DeleteJobs(accountID, projectID int64, jobIDs []int64) error

	// Cancel a list of jobs.
	CancelJobs(accountID, projectID int64, jobIDs []int64) error

	// Get a job's info.
	Get(jobID int64) (Job, error)

	// GetAll returns a project's jobs.
	GetAll(accountID, projectID int64) ([]Job, error)

	// Create a job based on the GitHub hook info.
	CreateJob(accountID, projectID int64, submitter, deliveryID, revisionID string) (jobID int64, err error)

	// UpdateJob persists changes to a job.
	UpdateJob(job event.JobUpdated) error
}
