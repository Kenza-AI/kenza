import React from 'react'
import { header } from './JobsHeader.jsx'
import { jobsList } from './JobsList.jsx'
import { jobsComparison } from './JobsComparison.jsx'

const MAX_SELECTED_JOBS = 3

export const jobsSection = (jobs, filterChangeHandler, viewJobDetailsHandler, viewLogsHandler, filtersChangeHandler, selectJobForComparisonHandler, compareJobsHandler, dismissJobsComparisonHandler, deleteJobsHandler, cancelJobsHandler) => (
    <section id="main-content">
        {header(jobs.list.length > 0 ? "> " + jobs.list[0].project.name : '', jobs.textFilter, filterChangeHandler, filtersChangeHandler, jobs.typeFilter)}
        {jobsList(jobs, viewJobDetailsHandler, viewLogsHandler, selectJobForComparisonHandler, MAX_SELECTED_JOBS, deleteJobsHandler, cancelJobsHandler)}
        {jobsComparison(jobs.inComparisonMode, jobs.isFetchingJobDetailsForComparison, jobs.fetchJobDetailsForComparisonError, jobs.list.filter(job => job.isSelectedForComparison), dismissJobsComparisonHandler, compareJobsHandler, MAX_SELECTED_JOBS)}
    </section>
)
