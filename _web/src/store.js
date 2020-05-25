import { connectRoutes } from 'redux-first-router'
const createHistory = require('history').createHashHistory
import { combineReducers, createStore, applyMiddleware, compose } from 'redux'

// Middleware
import thunk from 'redux-thunk'
import logger from 'redux-logger'

// Reducers
import { jobs } from './jobs/reducers'
import { account } from './auth/reducers'
import { projects, pendingProject } from './projects/reducers'

const history = createHistory()

// Routes.
const routesMap = {
  // Auth
  SIGNIN: '/signin',
  SIGNUP: '/signup',
  FORGOTPASSWORD: '/forgotpassword',

  // Projects / jobs
  PROJECTS: '/projects',
  PROJECT: '/projects/:projectID/jobs',
  JOB: '/projects/:projectID/jobs/:jobID',
}

const options = {
  initialDispatch: false
}

const { reducer, middleware, enhancer } = connectRoutes(history, routesMap, options)

const appReducer = combineReducers({
  location: reducer,
  jobs,
  account,
  projects,
  pendingProject
})

const middlewares = applyMiddleware(middleware, thunk, logger)
const store = createStore(appReducer, compose(enhancer, middlewares))

export default store