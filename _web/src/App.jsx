import React from 'react'

import VersionContainer from './VersionContainer.jsx'
import HeaderContainer from './HeaderContainer.jsx'
import AuthContainer from './auth/AuthContainer.jsx'
import JobsContainer from './jobs/JobsContainer.jsx'
import ProjectsContainer from './projects/ProjectsContainer.jsx'
import JobDetailsContainer from './jobs/JobDetailsContainer.jsx'

const App = ({ location, account, transitionToSignIn, transitionToProjects, config, inJobsComparisonMode }) => (
  componentForLocation(location, account, transitionToSignIn, transitionToProjects, config, inJobsComparisonMode)
)

const componentForLocation = (location, account, transitionToSignIn, transitionToProjects, config, inJobsComparisonMode) => {
  disableScrollingIfNeeded(inJobsComparisonMode)

  if (!isLoggedIn(account) && location != "SIGNIN" && location != "SIGNUP") {
    transitionToSignIn()
    return <div></div>
  }

  var page = null
  switch (location) {
    case 'JOB':
      page = <JobDetailsContainer config={config} account={account} />
      break
    case 'PROJECT':
      page = <JobsContainer config={config} account={account} />
      break
    case 'PROJECTS':
      page = <ProjectsContainer config={config} account={account} />
      break
    case 'SIGNIN':
    case 'SIGNUP':
      if (isLoggedIn(account)) {
        transitionToProjects()
        break
      }
      return <div>
                <AuthContainer config={config} account={account} />
                <VersionContainer version={config.appVersion}/>
            </div>
    default:
      page = <div id='404'>{location} not found (404)</div>
  }

  return <div>
    <HeaderContainer />
    {page}
    <VersionContainer version={config.appVersion}/>
  </div>
}

const isLoggedIn = account => (
  Object.keys(account.account).length != 0 && account.account.constructor === Object
)

const disableScrollingIfNeeded = inJobsComparisonMode => {
  document.body.style.overflow = inJobsComparisonMode ? 'hidden' : 'scroll'
}

export default App
