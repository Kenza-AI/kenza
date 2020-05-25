import axios from 'axios'
import { signOutIfUnauthorized, apiBaseURL } from '../httputil'
import { 
  fetchJobs as fetchJobsAction,
  fetchJobsSuccess,
  fetchJobsFailure,
  fetchJobDetails as fetchJobDetailsAction,
  fetchJobDetailsSuccess,
  fetchJobDetailsFailure,
  fetchLogs as fetchLogsAction,
  compareJobs as compareJobsAction,
  compareJobsSuccess,
  compareJobsFailure,
  deleteJobs as deleteJobsAction,
  deleteJobsSuccess,
  deleteJobsFailure,
  cancelJobs as cancelJobsAction,
  cancelJobsSuccess,
  cancelJobsFailure
} from './actions'

const headers = { 'Content-Type': 'application/json', 'Accept': 'application/json' }

export const fetchJobs = async (dispatch, projectID, config, account) => {
  dispatch(fetchJobsAction())

  try {
    const response = await axios.get(jobListEndpoint(config, account, projectID),
      { headers: { ...headers, 'Authorization': account.account.accessToken } })
    dispatch(fetchJobsSuccess(response.data))
  } catch (err) {
    dispatch(fetchJobsFailure(err.message != undefined ? err.message : err))
    signOutIfUnauthorized(err, dispatch)
  }
}

export const compareJobs = async (dispatch, projectID, jobs, config, account) => {
  dispatch(compareJobsAction())

  const urls = []
  jobs.forEach(job => {
    urls.push(jobDetailsEndpoint(config, account, projectID, job.id))
  })

  try {
    const jobs = await Promise.all(
      urls.map(url => axios.get(url,
        { headers: { ...headers, 'Authorization': account.account.accessToken } }))
    )
    dispatch(compareJobsSuccess(jobs.map(job => job.data)))    
  } catch (err) {
    dispatch(compareJobsFailure("Please try again"))
    signOutIfUnauthorized(err, dispatch)
  }
}

export const fetchJobDetails = async (dispatch, projectID, jobID, config, account) => {
  dispatch(fetchJobDetailsAction(jobID))

  try {
    const response = await axios.get(jobDetailsEndpoint(config, account, projectID, jobID),
      { headers: { ...headers, 'Authorization': account.account.accessToken } })
    dispatch(fetchJobDetailsSuccess(response.data))
  } catch (err) {
    dispatch(fetchJobDetailsFailure(err.message != undefined ? err.message : err))
    signOutIfUnauthorized(err, dispatch)
  }
}

export const fetchLogs = async (dispatch, projectID, jobID, config, account) => {
  dispatch(fetchLogsAction())

  try {
    const response = await axios.get(jobLogsEndpoint(config, account, projectID, jobID),
      { headers: { ...headers, 'Authorization': account.account.accessToken } })
      logsResponseAsDownload(response, jobID)
    } catch (err) {
      if (signOutIfUnauthorized(err, dispatch)) {
        window.alert("Could not locate log file for job " + jobID)
      }
  }
}

export const cancelJobs = async (dispatch, projectID, jobIDs, config, account) => {
  dispatch(cancelJobsAction(jobIDs))

  let url = new URL(jobsCancellationEndpoint(config, account, projectID))
  jobIDs.forEach(id => url.searchParams.append('id', id)) 

  try {
    await axios.post(url, { headers: { ...headers, 'Authorization': account.account.accessToken } })
    dispatch(cancelJobsSuccess())
    fetchJobs(dispatch, projectID, config, account)
  } catch (err) {
    dispatch(cancelJobsFailure(err.message != undefined ? err.message : err))
    signOutIfUnauthorized(err, dispatch)
  }
}

export const deleteJobs = async (dispatch, projectID, jobIDs, config, account) => {
  dispatch(deleteJobsAction(jobIDs))

  let url = new URL(jobListEndpoint(config, account, projectID))
  jobIDs.forEach(id => url.searchParams.append('id', id)) 

  try {
    await axios.delete(url, { headers: { ...headers, 'Authorization': account.account.accessToken } })
    dispatch(deleteJobsSuccess())
    fetchJobs(dispatch, projectID, config, account)
  } catch (err) {
    dispatch(deleteJobsFailure(err.message != undefined ? err.message : err))
    signOutIfUnauthorized(err, dispatch)
  }
}

const jobsCancellationEndpoint = (config, account, projectID) => (
  apiBaseURL(config) + '/' + 'accounts' + '/' + account.account.accountID + '/' + 'projects' + '/' + projectID + '/' + 'jobs' + '/cancellations'
)

const jobListEndpoint = (config, account, projectID) => (
  apiBaseURL(config) + '/' + 'accounts' + '/' + account.account.accountID + '/' + 'projects' + '/' + projectID + '/' + 'jobs'
)

const jobDetailsEndpoint = (config, account, projectID, jobID) => (
  apiBaseURL(config) + '/' + 'accounts' + '/' + account.account.accountID + '/' + 'projects' + '/' + projectID + '/' + 'jobs' + '/' + jobID
)

const jobLogsEndpoint = (config, account, projectID, jobID) => (
  apiBaseURL(config) + '/' + 'accounts' + '/' + account.account.accountID + '/' + 'projects' + '/' + projectID + '/' + 'jobs' + '/' + jobID + '/' + 'logs'
)

const logsResponseAsDownload = (response, jobID) => {
  const url = window.URL.createObjectURL(new Blob([response.data]))
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', jobID + ".log")
  document.body.appendChild(link)
  link.click()
}