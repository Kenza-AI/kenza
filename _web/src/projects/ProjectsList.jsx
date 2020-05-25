import React from 'react'
import Filter from '../Filter.jsx'
import Project from './Project.jsx'
import AddProject from './AddProject.jsx'
import { listIcon, gridIcon } from '../icons'

export const emptyProjects = (addProjectHandler, hideAddProjectHandler, createProjectHandler, pendingProject, pendingProjectChangeHandler) => (
  <section className="main-content">

    {pendingProject != null && <AddProject
      pendingProject={pendingProject}
      createProjectHandler={createProjectHandler}
      hideAddProjectHandler={hideAddProjectHandler}
      pendingProjectChangeHandler={pendingProjectChangeHandler} />}

    <div id='no-projects'>
      <p>It's lonely in here...</p>
      <p><b>Create a new project</b> to get started.</p>
      <div>
        <button onClick={addProjectHandler}>Create New Project</button>
      </div>
    </div>

  </section>
)

export const projectsSection = (
  projects,
  accountName,
  filterText,
  filterChangeHandler,
  showAddProjectHandler,
  hideAddProjectHandler,
  createProjectHandler,
  projectsLayout,
  layoutSwitchHandler,
  projectSelectionHandler,
  pendingProject,
  pendingProjectChangeHandler,
  createJobHandler,
  deleteProjectHandler) => (
    <section id="main-content">

      {header(accountName, filterText, filterChangeHandler, showAddProjectHandler, layoutSwitchHandler, projectsLayout)}

      {pendingProject != null && <AddProject
        pendingProject={pendingProject}
        createProjectHandler={createProjectHandler}
        hideAddProjectHandler={hideAddProjectHandler}
        pendingProjectChangeHandler={pendingProjectChangeHandler} />}

      <ul className={projectsListClass(projectsLayout)}>
        {projects.list.filter(project => (project.title.toUpperCase().includes(filterText.toUpperCase())))
          .sort((a, b) => { return new Date(b.created).getTime() - new Date(a.created).getTime() })
          .map(project => <Project
            key={project.projectID}
            layout={projectsListClass(projectsLayout)}
            projectSelectionHandler={projectSelectionHandler}
            createJobHandler={createJobHandler}
            deleteProjectHandler={deleteProjectHandler}
            project={project} />)}
      </ul>
    </section>
  )

const header = (accountName, filterText, filterChangeHandler, showAddProjectHandler) => (
  <header>
    <div>
      <h2 className='projects-header'>{accountName} Projects</h2>
      <button onClick={showAddProjectHandler}>Create New Project</button>
      <Filter text={filterText} placeholder='Search for projects by name' onChangeHandler={filterChangeHandler} />
    </div>
  </header>
)

const projectsListClass = projectsLayout => (
  projectsLayout === 'card' ? 'card' : ''
)
