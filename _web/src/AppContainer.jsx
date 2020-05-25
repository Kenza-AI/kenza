import { connect } from 'react-redux'
import App from './App.jsx'

const mapStateToProps = (state, ownProps) => ({
  location: state.location.type,
  config: ownProps.config,
  account: state.account,
  inJobsComparisonMode: state.jobs.inComparisonMode
})

const mapDispatchToProps = dispatch => ({
  transitionToSignIn: () => dispatch({ type: 'SIGNIN' }),
  transitionToProjects: () => dispatch({ type: 'PROJECTS' })
})

const AppContainer = connect(mapStateToProps, mapDispatchToProps)(App)

export default AppContainer



