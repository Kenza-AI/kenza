import React from 'react'
import PropTypes from 'prop-types'

const Filter = ({placeholder, text, onChangeHandler}) => (
  <div className='filter'>
    <input type='text' placeholder={placeholder} value={text} onChange={ e => onChangeHandler(e.target.value)} />
  </div>
)

Filter.propTypes = {
  text: PropTypes.string.isRequired,
  placeholder: PropTypes.string.isRequired,
  onChangeHandler: PropTypes.func.isRequired
}

Filter.defaultProps = {
  text: '',
  placeholder: 'Search...'
}

export default Filter
