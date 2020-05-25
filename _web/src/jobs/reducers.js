import {
  FETCH_JOBS,
  FETCH_JOBS_SUCCESS,
  FETCH_JOBS_FAILURE,
  FILTER_JOBS,
  FILTER_JOBS_TYPE,
  FETCH_JOB_DETAILS,
  FETCH_JOB_DETAILS_FAILURE,
  FETCH_JOB_DETAILS_SUCCESS,
  SELECT_JOB_FOR_COMPARISON,
  CLEAR_JOBS_SELECTED_FOR_COMPARISON,
  COMPARE_JOBS,
  DISMISS_JOBS_COMPARISON,
  COMPARE_JOBS_FAILURE,
  COMPARE_JOBS_SUCCESS,
  DELETE_JOBS,
  DELETE_JOBS_SUCCESS,
  DELETE_JOBS_FAILURE,
  CANCEL_JOBS,
  CANCEL_JOBS_SUCCESS,
  CANCEL_JOBS_FAILURE
} from './actions'
import { SIGN_OUT } from '../auth/actions'

export function jobs(state = initialState, action) {
  switch (action.type) {
    case FETCH_JOBS_SUCCESS:
      return {
        ...state,
        isFetching: false,
        list: diffExistingJobsWithIncomingJobs(state.list, action.jobs),
        fetchJobsError: null
      }
    case FETCH_JOBS:
      return { ...state, isFetching: true }
    case SIGN_OUT:
      return initialState
    case FETCH_JOBS_FAILURE:
      return {
        ...state,
        isFetching: false,
        fetchJobsError: action.error
      }
    case FILTER_JOBS:
      return {
        ...state,
        textFilter: action.text
      }
    case FILTER_JOBS_TYPE:
      return {
        ...state,
        typeFilter: action.filter
      }
    case FETCH_JOB_DETAILS:
      return {
        ...state,
        isFetchingJobDetails: true
      }
    case FETCH_JOB_DETAILS_SUCCESS:
      return {
        ...state,
        list: state.list.length > 0 ? state.list.map(job => job.id == action.job.id ? action.job : job) : [action.job],
        isFetchingJobDetails: false,
        fetchJobDetailsError: null
      }
    case FETCH_JOB_DETAILS_FAILURE:
      return {
        ...state,
        isFetchingJobDetails: false,
        fetchJobDetailsError: action.error
      }
    case SELECT_JOB_FOR_COMPARISON:
      return {
        ...state,
        list: state.list.map(job => (job.id != action.id ? job :
          { ...job, isSelectedForComparison: !job.isSelectedForComparison }))
      }
    case CLEAR_JOBS_SELECTED_FOR_COMPARISON:
      return {
        ...state,
        fetchJobDetailsForComparisonError: null,
        list: state.list.map(job => ({ ...job, isSelectedForComparison: false }))
      }
    case COMPARE_JOBS:
      return { ...state, inComparisonMode: true, isFetchingJobDetailsForComparison: true }
    case COMPARE_JOBS_FAILURE:
      return {
        ...state,
        inComparisonMode: false,
        isFetchingJobDetailsForComparison: false,
        fetchJobDetailsForComparisonError: action.error,
      }
    case COMPARE_JOBS_SUCCESS:
      return {
        ...state,
        list: updateDetailsForJobsBeingCompared(state.list, action.jobs),
        isFetchingJobDetailsForComparison: false,
        fetchJobDetailsForComparisonError: null
      }
    case DISMISS_JOBS_COMPARISON:
      return { ...state, inComparisonMode: false }
    case DELETE_JOBS:
      return { ...state, isDeleting: true }
    case DELETE_JOBS_SUCCESS:
      return { ...state, isDeleting: false, deleteJobsError: null }
    case DELETE_JOBS_FAILURE:
      return { ...state, isDeleting: false, deleteJobsError: action.error }
    case CANCEL_JOBS:
      return { ...state, isCancelling: true }
    case CANCEL_JOBS_SUCCESS:
      return { ...state, isCancelling: false, cancelJobsError: null }
    case CANCEL_JOBS_FAILURE:
      return { ...state, isCancelling: false, cancelJobsError: action.error }
    default:
      return state
  }
}

const initialState = {
  list: [],
  typeFilter: [],
  textFilter: "",
  inComparisonMode: false,

  // Fetch all jobs
  isFetching: false,
  fetchJobsError: null,

  // Fetch details for one job
  fetchJobDetailsError: null,
  isFetchingJobDetails: false,

  // Fetch details of jobs selected for comparison
  fetchJobDetailsForComparisonError: null,
  isFetchingJobDetailsForComparison: false,

  // Delete jobs
  isDeleting: false,
  deleteJobsError: null,

  // Cancel jobs
  isCancelling: false,
  cancelJobsError: null
}

// Returns the existing list of jobs but with the jobs being compared updated with their details we just fetched 
const updateDetailsForJobsBeingCompared = (jobs, jobsForComparison) => {
  const jobsForComparisonMap = jobsForComparison.reduce((map, job) => { map[job.id] = job; return map }, {})
  return jobs.map(job => (jobsForComparisonMap[job.id] != undefined ? { ...jobsForComparisonMap[job.id], isSelectedForComparison: true } : job))
}

// Allows existing local UI state (e.g. selected for comparison) to persist when the list of jobs is refreshed.
// If not, one would select a job for comparison and as soon as the list is refreshed (from auto-refresh / polling)
// the job would be automatically unselected, an annoying experience.
const diffExistingJobsWithIncomingJobs = (existing, incoming) => {
  return incoming.map(job => (
    {
      ...job, isSelectedForComparison: existing.find(j => j.id == job.id) != undefined
        ? existing.find(j => j.id == job.id).isSelectedForComparison : false
    }))
}
