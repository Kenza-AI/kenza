import { 
  SIGN_OUT,
  TRY_SIGN_IN,
  SIGN_IN_SUCCESS,
  SIGN_IN_FAILURE,
  TRY_SIGN_UP,
  SIGN_UP_SUCCESS,
  SIGN_UP_FAILURE } from './actions'

export const account = (state=initialState(), action) => {
  switch (action.type) {
    case SIGN_IN_SUCCESS:
      setAccount(action.account)
      return {
        ...state, 
        account: getAccount(), 
        isSigningIn: false,
        signInError: null
      }
    case SIGN_UP_SUCCESS:
      setAccount(action.account)
      return {
        ...state,
        account: getAccount(),
        isSigningUp: false,
        signUpError: null
      }
    case SIGN_IN_FAILURE:
      setAccount(null)
      return {
        ...state,
        account: getAccount(),
        isSigningIn: false,
        signInError: action.error
      }
    case SIGN_UP_FAILURE:
      setAccount(null)
      return {
        ...state,
        account: getAccount(),
        isSigningUp: false,
        signUpError: action.error
      }
    case TRY_SIGN_IN:
      return {
        ...state,
        isSigningIn: true,
        signInError: null,
        isSigningUp: false,
        signUpError: null
      }
    case TRY_SIGN_UP:
      return {
        ...state,
        isSigningIn: false,
        signInError: null,
        isSigningUp: true,
        signUpError: null
      }
    case SIGN_OUT:
      setAccount(null)
      return initialState()
    default:
      return state
  }
}

const setAccount = account => {
  if (account === null) {
    sessionStorage.clear()
  } else {
    sessionStorage.setItem('account', JSON.stringify(account))
  }
}

const getAccount = () => {
  const account = sessionStorage.getItem('account')
  return account === null || account === undefined ? {} : JSON.parse(account)
}

const initialState = () => ({
  account: getAccount(),
  isSigningIn: false,
  isSigningUp: false,
  signInError: null,
  signUpError: null
})