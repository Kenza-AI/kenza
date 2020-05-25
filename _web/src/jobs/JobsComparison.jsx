import React from 'react'
import { closeIcon } from '../icons'
import Error from '../Error.jsx'
import JobSummary from './JobSummary.jsx'
import { loadingJobs } from './JobsLoading.jsx'

export const jobsComparison = (inComparisonMode, isFetchingJobDetailsForComparison, fetchJobDetailsForComparisonError, jobsSelectedForComparison, dismissJobsComparisonHandler, compareJobsHandler, maxAllowedJobsForComparison) => {
    if (inComparisonMode && !fetchJobDetailsForComparisonError) {
        return jobsComparisonPopUp(jobsSelectedForComparison, dismissJobsComparisonHandler, isFetchingJobDetailsForComparison)
    } else {
        return jobsSelectedForComparison.length > 0 ?
            jobsPreComparisonSection(jobsSelectedForComparison, fetchJobDetailsForComparisonError, compareJobsHandler, maxAllowedJobsForComparison) : null
    }
}

const jobsPreComparisonSection = (jobsSelectedForComparison, fetchJobDetailsForComparisonError, compareJobsHandler, maxAllowedJobsForComparison) => (
    <section className='jobs-selection'>
        <p>
            <b>{`${jobsSelectedForComparison.length} job${jobsSelectedForComparison.length > 1 ? 's' : ''}`} selected </b> ·
            <i>{`Select up to ${maxAllowedJobsForComparison} jobs to compare`} · <b>esc</b> to deselect all</i>
        </p>
        <div>
            {fetchJobDetailsForComparisonError && <b className='jobs-comparison-error'>{fetchJobDetailsForComparisonError}</b>}
            <button className={`${jobsSelectedForComparison.length <= 1 ? 'disabled' : null}`} onClick={() => compareJobsHandler(jobsSelectedForComparison)}>
                Compare jobs
            </button>
        </div>
    </section>
)

const jobsComparisonPopUp = (jobs, dismissJobsComparisonHandler, isFetchingJobDetailsForComparison) => (
    <section className='fullscreen-overlay' onClick={() => dismissJobsComparisonHandler()}>
        <ul className={`fullscreen-overlay-content jobs-comparison ${isFetchingJobDetailsForComparison ? 'loading' : null}`} onClick={e => e.stopPropagation()}>
            <h2>Comparing {jobs.length} Jobs </h2>

            {isFetchingJobDetailsForComparison ? loadingJobs() :
                jobs.map(job => <JobSummary
                    job={job}
                    key={job.id}
                    inComparisonMode={true}
                    viewLogsHandler={noop}
                    viewJobDetailsHandler={noop}
                    selectForComparisonHandler={noop} />)}

            <button className='dismiss-modal-job-comparison' onClick={() => dismissJobsComparisonHandler()}>{closeIcon()}</button>
        </ul>
    </section>
)

const noop = () => { }
