import { SIGN_OUT } from '../auth/actions'
import {
  FETCH_PROJECTS,
  FETCH_PROJECTS_FAILURE,
  FETCH_PROJECTS_SUCCESS,
  FILTER_PROJECTS,
  ADD_PROJECT_PROMPT,
  ADD_PROJECT_PROMPT_DISMISS,
  PENDING_PROJECT_CHANGED,
  CREATE_PROJECT,
  CREATE_PROJECT_SUCCESS,
  CREATE_PROJECT_FAILURE,
  DELETE_PROJECT,
  DELETE_PROJECT_SUCCESS,
  DELETE_PROJECT_FAILURE
} from './actions'

export const projects = (state = initialState(), action) => {
  switch (action.type) {
    case FETCH_PROJECTS:
    case DELETE_PROJECT:
      return {
        ...state,
        isFetching: true
      }
    case DELETE_PROJECT_SUCCESS:
      return { ...state, isFetching: false, deleteProjectError: action.error }
    case FETCH_PROJECTS_FAILURE:
      return {
        ...state,
        isFetching: false,
        hasFetchedOnce: true,
        fetchProjectsError: action.error
      }
    case DELETE_PROJECT_FAILURE:
      return {
        ...state,
        isFetching: false,
        deleteProjectError: action.error
      }
    case FETCH_PROJECTS_SUCCESS:
      return {
        ...state,
        isFetching: false,
        hasFetchedOnce: true,
        fetchProjectsError: null,
        list: action.projects
      }
    case FETCH_PROJECTS_SUCCESS:
      return {
        ...state,
        isFetching: false,
        deleteProjectError: null,
      }
    case FILTER_PROJECTS:
      return {
        ...state,
        textFilter: action.filterText
      }
    case SIGN_OUT:
      return initialState()
    default:
      return state
  }
}

export const pendingProject = (state = null, action) => {
  switch (action.type) {
    case ADD_PROJECT_PROMPT:
      return pendingNewProjectInitialState
    case ADD_PROJECT_PROMPT_DISMISS:
      return null
    case PENDING_PROJECT_CHANGED:
      return action.project
    case CREATE_PROJECT:
      return { ...state, isCreating: true }
    case CREATE_PROJECT_SUCCESS:
    case CREATE_PROJECT_FAILURE: // TODO(ilazakis): handle error (display error on UI)
      return { ...state, isCreating: false }
    case SIGN_OUT:
      return null
    default:
      return state
  }
}

const initialState = () => ({
  list: [],
  textFilter: "",
  isFetching: false,
  deleteProjectError: null,
  fetchProjectsError: null,
  hasFetchedOnce: false
})

const pendingNewProjectInitialState = {
  title: '',
  description: '',
  repo: '',
  branch: '',
  accessToken: '',
  isCreating: false
}