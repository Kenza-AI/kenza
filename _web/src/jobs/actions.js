export const FETCH_JOB_DETAILS = 'FETCH_JOB_DETAILS'
export const FETCH_JOB_DETAILS_SUCCESS = 'FETCH_JOB_DETAILS_SUCCESS'
export const FETCH_JOB_DETAILS_FAILURE = 'FETCH_JOB_DETAILS_FAILURE'
export const FETCH_JOBS = 'FETCH_JOBS'
export const FETCH_JOBS_SUCCESS = 'FETCH_JOBS_SUCCESS'
export const FETCH_JOBS_FAILURE = 'FETCH_JOBS_FAILURE'
export const FILTER_JOBS = 'FILTER_JOBS'
export const FILTER_JOBS_TYPE = 'FILTER_JOBS_TYPE'
export const VIEW_JOB_DETAILS = 'VIEW_JOB_DETAILS'
export const FETCH_LOGS = 'FETCH_LOGS'
export const SELECT_JOB_FOR_COMPARISON = 'SELECT_JOB_FOR_COMPARISON'
export const CLEAR_JOBS_SELECTED_FOR_COMPARISON = 'CLEAR_JOBS_SELECTED_FOR_COMPARISON'
export const DISMISS_JOBS_COMPARISON = 'DISMISS_JOBS_COMPARISON'
export const COMPARE_JOBS = 'COMPARE_JOBS'
export const COMPARE_JOBS_SUCCESS = 'COMPARE_JOBS_SUCCESS'
export const COMPARE_JOBS_FAILURE = 'COMPARE_JOBS_FAILURE'
export const DELETE_JOBS = 'DELETE_JOBS'
export const DELETE_JOBS_SUCCESS = 'DELETE_JOBS_SUCCESS'
export const DELETE_JOBS_FAILURE = 'DELETE_JOBS_FAILURE'
export const CANCEL_JOBS = 'CANCEL_JOBS'
export const CANCEL_JOBS_SUCCESS = 'CANCEL_JOBS_SUCCESS'
export const CANCEL_JOBS_FAILURE = 'CANCEL_JOBS_FAILURE'

export const fetchJobDetailsSuccess = job => ({
  type: FETCH_JOB_DETAILS_SUCCESS,
  job: job
})

export const fetchJobDetails = jobID => ({
  type: FETCH_JOB_DETAILS,
  jobID: jobID
})

export const fetchJobDetailsFailure = error => ({
  type: FETCH_JOB_DETAILS_FAILURE,
  error: error
})

export const fetchJobsSuccess = jobs => ({
  type: FETCH_JOBS_SUCCESS,
  jobs: jobs
})

export const fetchJobs = () => ({
  type: FETCH_JOBS
})

export const fetchJobsFailure = error => ({
  type: FETCH_JOBS_FAILURE,
  error: error
})

export const filterJobs = text => ({
  type: FILTER_JOBS,
  text: text
})

export const changeFilters = filter => ({
  type: FILTER_JOBS_TYPE,
  filter: filter
})

export const fetchLogs = () => ({
  type: FETCH_LOGS
})

export const selectJobForComparison = id => ({
  type: SELECT_JOB_FOR_COMPARISON,
  id: id
})

export const clearJobsSelectedForComparison = () => ({
  type: CLEAR_JOBS_SELECTED_FOR_COMPARISON
})

export const dismissJobsComparison = () => ({
  type: DISMISS_JOBS_COMPARISON
})

export const compareJobs = () => ({
  type: COMPARE_JOBS
})

export const compareJobsSuccess = jobs => ({
  type: COMPARE_JOBS_SUCCESS,
  jobs: jobs
})

export const compareJobsFailure = error => ({
  type: COMPARE_JOBS_FAILURE,
  error: error
})

export const deleteJobs = jobs => ({
  type: DELETE_JOBS,
  jobs: jobs
})

export const deleteJobsSuccess = () => ({
  type: DELETE_JOBS_SUCCESS
})

export const deleteJobsFailure = error => ({
  type: DELETE_JOBS_FAILURE,
  error: error
})

export const cancelJobs = jobs => ({
  type: CANCEL_JOBS,
  jobs: jobs
})

export const cancelJobsSuccess = () => ({
  type: CANCEL_JOBS_SUCCESS
})

export const cancelJobsFailure = error => ({
  type: CANCEL_JOBS_FAILURE,
  error: error
})
