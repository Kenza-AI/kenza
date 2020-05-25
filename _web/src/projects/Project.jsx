import React from 'react'
import PropTypes from 'prop-types'
import { rightArrowIcon, runIcon } from '../icons'

const Project = ({ project, layout, projectSelectionHandler, createJobHandler, deleteProjectHandler }) => (

  <li className={listItemClass(layout)}>
    <button id={project.projectID} className={idClass(layout)}> id: {project.projectID}</button>
    <h3 className={titleClass(layout)}>
      <div>
        {project.title} ({project.branch})
      </div>
      <i>{project.repo}</i>
    </h3>
    <button className="create-job" onClick={() => createJobHandler(project)} title="Run the most recently pushed commit">{runIcon()} Run</button>
    <hr className={separatorClass(layout)} />
    <i className={createdByClass(layout)}>Created by {project.creator}</i>
    <i className={updatedByClass(layout)}>Last updated on {new Date(project.updated).toDateString()}</i>
    <p className={descriptionClass(layout)}>{project.description}</p>
    {<button className="delete-project" onClick={() => deleteProjectHandler(project.projectID)}> Delete </button>}
    {<button className={trainingJobsButtonClass(layout)} onClick={() => projectSelectionHandler({ 'title': project.title, 'id': project.projectID })}> {rightArrowIcon()} </button>}
  </li>
)

const idClass = layout => ('project-id ' + layout)
const titleClass = layout => ('project-title ' + layout)
const separatorClass = layout => ('project-separator ' + layout)
const createdByClass = layout => ('project-created-by ' + layout)
const updatedByClass = layout => ('project-updated-by ' + layout)
const descriptionClass = layout => ('project-description ' + layout)
const listItemClass = layout => ('project ' + layout)
const trainingJobsButtonClass = layout => ('project-training-jobs-button ' + layout)

Project.propTypes = {
  project: PropTypes.object.isRequired,
  layout: PropTypes.string.isRequired,
  createJobHandler: PropTypes.func.isRequired,
  projectSelectionHandler: PropTypes.func.isRequired,
}

export default Project