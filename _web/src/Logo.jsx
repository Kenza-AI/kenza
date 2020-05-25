import React from 'react'

const Logo = ({navigateToHomeHandler}) => (
  <button to="/app/projects" className="logo" onClick={navigateToHomeHandler} >
    <img src="kenza_logo.svg" alt="kenza logo" />
  </button>
)

export default Logo
