import React from 'react'
import PropTypes from 'prop-types'

import Logo from './Logo.jsx'
import { dropdownIcon } from './icons'

const Header = ({ username, signOutHandler, navigateToSettingsHandler, navigateToHomeHandler }) => {

  let dropDownMenu = null

  return (
    <header>
      <Logo navigateToHomeHandler={navigateToHomeHandler}/>
      <nav id='account-menu' onMouseOver={() => toggleMenu()} onMouseOut={() => toggleMenu()}>
        <button className='dropdown-toggle-button'>{username}{dropdownIcon()}</button>
        {dropdown(signOutHandler, navigateToSettingsHandler)}
      </nav>
    </header>
  )

  function dropdown(signOutHandler, navigateToSettingsHandler) {
    return (
      <ul id="settings-dropdown" ref={menu => dropDownMenu = menu}>
        <li><button onClick={signOutHandler}>Sign out</button></li>
      </ul>
    )
  }

  function toggleMenu() {
    dropDownMenu.classList.toggle('visible')
  }
}

Header.propTypes = {
  username: PropTypes.string.isRequired,
  signOutHandler: PropTypes.func.isRequired,
  navigateToHomeHandler: PropTypes.func.isRequired,
  navigateToSettingsHandler: PropTypes.func.isRequired
}

export default Header
