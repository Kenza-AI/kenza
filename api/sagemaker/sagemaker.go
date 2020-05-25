package sagemaker

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
)

// Client for AWS SageMaker
type Client struct {
	*sagemaker.SageMaker
}

// NewClient - SageMaker client initializer
func NewClient(profile, region, sessionName string, sharedConfigFiles []string) (*Client, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
		SharedConfigFiles: sharedConfigFiles,
	})
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		sagemaker.New(sess),
	}, nil
}

// GetTrainingJob returns the current status of the provided job.
func (sm *Client) GetTrainingJob(jobID, region, iamRole string) (*sagemaker.DescribeTrainingJobOutput, error) {
	log.Print(sm.Config.Region, region)
	sm.Config.Region = &region
	log.Print(sm.Config.Region, region)

	jobInput := &sagemaker.DescribeTrainingJobInput{
		TrainingJobName: aws.String(jobID),
	}

	job, err := sm.DescribeTrainingJob(jobInput)
	if err != nil {
		return &sagemaker.DescribeTrainingJobOutput{}, err
	}

	return job, nil
}

// GetTransformJob returns the job description as provided by AWS.
func (sm *Client) GetTransformJob(jobID, region, iamRole string) (*sagemaker.DescribeTransformJobOutput, error) {
	jobInput := &sagemaker.DescribeTransformJobInput{
		TransformJobName: aws.String(jobID),
	}

	job, err := sm.DescribeTransformJob(jobInput)
	if err != nil {
		return &sagemaker.DescribeTransformJobOutput{}, err
	}

	return job, nil
}

// GetTuningJob returns the job description as provided by AWS.
func (sm *Client) GetTuningJob(jobID, region, iamRole string) (*sagemaker.DescribeHyperParameterTuningJobOutput, error) {
	jobInput := &sagemaker.DescribeHyperParameterTuningJobInput{
		HyperParameterTuningJobName: aws.String(jobID),
	}

	job, err := sm.DescribeHyperParameterTuningJob(jobInput)
	if err != nil {
		return &sagemaker.DescribeHyperParameterTuningJobOutput{}, err
	}

	return job, nil
}

// GetEndpoint returns the specified endpoint.
func (sm *Client) GetEndpoint(name, region, iamRole string) (*sagemaker.DescribeEndpointOutput, error) {
	endpointInput := &sagemaker.DescribeEndpointInput{
		EndpointName: aws.String(name),
	}

	endpoint, err := sm.DescribeEndpoint(endpointInput)
	if err != nil {
		return &sagemaker.DescribeEndpointOutput{}, err
	}

	return endpoint, nil
}
