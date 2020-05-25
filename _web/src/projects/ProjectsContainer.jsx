import { connect } from 'react-redux'

import Projects from './Projects.jsx'
import { fetchProjects, createProject, submitJob, deleteProject } from './api'
import { filterProjects, showAddProject, hideAddProject, pendingProjectChanged } from './actions'

const mapStateToProps = state => ({
    projects: state.projects,
    pendingProject: state.pendingProject
})

const mapDispatchToProps = (dispatch, ownProps) => ({
    // Fetch and filter projects
    fetchProjects: () => fetchProjects(dispatch, ownProps.config, ownProps.account),
    filterChangeHandler: text => dispatch(filterProjects(text)),
    
    // Add new project
    showAddProjectHandler: () => dispatch(showAddProject()),
    hideAddProjectHandler: () => dispatch(hideAddProject()),
    createProjectHandler: project => 
            createProject(dispatch, project, ownProps.config, ownProps.account),
    pendingProjectChangeHandler: project => 
            dispatch(pendingProjectChanged(project)),

    // Select project (go to jobs list)
    projectSelectionHandler: project => dispatch({type: 'PROJECT', payload: { projectID: project.id }}),

    // "Build now", only applicable if building a specific git ref
    createJobHandler: project => submitJob(dispatch, project, ownProps.config, ownProps.account),

    // Delete project
    deleteProjectHandler: projectID => { 
        if (confirm("Are you sure you want to delete this project?")) {
            deleteProject(dispatch, projectID, ownProps.config, ownProps.account)}
        }
})

const ProjectsContainer = connect(mapStateToProps, mapDispatchToProps)(Projects)

export default ProjectsContainer