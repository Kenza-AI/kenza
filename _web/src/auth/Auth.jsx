import React from 'react'
import PropTypes from 'prop-types';
//import LoadingIndicator from '../../../components/LoadingIndicator.jsx'

const Auth = ({ signInHandler, signUpHandler, isSigningIn, isSigningUp, signInError, signUpError, isSignUpMode, toggleModeHandler }) => {
 
  let emailError = null
  let passwordError = null

  let emailInput = null
  let passwordInput = null
  
  return (
    <section id='auth-container'>
      <header>
        <img src='kenza_logo.svg' />
      </header>
      <p>{formTitle(isSignUpMode)}</p>
      {isSignUpMode && <p>While Kenza is in a limited preview release, a new batch of invites will be sent out weekly.</p>}
      {isSignUpMode && <p> Sign up to be added to the waitlist.</p>}
      {signInError && <section className='auth-notification'>{signInError}</section>}
      {signUpError && <section className='auth-notification'>{signUpError}</section>}
      <form className='auth-form'>
        {emailErrorElement()}
        {emailFormElement()}
        {passwordErrorElement()}
        {passwordFormElement()}
        <button type='submit' disabled={isSigningIn || isSigningUp} onClick={ e => validateAndSubmit(isSignUpMode, signInHandler, signUpHandler)}>{mainActionButtonTitle(isSigningIn, isSigningUp, isSignUpMode)}</button>
      </form>
      <div className='toggle-access-mode'> 
        <p>{toggleModePromptTitle(isSignUpMode)}</p>
        <button onClick={ () => toggleModeHandler(isSignUpMode)}>{toggleModeButtonTitle(isSignUpMode)}</button>
      </div>
    </section>
  )

  function emailFormElement() {
    return ( <input type='text' defaultValue='' ref={input => { emailInput = input }} placeholder='Email Address' required /> ) 
  }

  function emailErrorElement() {
    return ( <span ref={error => emailError = error} className='form-error'></span> )
  }

  function passwordFormElement() {
    return ( <input type='password' autoCapitalize='off' autoComplete='off' defaultValue='' ref={input => { passwordInput = input }} placeholder='Password' required /> ) 
  }

  function passwordErrorElement() {
    return ( <span ref={error => passwordError = error} className='form-error'></span> )
  }

  function validateAndSubmit(isSignUpMode, signInHandler, signUpHandler) {
    const email = emailInput.value
    const password = passwordInput.value
    
    let emailErrorText = validateEmail(email)
    let passwordErrorText = validatePassword(password)

    if (emailErrorText) {
      emailInput.classList.add('error')
      emailError.textContent = emailErrorText
    } else {
      emailInput.classList.remove('error')
      emailError.textContent = ''
    }
    if (passwordErrorText) {
      passwordInput.classList.add('error')
      passwordError.textContent = passwordErrorText
    } else {
      passwordInput.classList.remove('error')
      passwordError.textContent = ''
    }

    let isFormValid = emailErrorText == null && passwordErrorText == null
    if (isFormValid) {
      if (isSignUpMode) {
        signUpHandler(email, password)
      } else {
        signInHandler(email, password)
      }
    }
  }
  
  function validateEmail(email) {
    return email.trim() === '' ? 'Email cannot be left blank' : null // TODO: validate email properly, not just empty strings.
  }

  function validatePassword(password) {
    return password.trim().length < 8 ? 'Password should be at least 8 characters long and contain a number, special character, lowercase and uppercase letter' : null
  }
}

const mainActionButtonTitle = (isSigningIn, isSigningUp, signUpMode) => {
  let title = ''
  if (signUpMode) {
    title = isSigningUp ? 'Signing up...' : 'Sign Up'
  } else {
    title = isSigningIn ? 'Signing in...' : 'Sign in'
  }
  return title
}

const formTitle = signUpMode => (
    signUpMode ? 'Sign up for the Private Alpha' : 'Hi there, come on in!'
)

const toggleModePromptTitle = signUpMode => (
    signUpMode ? 'Already have a Kenza account?' : 'Don’t have a Kenza account yet?'
)

const toggleModeButtonTitle = signUpMode => (
    signUpMode ? 'Sign In' : 'Sign Up'
)

Auth.propTypes = {
  isSigningIn: PropTypes.bool.isRequired,
  signInHandler: PropTypes.func.isRequired,
  signInError: PropTypes.string,
  isSignUpMode: PropTypes.bool.isRequired,
  toggleModeHandler: PropTypes.func.isRequired,
  signUpHandler: PropTypes.func.isRequired,
  isSigningUp: PropTypes.bool,
  signUpError: PropTypes.string
}

Auth.defaultProps = {
  isSigningIn: false,
  signInError: null,
  isSignUpMode: false
}

export default Auth
