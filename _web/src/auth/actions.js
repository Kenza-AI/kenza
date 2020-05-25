export const SIGN_OUT = 'SIGN_OUT'
export const TRY_SIGN_IN = 'TRY_SIGN_IN'
export const TRY_SIGN_UP = 'TRY_SIGN_UP'
export const SIGN_IN_SUCCESS = 'SIGN_IN_SUCCESS'
export const SIGN_IN_FAILURE = 'SIGN_IN_FAILURE'
export const SIGN_UP_SUCCESS = 'SIGN_UP_SUCCESS'
export const SIGN_UP_FAILURE = 'SIGN_UP_FAILURE'

export const signOut = () => ({
    type: SIGN_OUT
})

export const signIn = (email, password) => ({
    type: TRY_SIGN_IN,
    email: email,
    password: password
})

export const signInSuccess = account => ({
    type: SIGN_IN_SUCCESS,
    account: account
})

export const signInFailure = error => ({
    type: SIGN_IN_FAILURE,
    error: error
})

export const signUp = (email, password, firstName, lastName) => ({
    type: TRY_SIGN_UP,
    email: email,
    password: password,
    firstName: firstName,
    lastName: lastName
})

export const signUpSuccess = account => ({
    type: SIGN_UP_SUCCESS,
    account: account
})

export const signUpFailure = error => ({
    type: SIGN_UP_FAILURE,
    error: error
})

export const toggleAuthMode = isSignUpMode => ({
    type: isSignUpMode ? 'SIGNIN' : 'SIGNUP'
})
