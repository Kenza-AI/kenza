import React from 'react'

export const metricsSection = (job) => (
    <section id={`results-${job.id}`} className='job-results'>
        <h4>Metrics</h4>
        {hasMetrics(job) ? metricsList(job.sagemaker_training_job.metrics) : <p>No metrics recorded</p>}
    </section>
)

const hasMetrics = job => (
    job.sagemaker_training_job.metrics != null && job.sagemaker_training_job.metrics.length > 0
)

const metricsList = metrics => (
    <ul>
        {metrics.map(metric =>
            <li key={metric.name}>
                {metric.name}: {metric.value}
            </li>)
        }
    </ul>
)

