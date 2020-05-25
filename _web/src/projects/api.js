import axios from 'axios'
import { signOutIfUnauthorized, apiBaseURL } from '../httputil'
import { 
  fetchProjects as fetchProjectsAction,
  fetchProjectsSuccess,
  fetchProjectsFailure,
  createProject as createProjectAction,
  createProjectSuccess,
  createProjectFailure,
  hideAddProject,
  createJob as createJobAction,
  deleteProject as deleteProjectAction,
  deleteProjectSuccess,
  deleteProjectFailure } from './actions'

const headers = { 'Content-Type': 'application/json', 'Accept': 'application/json' }

export const fetchProjects = async (dispatch, config, account) => {
  dispatch(fetchProjectsAction())

  try {
    const response = await axios.get(projectsEndpoint(config, account),
      { headers: { ...headers, 'Authorization': account.account.accessToken } })
    dispatch(fetchProjectsSuccess(response.data))
  } catch (err) {
    dispatch(fetchProjectsFailure(err.message != undefined ? err.message : err))
    signOutIfUnauthorized(err, dispatch)
  }
}

export const createProject = async (dispatch, project, config, account) => {
  const sanitizedDescription = project.description === undefined || project.description === null || project.description.trim().length == 0 ? 'No Description' : project.description
  dispatch(createProjectAction(project))

  const data = { title: project.title, description: sanitizedDescription, repo: project.repo, branch: project.branch, accessToken: project.accessToken }
  try {
    await axios.post(projectsEndpoint(config, account), data, { headers: { ...headers, 'Authorization': account.account.accessToken } })
    dispatch(createProjectSuccess())
    dispatch(hideAddProject())
    fetchProjects(dispatch, config, account)
  } catch (err) {
    dispatch(createProjectFailure(err.message != undefined ? err.message : err))
    signOutIfUnauthorized(err, dispatch)
  }
}

export const deleteProject = async (dispatch, projectID, config, account) => {
  dispatch(deleteProjectAction(projectID))

  try {
    await axios.delete(deleteProjectEndpoint(config, account, projectID), { headers: { ...headers, 'Authorization': account.account.accessToken } })
    dispatch(deleteProjectSuccess())
    fetchProjects(dispatch, config, account)
  } catch (err) {
    dispatch(deleteProjectFailure(err.message != undefined ? err.message : err))
    signOutIfUnauthorized(err, dispatch)
  }
}

export const submitJob = async (dispatch, project, config, account) => {
  dispatch(createJobAction(project))

  // https://developer.github.com/v3/repos/hooks/#create-a-hook
  const data = {
    ref: project.branch, 
    sender: {
      login: account.account.username
    },
    repository: {
      clone_url: project.repo
    },
    head_commit: {
      id: "" // will build HEAD
    }
  }

  try {
    await axios.post(submitJobEndpoint(config, account, project.projectID), data, { headers: { ...headers, 'Authorization': account.account.accessToken } })
  } catch (err) {
    signOutIfUnauthorized(err, dispatch)
  }
}

const projectsEndpoint = (config, account) => (
  apiBaseURL(config) + '/' + 'accounts' + '/' + account.account.accountID + '/' + 'projects'
)

const submitJobEndpoint = (config, account, projectID) => (
  apiBaseURL(config) + '/' + 'accounts' + '/' + account.account.accountID + '/' + 'projects' + '/' + projectID + '/' + 'jobs' + '/' + 'submissions'
)

const deleteProjectEndpoint = (config, account, projectID) => (
  apiBaseURL(config) + '/' + 'accounts' + '/' + account.account.accountID + '/' + 'projects' + '/' + projectID
)