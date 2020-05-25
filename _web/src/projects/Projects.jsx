import React from 'react'
import Error from '../Error.jsx'
import { loadingProjects } from './ProjectsLoading.jsx'
import { emptyProjects, projectsSection } from './ProjectsList.jsx'

import PropTypes from 'prop-types'

export default class Projects extends React.Component {

  componentDidMount() {
    this.props.fetchProjects()
  }

  render() {
    const {
      projects,
      accountName,
      filterChangeHandler,
      showAddProjectHandler,
      hideAddProjectHandler,
      createProjectHandler,
      layoutSwitchHandler,
      projectSelectionHandler,
      pendingProject,
      pendingProjectChangeHandler,
      createJobHandler,
      deleteProjectHandler } = this.props

    if (projects.fetchProjectsError != null) {
      return <Error message={projects.fetchProjectsError} />
    } else if (projects.deleteProjectError != null) {
      return <Error message={projects.deleteProjectError} />
    } else if (projects.isFetching) {
      return loadingProjects()
    } else {
      if (!hasProjects(projects)) {
        return emptyProjects(
          showAddProjectHandler,
          hideAddProjectHandler,
          createProjectHandler,
          pendingProject,
          pendingProjectChangeHandler)
      } else {
        return projectsSection(
          projects,
          accountName,
          projects.textFilter,
          filterChangeHandler,
          showAddProjectHandler,
          hideAddProjectHandler,
          createProjectHandler,
          projects.layout,
          layoutSwitchHandler,
          projectSelectionHandler,
          pendingProject,
          pendingProjectChangeHandler,
          createJobHandler,
          deleteProjectHandler)
      }
    }
  }
}

const hasProjects = (projects) => (
  projects.list.length != 0 && projects.hasFetchedOnce
)
