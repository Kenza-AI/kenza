import React from 'react'

const TuningJobDetails = (jobID, job) => (
    <div>
        {basicJobDetailsSection(jobID, job)}

        {bestTrainingJobSection(job.BestTrainingJob)}

        {tuningConfigSection(job.HyperParameterTuningJobConfig)}
        
        {trainingJobStatusCountersSection(job.TrainingJobStatusCounters)}

        {objectiveStatusCountersSection(job.ObjectiveStatusCounters)}
    </div>
)

const basicJobDetailsSection = (jobID, job) => (
    <section id="job-details-basic">
        <h4>Job id: {jobID}</h4>
        <ul>
            <li><b>Job Type: </b> Hyperparameter Tuning</li>
            {objectToList(job)}
        </ul>
    </section>
)

const bestTrainingJobSection = job => (
    <section id="job-details-basic">
        <h4>Best Training Job</h4>
        <ul>
            {objectToList(job)}

            {listHeader("Tuned Hyper Parameters")}
            <ul>{objectToList(job.TunedHyperParameters)}</ul>

            {listHeader("Final HyperParameter Tuning Job Objective Metric")}
            <ul>{objectToList(job.FinalHyperParameterTuningJobObjectiveMetric)}</ul>
        </ul>
    </section>
)

const tuningConfigSection = job => (
    <section id="job-details-basic">
        <h4>Tuning Job Configuration</h4>
        <ul>
            {objectToList(job)}

            {listHeader("Objective")}
            <ul>{objectToList(job.HyperParameterTuningJobObjective)}</ul>

            {listHeader("Resource Limits")}
            <ul>{objectToList(job.ResourceLimits)}</ul>

            {listHeader("Parameter Ranges")}
            {listHeader("Categorical")}
            {parameterRangeToList(job.ParameterRanges.CategoricalParameterRanges)}

            {listHeader("Integer")}
            {parameterRangeToList(job.ParameterRanges.IntegerParameterRanges)}

            {listHeader("Continuous")}
            {parameterRangeToList(job.ParameterRanges.ContinuousParameterRanges)}
        </ul>
    </section>
)

const objectiveStatusCountersSection = counters => (
    <section id="job-details-basic">
        <h4>Objective Status Counters</h4>
        <ul>
            {objectToList(counters)}
        </ul>
    </section>
)

const trainingJobStatusCountersSection = counters => (
    <section id="job-details-basic">
        <h4>Training Job Status Counters</h4>
        <ul>
            {objectToList(counters)}
        </ul>
    </section>
)


const listHeader = title => (
    <li><b>{title}</b></li>
)

const objectToList = object => {
    let list = []
    for (var key in object) {
        list.push(typeof (object[key]) == "object" ? null : <li key={key}><b>{key}: </b>{object[key]}</li>)
    }
    return list
}

const parameterRangeToList = range => (
    <ul>
        {range.map(param =>
            <li key={param["Name"]}>
                <ul>
                    {Object.keys(param).map(p =>
                        <li key={p}>
                            <b>{p}: </b> {param[p]}
                        </li>)}
                </ul>
            </li>)}
    </ul>
)

export default TuningJobDetails