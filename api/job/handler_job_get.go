package job

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/kenza-ai/kenza/api/httputil"
	"github.com/kenza-ai/kenza/api/sagemaker"
)

// Get returns a job's info.
func Get(jobs Store, sm *sagemaker.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobID, err := strconv.ParseInt(httputil.Param(r, "jobID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		job, err := jobs.Get(jobID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := appendSageMakerInfoIfNeeded(&job, sm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := appendEndpointInfoIfNeeded(&job, sm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}

// GetLog returns a job's log file
func GetLog(logfileDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobID, projectID := httputil.Param(r, "jobID"), httputil.Param(r, "projectID")
		pathToLogfile := filepath.Join(logfileDir, projectID+"-"+jobID+".log")

		i("Serving log for project %s, job %s from path %s", projectID, jobID, pathToLogfile)
		http.ServeFile(w, r, pathToLogfile)
	})
}

// GetAll returns a project's jobs.
// TODO(ilazakis): paging, hard limit or both; job count can get out of hand quickly.
func GetAll(jobs Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID, err := strconv.ParseInt(httputil.Param(r, "accountID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		projectID, err := strconv.ParseInt(httputil.Param(r, "projectID"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jobs, err := jobs.GetAll(accountID, projectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(jobs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}

func appendSageMakerInfoIfNeeded(job *Job, sm *sagemaker.Client) error {
	if len(job.SageMakerID) < 1 {
		return nil
	}

	switch job.Type {
	case "batch":
		return appendSageMakerTransformJobInfo(job, sm)
	case "tuning":
		return appendSageMakerTuningJobInfo(job, sm)
	case "training":
		return appendSageMakerTrainingJobInfo(job, sm)
	default:
		return nil
	}
}

func appendSageMakerTrainingJobInfo(job *Job, sm *sagemaker.Client) error {
	sageMakerTrainingJob := SageMakerTrainingJob{}

	sagemakerJob, err := sm.GetTrainingJob(job.SageMakerID, job.Region, "")
	if err != nil {
		return awsErr(err)
	}

	metrics := []Metric{}
	for _, metric := range sagemakerJob.FinalMetricDataList {
		m := Metric{
			Name:  *metric.MetricName,
			Value: *metric.Value,
		}
		metrics = append(metrics, m)
	}

	sageMakerTrainingJob.Metrics = metrics

	sageMakerTrainingJob.ResourceConfig.InstanceType = *sagemakerJob.ResourceConfig.InstanceType
	sageMakerTrainingJob.ResourceConfig.InstanceCount = *sagemakerJob.ResourceConfig.InstanceCount
	sageMakerTrainingJob.ResourceConfig.VolumeSizeInGB = *sagemakerJob.ResourceConfig.VolumeSizeInGB

	sageMakerTrainingJob.HyperParameters = sagemakerJob.HyperParameters

	sageMakerTrainingJob.TrainingTimeInSeconds = *sagemakerJob.TrainingTimeInSeconds

	sageMakerTrainingJob.S3ModelLocation = *sagemakerJob.ModelArtifacts.S3ModelArtifacts

	job.SageMakerTrainingJob = sageMakerTrainingJob
	job.SageMakerTrainingJob.DescribeTrainingJobOutput = sagemakerJob
	return nil
}

func appendSageMakerTransformJobInfo(job *Job, sm *sagemaker.Client) error {
	sageMakerTransformJob := SageMakerTransformJob{}

	sagemakerJob, err := sm.GetTransformJob(job.SageMakerID, job.Region, "")
	if err != nil {
		return awsErr(err)
	}

	sageMakerTransformJob.S3InputLocation = *sagemakerJob.TransformInput.DataSource.S3DataSource.S3Uri
	sageMakerTransformJob.S3OutputLocation = *sagemakerJob.TransformOutput.S3OutputPath
	job.SageMakerTransformJob = sageMakerTransformJob
	job.SageMakerTransformJob.DescribeTransformJobOutput = sagemakerJob

	return nil
}

func appendSageMakerTuningJobInfo(job *Job, sm *sagemaker.Client) error {
	sageMakerTuningJob := SageMakerTuningJob{}

	sagemakerJob, err := sm.GetTuningJob(job.SageMakerID, job.Region, "")
	if err != nil {
		return awsErr(err)
	}

	sageMakerTuningJob.S3InputLocation = *sagemakerJob.TrainingJobDefinition.InputDataConfig[0].DataSource.S3DataSource.S3Uri
	sageMakerTuningJob.S3OutputLocation = *sagemakerJob.TrainingJobDefinition.OutputDataConfig.S3OutputPath

	sageMakerTuningJob.DescribeHyperParameterTuningJobOutput = sagemakerJob

	job.SageMakerTuningJob = sageMakerTuningJob

	return nil
}

func appendEndpointInfoIfNeeded(job *Job, sm *sagemaker.Client) error {
	if job.Endpoint == "" {
		return nil
	}

	endpoint, err := sm.GetEndpoint(job.Endpoint, job.Region, "")
	if err != nil {
		return awsErr(err)
	}

	job.EndpointInfo = endpoint

	return nil
}

func awsErr(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		// AWS does not return a 404 for this, we have to check the error message
		if strings.HasPrefix(awsErr.Message(), "Could not find endpoint") ||
			strings.HasPrefix(awsErr.Message(), "Requested resource not found") {
			return nil
		}
	}
	return err
}
