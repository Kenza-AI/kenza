import React from 'react'
import JobSummary from './JobSummary.jsx'

export const jobsList = (jobs, viewJobDetailsHandler, viewLogsHandler, selectForComparisonHandler, maxAllowedJobsForComparison, deleteJobsHandler, cancelJobsHandler) => (
        <ul className='card'>
            {jobs.list.filter(job => (job.commitID.toLowerCase().includes(jobs.textFilter.toLowerCase())))
                .filter(job => jobs.typeFilter.length > 0 ? jobs.typeFilter.includes(job.type) || (jobs.typeFilter.includes("endpoints") ? job.hasOwnProperty("endpoint") : false) : true)
                .sort((a, b) => { return new Date(b.created).getTime() - new Date(a.created).getTime() })
                .map(job => <JobSummary
                    key={job.id}
                    job={{...job, "isSelectable": isJobSelectable(jobs, job, maxAllowedJobsForComparison)}}
                    inComparisonMode={jobs.inComparisonMode}
                    viewLogsHandler={viewLogsHandler}
                    viewJobDetailsHandler={viewJobDetailsHandler}
                    selectForComparisonHandler={selectForComparisonHandler}
                    cancelJobsHandler={cancelJobsHandler}
                    deleteJobsHandler={deleteJobsHandler} />)}
        </ul>
)

const isJobSelectable = (jobs, job, maxAllowedJobsForComparison) => (
    job.isSelectedForComparison ||
    jobs.list.filter(job => job.isSelectedForComparison).length < maxAllowedJobsForComparison &&
    job.status.toUpperCase() == 'COMPLETED' || job.status.toUpperCase() == 'FAILED'
)

