import React from 'react'
import PropTypes from 'prop-types'
import { metricsSection } from './JobMetrics.jsx'
import { parametersSection } from './JobHyperparameters.jsx'

import { selectedIcon, unselectedIcon } from '../icons'

const JobSummary = ({ job, viewJobDetailsHandler, viewLogsHandler, inComparisonMode, selectForComparisonHandler, deleteJobsHandler, cancelJobsHandler }) => (
    <li id={job.id} className={className(job.codeLocation)}>
        {jobComparisonSection(job, inComparisonMode, selectForComparisonHandler)}
        <button className='job-id' onClick={() => viewJobDetailsHandler(job)}>id: {job.id}</button>
        <h3 className='job-title collapsed'>{sanitizedJobType(job)}</h3>
        <i className={`job-status ${job.status.toLowerCase()}`}>{job.status}</i>
        <i className={`job-running-time ${job.status.toLowerCase()}`}>{jobTimerSection(job)}</i>
        {!inComparisonMode && <i className='job-created-by'>Submitted by {job.submitter} on {new Date(job.created).toLocaleString()}</i>}
        {!inComparisonMode && jobOptionsSection(job, viewLogsHandler, deleteJobsHandler, viewJobDetailsHandler, cancelJobsHandler)}
        {inComparisonMode && job.isSelectedForComparison ? metricsSection(job) : null}
        {inComparisonMode && job.isSelectedForComparison ? parametersSection(job) : null}
    </li>
)

const className = commitID => {
    let jobClass = 'training-job'
    if (commitID) jobClass += ' no-scm-info collapsed'
    return jobClass
}

const jobTimerSection = job => {
    if (job.status == 'submitted') return null
    const endDateToUseForTimer = (job.status == 'running') ? new Date() : new Date(job.updated)
    const startDateToUseForTimer = new Date(job.started)

    const timeSinceJobStarted = new Date(endDateToUseForTimer - startDateToUseForTimer)
    const hoursValue = timeSinceJobStarted.getUTCHours()
    const minsValue = timeSinceJobStarted.getUTCMinutes()
    const secondsValue = timeSinceJobStarted.getUTCSeconds()

    const hours = hoursValue == 0 ? "" : 
        hoursValue < 10 ? '0' + hoursValue + "h:" :  hoursValue + "h:"
    
    const mins = minsValue < 10 ? '0' + minsValue + "m:" :  minsValue + "m:"
    
    const seconds = secondsValue < 10 ? '0' + secondsValue + "s" :  secondsValue + "s"

    return hours + mins + seconds
}

const jobOptionsSection = (job, viewLogsHandler, deleteJobsHandler, viewJobDetailsHandler, cancelJobsHandler) => (
    <div className='job-options'>
        {job.commitID != null && job.commitID != '' && codeSection(sanitizedRepoLocation(job.project.repo) + '/commit/' + job.commitID)}
        {isFinished(job) && viewJobDetailsCTA(job, viewJobDetailsHandler)}
        {isFinished(job) && logCTA(job.id, viewLogsHandler, job.project.id)}
        {isFinished(job) && deleteCTA(job.id, deleteJobsHandler, job.project.id)}
        
        {/* temporarily hide the "cancel job" CTA while the "cancel job" feature is under development*/}
        {/* {!isFinished(job) && cancelCTA(job.id, cancelJobsHandler, job.project.id)} */}
        
    </div>
)

const logCTA = (id, viewLogsHandler, projectID) => (
    <button onClick={() => viewLogsHandler(projectID, id)}>| View Build Log |</button>
)

const viewJobDetailsCTA = (job, viewJobDetailsHandler) => (
    <button onClick={() => viewJobDetailsHandler(job)}>| View Details |</button>
)

const cancelCTA = (id, cancelJobsHandler, projectID) => (
    <button className='job-cancellation' onClick={() => cancelJobsHandler(projectID, [id])}>Cancel</button>
)

const deleteCTA = (id, deleteJobsHandler, projectID) => (
    <button className='job-deletion' onClick={() => deleteJobsHandler(projectID, [id])}>Delete</button>
)

const codeSection = (codeLocation) => (
    <button className='job-code'>
        <a href={codeLocation} target='_blank'>| View Source Code |</a>
    </button>
)

const jobComparisonSection = (job, inComparisonMode, selectForComparisonHandler) => {
    if (inComparisonMode) {
        return <div className='job-selection'></div>
    } else {
        return <button className='job-selection' onClick={() => selectForComparisonHandler(job.id, job.isSelectable, job.isSelectedForComparison)}>
            {comparisonIcon(job.isSelectedForComparison, job.status.toUpperCase())}{comparisonText(job.status.toUpperCase())}
        </button>
    }
}

const comparisonText = status => (
    status == 'COMPLETED' || status == 'FAILED' ? 'Compare job' : null
)

const comparisonIcon = (isSelectedForComparison, status) => (
    status == 'COMPLETED' || status == 'FAILED' ? (isSelectedForComparison ? selectedIcon() : unselectedIcon()) : null
)

// Drops `.git` suffix from repo name if needed
const sanitizedRepoLocation = repo => (
    repo.endsWith(".git") ? repo.substring(0, repo.length - 4) : repo
)

const sanitizedJobType = job => (
    job.type == 'unknown' ? "" :  
        job.type == 'batchtransform' ? 'batch transform' : job.type
)

const isFinished = job => (
    job.status == 'completed' || job.status == 'failed'
)

export default JobSummary

JobSummary.propTypes = {
    job: PropTypes.object.isRequired,
    viewLogsHandler: PropTypes.func.isRequired,
    inComparisonMode: PropTypes.bool.isRequired,
    viewJobDetailsHandler: PropTypes.func.isRequired,
    selectForComparisonHandler: PropTypes.func.isRequired
}
