import React from 'react'
import Error from '../Error.jsx'
import { metricsSection } from './JobMetrics.jsx'
import { parametersSection } from './JobHyperparameters.jsx'
import LoadingIndicator from '../LoadingIndicator.jsx'
import TuningJobDetails from './TuningJobDetails.jsx'

export default class JobDetails extends React.Component {

    // Component lifecycle
    componentDidMount() {
        this.props.fetchJobDetails(this.props.job.project.id, this.props.job.id)
    }

    componentDidUpdate(prevProps) {
        if (this.props.job.id != prevProps.job.id) {
            this.props.fetchJobDetails(this.props.job.project.id, this.props.job.id)
        }
    }

    render() {
        const { job, fetchJobDetailsError, isFetchingJobDetails } = this.props

        // Display any errors if present.
        if (fetchJobDetailsError != null) {
            return <Error message={fetchJobDetailsError} />
        }
        if (isFetchingJobDetails) {
            return loadingJobDetailsIndicator()
        }
        return jobDetails(job)
    }
}

const jobDetails = job => (
    <section id="job-details">
        {job.type == "tuning" &&
            job.sagemaker_tuning_job.DescribeHyperParameterTuningJobOutput != null &&
            TuningJobDetails(job.id, job.sagemaker_tuning_job.DescribeHyperParameterTuningJobOutput)}


        <section id="job-details-basic">
            <h4>Job id: {job.id}</h4>
            <ul>
                <li><b>Job Type: </b> {job.type}</li>
                {job.type == "training" && <li><b>Training Time: </b> {job.sagemaker_training_job.training_time_seconds} seconds</li>}
                {job.type == "training" && s3ModelLocation(job)}
                {job.type == "batch" || job.type == "tuning" && s3FeaturesInputLocation(job)}
                {job.type == "batch" || job.type == "tuning" && s3FeaturesOutputLocation(job)}
            </ul>
        </section>

        {job.type == "training" && resourceConfig(job)}
        {job.endpoint_info != undefined && endpointSection(job.endpoint_info, job.region)}
        {job.type == "training" && metricsSection(job)}
        {job.type == "training" && parametersSection(job)}
    </section>
)


const resourceConfig = job => (
    <section id="job-details-basic">
        <h4>Training configuration</h4>
        <ul>
            <li><b>Instance type: </b> {job.sagemaker_training_job.resource_config.instance_type}</li>
            <li><b>Instance count: </b> {job.sagemaker_training_job.resource_config.instance_count}</li>
            <li><b>Volume size: </b> {job.sagemaker_training_job.resource_config.volume_size_gb}</li>
        </ul>
    </section>
)

const s3ModelLocation = job => (
    <li><b>S3 model location:</b> {job.sagemaker_training_job.s3_model_location}</li>
)

const s3FeaturesInputLocation = job => (
    <li>S3 features input location: {job.type == 'batch' ? job.sagemaker_transform_job.s3_input_location : job.sagemaker_tuning_job.s3_input_location}</li>
)

const s3FeaturesOutputLocation = job => (
    <li>S3 features output location: {job.type == 'batch' ? job.sagemaker_transform_job.s3_output_location : job.sagemaker_tuning_job.s3_output_location}</li>
)

const loadingJobDetailsIndicator = () => (
    <section id="main-content-loading">
        <h2>Loading training job details ...</h2>
        <LoadingIndicator />
    </section>
)

const endpointSection = (endpoint, region) => (
    <section id="job-details-endpoint">
        <h4>Endpoint info</h4>
        <ul>
            {<li className={`${endpoint.EndpointStatus == 'InService' ? 'endpoint-up' : null}`}><b>Status:</b> {endpoint.EndpointStatus}</li>}
            {<li><b>Created: </b>{new Date(endpoint.CreationTime).toLocaleString()}</li>}
            {<li><b>Last updated: </b>{new Date(endpoint.LastModifiedTime).toLocaleString()}</li>}
            {<li><b>Endpoint Name: </b>{endpoint.EndpointName}</li>}
            {<li><b>Endpoint Config Name: </b>{endpoint.EndpointConfigName}</li>}
            {<li><b>Endpoint URL: </b>{`https://runtime.sagemaker.${region}.amazonaws.com/endpoints/${endpoint.EndpointName}/invocations`}</li>}
        </ul>
        {/* {endpointInvocationSection(endpoint, region)} */}
    </section>
)

const endpointInvocationSection = (endpoint, region) => (
    <section id="endpoint-invocation">
        <h4>Endpoint invocation</h4>

        <section className='endpoint-invocation'>
            <section>
                <label for="endpoint-invocation-payload" className='endpoint-invocation-label'>Request:</label>
                <textarea id="endpoint-invocation-payload">Place your JSON payload here.</textarea>
            </section>
            <button id='endpoint-invocation-cta'>Call endpoint</button>

            <section>
                <label for="endpoint-invocation-response" className='endpoint-invocation-label'>Response:</label>
                <textarea id="endpoint-invocation-response"></textarea>
            </section>
        </section>
    </section>
)

