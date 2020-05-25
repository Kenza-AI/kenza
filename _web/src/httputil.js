export const signOutIfUnauthorized = (err, dispatch) => {
  if (!err.response) { return true }

  if (err.response.status == 401 || err.response.status == 403) {
    dispatch({ type: 'SIGN_OUT' })
    return false
  }
  return true
}

export const apiBaseURL = config => (
  config.apiProtocol + '://' + config.apiHost + ':' + config.apiPort + '/' + config.apiVersion
)
