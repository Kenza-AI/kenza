import { connect } from 'react-redux'

import Auth from './Auth.jsx'
import { signIn, signUp } from './api'
import { toggleAuthMode } from './actions'

const matchStateToProps = state => ({
  isSigningIn: state.account.isSigningIn,
  isSigningUp: state.account.isSigningUp,
  signInError: state.account.signInError,
  signUpError: state.account.signUpError,
  isSignUpMode: state.location.type == 'SIGNUP'
})

const mapDispatchToProps = (dispatch, ownProps) => ({
  signInHandler: (email, password) => signIn(dispatch, ownProps.config, email, password),
  signUpHandler: (email, password) => signUp(dispatch, ownProps.config, email, password),
  toggleModeHandler: isSignUpMode => dispatch(toggleAuthMode(isSignUpMode))
})

const AuthContainer = connect(matchStateToProps, mapDispatchToProps)(Auth)
export default AuthContainer
