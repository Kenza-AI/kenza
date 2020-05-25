import React from 'react'
import PropTypes from 'prop-types'

const Error = ({message}) => (
    <div className="error-banner">{message}</div>
)

Error.propTypes = {
  message: PropTypes.string.isRequired
}

Error.defaultProps = {
  message: 'Unkown error.'
}

export default Error
