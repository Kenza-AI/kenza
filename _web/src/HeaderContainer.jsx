import { connect } from 'react-redux'
import { signOut } from './auth/actions'
import Header from './Header.jsx'

const mapStateToProps = state => ({
  username: state.account.account.username
})

const mapDispatchToProps = dispatch => ({
    signOutHandler: () => dispatch(signOut()),
    navigateToHomeHandler: () => dispatch({type: 'PROJECTS'}),
    navigateToSettingsHandler: () => dispatch({type: 'SETTINGS'})
})

const HeaderContainer = connect(mapStateToProps, mapDispatchToProps)(Header)

export default HeaderContainer
