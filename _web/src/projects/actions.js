export const DELETE_PROJECT = 'DELETE_PROJECT'
export const DELETE_PROJECT_SUCCESS = 'DELETE_PROJECT_SUCCESS'
export const DELETE_PROJECT_FAILURE = 'DELETE_PROJECT_FAILURE'
export const CREATE_PROJECT = 'CREATE_PROJECT'
export const CREATE_PROJECT_SUCCESS = 'CREATE_PROJECT_SUCCESS'
export const CREATE_PROJECT_FAILURE = 'CREATE_PROJECT_FAILURE'
export const FILTER_PROJECTS = 'FILTER_PROJECTS'
export const FETCH_PROJECTS = 'FETCH_PROJECTS'
export const FETCH_PROJECTS_SUCCESS = 'FETCH_PROJECTS_SUCCESS'
export const FETCH_PROJECTS_FAILURE = 'FETCH_PROJECTS_FAILURE'
export const ADD_PROJECT_PROMPT = 'ADD_PROJECT_PROMPT'
export const ADD_PROJECT_PROMPT_DISMISS = 'ADD_PROJECT_PROMPT_DISMISS'
export const PENDING_PROJECT_CHANGED = 'PENDING_PROJECT_CHANGED'
export const CREATE_JOB = 'CREATE_JOB'

export const createProject = (project) => ({
  type: CREATE_PROJECT,
  project: project
})

export const createProjectSuccess = () => ({
  type: CREATE_PROJECT_SUCCESS
})

export const createProjectFailure = error => ({
  type: CREATE_PROJECT_FAILURE,
  error: error
})

export const fetchProjects = () => ({
  type: FETCH_PROJECTS
})

export const fetchProjectsSuccess = projects => ({
  type: FETCH_PROJECTS_SUCCESS,
  projects: projects
})

export const fetchProjectsFailure = error => ({
  type: FETCH_PROJECTS_FAILURE,
  error: error
})

export const filterProjects = text => ({
  type: FILTER_PROJECTS,
  filterText: text
})

export const showAddProject = () => ({
  type: ADD_PROJECT_PROMPT
})

export const hideAddProject = () => ({
  type: ADD_PROJECT_PROMPT_DISMISS
})

export const pendingProjectChanged = project => ({
  type: PENDING_PROJECT_CHANGED,
  project: project
})

export const createJob = project => ({
  type: CREATE_JOB,
  project: project
})

export const deleteProject = (projectID) => ({
  type: DELETE_PROJECT,
  projectID: projectID
})

export const deleteProjectSuccess = () => ({
  type: DELETE_PROJECT_SUCCESS
})

export const deleteProjectFailure = error => ({
  type: DELETE_PROJECT_FAILURE,
  error: error
})
