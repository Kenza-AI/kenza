import React from 'react'

export const parametersSection = (job) => (
    <section className='job-algorithms'>
        <h4>Parameters</h4>
        {hasParameters(job) ? parametersList(job.sagemaker_training_job.hyperparameters) : <p>No parameters recorded</p>}
    </section>
)

const hasParameters = job => (
    job.sagemaker_training_job.hyperparameters != null && Object.keys(job.sagemaker_training_job.hyperparameters).length > 0
)

const parametersList = params => (
    <ul>
        {Object.keys(params).map(param =>
            <li key={param}>
                {param}: {params[param]}
            </li>)
        }
    </ul>
)

