import axios from 'axios'
import 'regenerator-runtime/runtime'
import { apiBaseURL } from '../httputil'
import { signIn as signInAction, signInSuccess, signInFailure, signUp as signUpAction, signUpSuccess, signUpFailure } from './actions'

const headers = { 'Content-Type': 'application/json' }

export const signIn = async (dispatch, config, email, password) => {
  dispatch(signInAction(email, password))

  const data = { email: email, password: password }
  try {
    const response = await axios.post(endpoint('SIGN_IN', config), data, { headers: headers })
    dispatch(signInSuccess({...response.data, accountID: accountID(response.data)}))
    dispatch({type: 'PROJECTS'})
  } catch (err) {
    dispatch(signInFailure(errorMessage(err)))
  }
}

export const signUp = async (dispatch, config, email, password) => {
  dispatch(signUpAction(email, password))

  const data = { email: email, password: password }
  try {
    await axios.post(endpoint("SIGN_UP", config), data, { headers: headers })
    dispatch(signUpSuccess({ email: email }))
    dispatch({type: 'SIGNIN'})
  } catch (err) {
    dispatch(signUpFailure(errorMessage(err)))
  }
}

const endpoint = (type, config) => (
  apiBaseURL(config) + '/' + (type == 'SIGN_IN' ? 'tokens' : 'users')
)

const errorMessage = err => (
  err.response != undefined ? err.response.data : 'Looks like something is wrong. Please try again.'
)

const accountID = account => (
  Object.keys(account.accounts)[0] // use top-most account as selected for this user
)